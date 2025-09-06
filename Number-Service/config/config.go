package config

import "github.com/ilyakaznacheev/cleanenv"

type Config struct {
	// ...
}

func LoadConfing() (*Config, error) {
	path := "../config.env"

	var cfg Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
