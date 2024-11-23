package cmd

import (
	"github.com/rauan06/own-redis/internal/config"
	"github.com/rauan06/own-redis/models"
)

var cfg *models.Config

func Init() {
	cfg = config.ParseFalgs()
}
