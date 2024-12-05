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

// ParseInput processes the input byte slice and returns a parsed Message or an error.
func ParseInput(input []byte) (*models.Messege, error) {
	args := strings.Fields(string(input))
	if len(args) == 0 {
		return nil, fmt.Errorf("(error) ERR empty input")
	}

	cmd := strings.ToLower(args[0])
	msg := &models.Messege{
		Cmd: cmd,
		PX:  0,
	}

	switch cmd {
	case "ping":
		return msg, nil
	case "set":
		return parseSetCommand(args, msg)
	case "get":
		return parseGetCommand(args, msg)
	default:
		return nil, fmt.Errorf("(error) ERR invalid command")
	}
}

func parseSetCommand(args []string, msg *models.Messege) (*models.Messege, error) {
	if len(args) < 3 {
		return nil, fmt.Errorf("(error) ERR wrong number of arguments for 'SET' command")
	}

	msg.Key = args[1]
	msg.Value = parseSetValue(args)
	if pxIndex := findPxIndex(args); pxIndex != -1 {
		if err := setPxValue(args, pxIndex, msg); err != nil {
			return nil, err
		}
	}
	return msg, nil
}

func parseGetCommand(args []string, msg *models.Messege) (*models.Messege, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("(error) ERR wrong number of arguments for 'GET' command")
	}
	msg.Key = args[1]
	return msg, nil
}

func parseSetValue(args []string) string {
	var valueBuilder strings.Builder
	for i := 2; i < len(args); i++ {
		if args[i] != "px" {
			valueBuilder.WriteString(" " + args[i])
		} else {
			break
		}
	}
	return strings.TrimSpace(valueBuilder.String())
}

func findPxIndex(args []string) int {
	for i, arg := range args {
		if arg == "px" {
			return i
		}
	}
	return -1
}

func setPxValue(args []string, pxIndex int, msg *models.Messege) error {
	if len(args) <= pxIndex+1 {
		return fmt.Errorf("(error) ERR missing PX value")
	}

	n, err := strconv.Atoi(args[pxIndex+1])
	if err != nil {
		return fmt.Errorf("(error) ERR PX value isn't a number")
	}

	if n < 0 {
		return fmt.Errorf("(error) ERR PX value cannot be negative")
	}

	msg.PX = time.Duration(n) * time.Millisecond
	return nil
}

// ChanTimer waits for the expiration time and closes the channel if it exists.
func ChanTimer(model *models.AsyncMap, key string, expiration time.Duration) {
	select {
	case <-model.Map[key].Cancel:
		return
	case <-time.After(expiration):
		if pair, exists := model.Map[key]; exists {
			close(pair.Cancel)
			close(pair.Data)
			delete(model.Map, key)
		}
	}
}

// HandleErr logs the error and writes the error message to the UDP connection.
func HandleErr(conn *net.UDPConn, addr *net.UDPAddr, msg string, err error) {
	logError(msg, err)
	writeResponseToUDP(conn, addr, fmt.Sprint(err))
}

// HandleOK logs the success message and writes "OK" to the UDP connection.
func HandleOK(conn *net.UDPConn, addr *net.UDPAddr, msg string) {
	slog.Info(msg, slog.String("status", "OK"))
	writeResponseToUDP(conn, addr, "OK")
}

// HandleGet logs the GET result and writes the returned value to the UDP connection.
func HandleGet(conn *net.UDPConn, addr *net.UDPAddr, msg string, value string) {
	slog.Error(msg, slog.String("returned", value))
	writeResponseToUDP(conn, addr, value)
}

func writeResponseToUDP(conn *net.UDPConn, addr *net.UDPAddr, response string) {
	if err := writeToUDP(conn, addr, []byte(response+"\n")); err != nil {
		logError("error writing to connection", err)
	}
}

func logError(msg string, err error) {
	slog.Error(msg, slog.String("error", fmt.Sprint(err)))
}

func writeToUDP(conn *net.UDPConn, addr *net.UDPAddr, buffer []byte) error {
	_, err := conn.WriteToUDP(buffer, addr)
	if err != nil {
		return fmt.Errorf("(error) ERR writing to connection: %v", err)
	}
	return nil
}
