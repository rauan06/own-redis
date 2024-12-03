package handler

import (
	"context"
	"log/slog"
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
	slog.Info("successfully responded to ping with pong")
}

func handleSet(conn *net.UDPConn, addr *net.UDPAddr, msg *models.Messege) {
	if pair, exists := dal.Data.Map[msg.Key]; exists && pair.Context != nil {
		pair.CancelFunc()
	}

	ctx, cancel := context.WithCancel(context.Background())

	dal.Data.Map[msg.Key] = models.ChanContextPair{
		Data:       make(chan string, 1),
		Context:    ctx,
		CancelFunc: cancel,
	}

	dal.Data.Map[msg.Key].Data <- msg.Value
	service.HandleOK(conn, addr, "successfully added a new data")

	if msg.PX != 0 {
		service.ChanTimer(dal.Data, msg.Key, msg.PX)
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
