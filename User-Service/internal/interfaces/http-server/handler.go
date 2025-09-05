package httpserver

import (
	"log"
	"net/http"
	dto "user-service/internal/app/DTO"
	"user-service/internal/app/DTO/request"
	"user-service/internal/app/usecase"

	"github.com/gin-gonic/gin"
)

type HttpServer struct {
	Usecase *usecase.UserUseCase
}

func NewHttpServer(usecase *usecase.UserUseCase) *HttpServer {
	return &HttpServer{Usecase: usecase}
}

func (h *HttpServer) AddUserHandler(c *gin.Context) {
	var user *request.RequestUserInfo
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"Status": dto.Success})
	}

	responseFromUC := h.Usecase.AddUser(user)
	if responseFromUC.Status == dto.Error {
		return
	}

	c.JSON(http.StatusOK, gin.H{"Status": dto.Success})

	go func() {
		log.Println("Отправка письма на поту")
	}()

}

func (h *HttpServer) CheckExistUserHandler(c *gin.Context) {
	username := c.Param("username")
	log.Println("(US) Username from AS", username)

	responseFromUC := h.Usecase.CheckExistUser(username)
	if responseFromUC.Status == dto.Exist {
		c.JSON(http.StatusFound, gin.H{"UUID": responseFromUC.UUID, "Password": responseFromUC.Password, "Role": responseFromUC.Role})
		return
	}

	c.JSON(http.StatusNotFound, nil)
}
