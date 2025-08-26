// Package config
package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	PORT     string
	Defaults Defaults
}

type Defaults struct {
	Colors map[string]string
}

var DefaulColors = map[string]string{
	"pink":   "#D50059",
	"red":    "#C50000",
	"blue":   "#1D03AF",
	"yellow": "#FFBB02",
}

func LoadConfig() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file")
	}

	return &Config{
		PORT: port,
		Defaults: Defaults{
			Colors: DefaulColors,
		},
	}
}
