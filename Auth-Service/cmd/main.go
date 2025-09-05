package main

import (
	"auth-service/config"
	"auth-service/internal/application/usecase"
	jwtt "auth-service/internal/infrastructure/JWT"
	httpserver "auth-service/internal/interfaces/http-server"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// Инициализация конфига
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Error load config: ", err)
	}

	urlUserService := "http://localhost:8081"
	client := http.Client{}

	// Инициализация слоя usecase
	usecases := usecase.NewClientUseCase(urlUserService, &client)

	// Инициализация JWT-сервиса
	JWT := jwtt.NewJWTService(cfg)

	// Инициализация слоя хэндлеров
	httpsrv := httpserver.NewHttpServer(usecases, JWT)

	// Создание сервера
	router := gin.Default()

	router.POST("/signup", httpsrv.Middleware, httpsrv.SignUpHandler)
	router.POST("/signin", httpsrv.Middleware, httpsrv.SignInHandler)
	router.GET("/validateJWT", httpsrv.MiddlewareJWT, httpsrv.Validate)

	router.Run(":" + cfg.ServerPort)
}
