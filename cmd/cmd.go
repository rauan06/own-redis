package cmd

import (
	"fmt"
	"log/slog"
	"net"
	"sync"

	"github.com/rauan06/own-redis/internal/config"
	"github.com/rauan06/own-redis/models"
)

var cfg *models.Config

func Init() {
	// Initialize config and logger
	cfg = config.SetupConfig()
	slog.SetDefault(cfg.Logger)

	addr := fmt.Sprintf("0.0.0.0:%d", cfg.Port)
	slog.Info("starting udp server", slog.String("Env", cfg.Env), slog.String("addr", addr))

	// Set up UDP listener
	udpAddr := net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: cfg.Port,
		Zone: "",
	}

	conn, err := net.ListenUDP("udp", &udpAddr)
	if err != nil {
		slog.Error("failed to listen on UDP", err)
		return
	}
	defer conn.Close()

	// Store remote connections in sync.Map
	remoteConns := new(sync.Map)

	for {
		buf := make([]byte, 1024)
		_, remoteAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			slog.Error("failed to read from UDP", err)
			continue
		}

		// Check if the remote address already exists in the map
		if _, ok := remoteConns.Load(remoteAddr.String()); !ok {
			remoteConns.Store(remoteAddr.String(), remoteAddr)
		}

		// Goroutine to handle broadcasting to all connected clients
		go func() {
			remoteConns.Range(func(key, value interface{}) bool {
				addr := value.(*net.UDPAddr)
				// Send the data to each connected remote address
				if _, err := conn.WriteToUDP(buf, addr); err != nil {
					slog.Error("failed to send UDP packet", slog.String("addr", addr.String()), err)
					// Remove address from map if it fails to send
					remoteConns.Delete(key)
				}
				return true
			})
		}()
	}
}
