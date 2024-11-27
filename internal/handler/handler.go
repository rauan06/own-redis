package handler

import (
	"net"

	"github.com/rauan06/own-redis/internal/dal"
	"github.com/rauan06/own-redis/internal/service"
)

func ServerHandle(conn *net.UDPConn) {
	for {
		buffer := make([]byte, 1024)

		n, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			service.HandleErr(conn, addr, "error reading from udp", err)
			continue
		}

		go func() {
			msg, err := service.ParseInput(buffer[:n])
			if err != nil {
				service.HandleErr(conn, addr, "error parsing inputs", err)
				return
			}

			switch msg.Cmd {
			case "ping":
				_, err = conn.WriteToUDP([]byte("PONG\n"), addr)
				if err != nil {
					service.HandleErr(conn, addr, "error writing to udp", err)
				}
			case "set":
				dal.Data.Map[msg.Key] = make(chan string, 1)
				dal.Data.Map[msg.Key] <- msg.Value
				service.HandleOK(conn, addr, "successfully added a new data")

				if msg.PX != 0 {
					go service.ChanTimer(dal.Data, msg.Key, msg.PX)
				}
			case "get":
				if _, ok := dal.Data.Map[msg.Key]; !ok {
					service.HandleGet(conn, addr, "successfully returned (nil) on invalid key", "(nil)")
				} else {
					item := <-dal.Data.Map[msg.Key]
					dal.Data.Map[msg.Key] <- item
					service.HandleGet(conn, addr, "successfully returned value on valid key", item)
				}
			}
		}()
	}
}
