package handler

import (
	"context"
	"net"

	"github.com/rauan06/own-redis/internal/dal"
	"github.com/rauan06/own-redis/internal/service"
	"github.com/rauan06/own-redis/models"
)

func ServerHandle(conn *net.UDPConn) {
	for {
		buffer := make([]byte, 1024)

		n, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			service.HandleErr(conn, addr, "error reading from udp", err)
			continue
		}

		go handleMessage(conn, addr, buffer[:n])
	}
}

func handleMessage(conn *net.UDPConn, addr *net.UDPAddr, buffer []byte) {
	msg, err := service.ParseInput(buffer)
	if err != nil {
		service.HandleErr(conn, addr, "error parsing inputs", err)
		return
	}

	switch msg.Cmd {
	case "ping":
		handlePing(conn, addr)
	case "set":
		handleSet(conn, addr, msg)
	case "get":
		handleGet(conn, addr, msg)
	}
}

func handlePing(conn *net.UDPConn, addr *net.UDPAddr) {
	_, err := conn.WriteToUDP([]byte("PONG\n"), addr)
	if err != nil {
		service.HandleErr(conn, addr, "error writing to udp", err)
	}
}

func handleSet(conn *net.UDPConn, addr *net.UDPAddr, msg *models.Messege) {
	// Cancel the existing context if the key already exists
	if pair, exists := dal.Data.Map[msg.Key]; exists && pair.Context != nil {
		pair.CancelFunc() // Cancel the existing context
	}

	// Create a new context with cancel functionality
	ctx, cancel := context.WithCancel(context.Background())

	// Store the context and the cancel function
	dal.Data.Map[msg.Key] = &models.ChanContextPair{
		Data:       make(chan string),
		Context:    ctx,
		CancelFunc: cancel,
	}

	// Send the value into the channel
	dal.Data.Map[msg.Key].Data <- msg.Value
	service.HandleOK(conn, addr, "successfully added a new data")

	// Set expiration if PX is provided
	if msg.PX != 0 {
		go service.ChanTimer(dal.Data, msg.Key, msg.PX)
	}
}

func handleGet(conn *net.UDPConn, addr *net.UDPAddr, msg *models.Messege) {
	if _, ok := dal.Data.Map[msg.Key]; !ok {
		service.HandleGet(conn, addr, "successfully returned (nil) on invalid key", "(nil)")
		return
	}

	item := <-dal.Data.Map[msg.Key].Data
	dal.Data.Map[msg.Key].Data <- item
	service.HandleGet(conn, addr, "successfully returned value on valid key", item)
}
