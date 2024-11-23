package models

import "log/slog"

type Config struct {
	Port   int
	Addr   string
	Env    string // ('dev', 'local', 'prod')
	Logger slog.Logger
}

func NewConfig() *Config {
	return &Config{}
}
