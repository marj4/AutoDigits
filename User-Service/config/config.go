package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Type       string `env:"DB_TYPE"`
	Host       string `env:"DB_HOST"`
	Port       string `env:"DB_PORT"`
	User       string `env:"DB_USER"`
	Name       string `env:"DB_NAME"`
	Password   string `env:"DB_PASSWORD"`
	SslMode    string `env:"DB_SSL"`
	ServerPort string `env:"SERVER_PORT"`
}

func LoadConfig() (*Config, error) {
	// Путь до .env файла
	path := "./config.env"

	var cfg Config

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil

}
