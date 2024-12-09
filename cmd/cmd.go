package cmd

import (
	"fmt"
	"log/slog"
	"net"
	"os"

	"github.com/rauan06/own-redis/internal/config"
	"github.com/rauan06/own-redis/internal/dal"
	"github.com/rauan06/own-redis/internal/handler"
	"github.com/rauan06/own-redis/models"
)

var cfg *models.Config

func Init() {
	_, err := os.Open("123")
	fmt.Printf("%+v\n\n", err)

	cfg = config.SetupConfig()
	slog.SetDefault(cfg.Logger)

	addr := fmt.Sprintf("0.0.0.0:%d", cfg.Port)
	slog.Info("starting udp server...", slog.String("Env", cfg.Env), slog.String("addr", addr))

	udpAddr := net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: cfg.Port,
	}

	conn, err := net.ListenUDP("udp", &udpAddr)
	if err != nil {
		slog.Error("error starting udp server", slog.String("error", fmt.Sprint(err)))
		return
	}
	defer conn.Close()

	dal.Data = &models.AsyncMap{
		Map: make(map[string](models.ChanContextPair)),
	}

	handler.ServerHandle(conn)
}
