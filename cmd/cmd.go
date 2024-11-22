package cmd

import (
	"own-redis/config"
	"own-redis/models"
)

var cfg *models.Config

func Init() {
	config.ParseFalgs()
}
