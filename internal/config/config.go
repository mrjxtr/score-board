// Package config
package config

import (
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
	if err := godotenv.Load(); err != nil {
		panic("Error loading .env file")
	}

	return &Config{
		PORT: os.Getenv("PORT"),
		Defaults: Defaults{
			Colors: DefaulColors,
		},
	}
}
