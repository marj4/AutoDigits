package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	JWTKey     string `env:"JWT_KEY"`
	ServerPort string `env:"SERVER_PORT"`
}

func LoadConfig() (*Config, error) {
	// Путь до .env файла
	path := "./config.env"

	var cfg Config

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("cannot read config: " + err.Error())
	}

	return &cfg, nil

}
