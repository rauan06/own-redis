package models

import "log/slog"

type Config struct {
	Port   int
	Env    string // ('dev', 'local', 'prod')
	Logger *slog.Logger
}
