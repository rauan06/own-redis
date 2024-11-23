package config

import (
	"flag"

	"github.com/rauan06/own-redis/models"
)

func ParseFalgs() *models.Config {
	var help bool
	var cfg = models.NewConfig()

	flag.BoolVar(&help, "help", false, "show usage information")
	flag.IntVar(&cfg.Port, "port", 8000, "set port number")
	flag.StringVar(&cfg.Env, "env", "local", "set environment of logger")
	flag.Parse()

	return nil
}
