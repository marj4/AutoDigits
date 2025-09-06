package main

import (
	"Number-Service/config"
	"log"
)

func main() {
	// Конфиг
	cfg, err := config.LoadConfing()
	if err != nil {
		log.Fatal("Error load config", err)
	}
	//Логгер

	// Подключение к БД
	// пулы + гоуртины

	// слои

	// запуск сервера...
	// Использовать GRPC + горутины
}
