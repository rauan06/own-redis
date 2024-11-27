package service

import (
	"fmt"
	"log/slog"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/rauan06/own-redis/models"
)

func ParseInput(input []byte) (*models.Messege, error) {
	args := strings.SplitAfter(string(input), " ")

	if len(args) == 0 {
		return nil, fmt.Errorf("(error) ERR empty input")
	}

	for i := range args {
		args[i] = strings.TrimSpace(args[i])
	}

	cmd := strings.ToLower(args[0])
	msg := &models.Messege{}
	msg.Cmd = cmd
	msg.PX = 0

	switch cmd {
	case "ping":
		return msg, nil
	case "set":
		if len(args) != 3 && len(args) != 5 {
			return nil, fmt.Errorf("(error) ERR wrong number of arguments for 'SET' command")
		}

		msg.Key = args[1]
		msg.Value = args[2]

		if len(args) == 5 {
			if args[3] != "px" {
				return nil, fmt.Errorf("(error) ERR incorrect arguments")
			}

			n, err := strconv.Atoi(args[4])
			if err != nil {
				return nil, fmt.Errorf("(error) ERR PX value isn't a nubmer")
			}

			if n < 0 {
				return nil, fmt.Errorf("(error) ERR PX value cannot be negative")
			}

			msg.PX = time.Duration(n) * time.Millisecond
		}

		return msg, nil

	case "get":
		if len(args) != 2 {
			return nil, fmt.Errorf("(error) ERR wrong number of arguments for 'GET' command")
		}
		msg.Key = args[1]

		return msg, nil

	default:
		return nil, fmt.Errorf("(error) ERR invalid command")
	}
}

func ChanTimer(model *models.AsyncMap, key string, expiration time.Duration) {
	time.Sleep(expiration)
	ch, exists := model.Map[key]
	if exists {
		close(ch)
		delete(model.Map, key)
	}
}

func HandleErr(conn *net.UDPConn, addr *net.UDPAddr, msg string, err error) {
	slog.Error(msg, slog.String("error", fmt.Sprint(err)))

	if err = writeToUDP(conn, addr, []byte(fmt.Sprint(err)+"\n")); err != nil {
		slog.Error("error writing to connection", slog.String("error", fmt.Sprint(err)))
	}
}

func HandleOK(conn *net.UDPConn, addr *net.UDPAddr, msg string) {
	slog.Error(msg, slog.String("status", "OK"))

	if err := writeToUDP(conn, addr, []byte("OK\n")); err != nil {
		slog.Error("error writing to connection", slog.String("error", fmt.Sprint(err)))
	}
}

func HandleGet(conn *net.UDPConn, addr *net.UDPAddr, msg string, value string) {
	slog.Error(msg, slog.String("returned", value))

	if err := writeToUDP(conn, addr, []byte(value+"\n")); err != nil {
		slog.Error("error writing to connection", slog.String("error", fmt.Sprint(err)))
	}
}

func writeToUDP(conn *net.UDPConn, addr *net.UDPAddr, buffer []byte) error {
	_, err := conn.WriteToUDP(buffer, addr)
	if err != nil {
		return fmt.Errorf("(error) ERR writing to connection: %v", err)
	}
	return nil
}
