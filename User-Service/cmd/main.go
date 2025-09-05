package main

import (
	"fmt"
	"log"
	"user-service/config"
	"user-service/internal/app/usecase"
	"user-service/internal/infrastructure/repository"
	httpserver "user-service/internal/interfaces/http-server"

	"github.com/gin-gonic/gin"
)

var connect_string = "postgres://%s:%s@%s:%s/%s?sslmode=%s"

func main() {
	// Инициализация конфига
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Error load config: ", err)
	}

	// Формирование строки подключения к БД
	connectString := fmt.Sprintf(connect_string, cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name, cfg.SslMode)

	// Создание пула соединений pgxpool
	DB, err := repository.InitDB(connectString)
	if err != nil {
		log.Fatal("Error init DB: ", err)
	}
	defer DB.Close()

	// Инициализация остальных слоев приложения
	repo := repository.NewPostgresRepository(DB)
	usercases := usecase.NewUserUseCase(repo)
	httpsrv := httpserver.NewHttpServer(usercases)

	// Создание и запуск сервера
	router := gin.Default()
	router.POST("/user", httpsrv.AddUserHandler)
	router.GET("/user/:username", httpsrv.CheckExistUserHandler)

	router.Run(":" + cfg.ServerPort)
}
