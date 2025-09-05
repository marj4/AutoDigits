package httpserver

import (
	dto "auth-service/internal/application/DTO"
	"auth-service/internal/application/usecase"
	jwtt "auth-service/internal/infrastructure/JWT"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type HttpServer struct {
	usecase *usecase.ClientUseCase
	JWT     *jwtt.JWTService
}

func NewHttpServer(usecase *usecase.ClientUseCase, JWT *jwtt.JWTService) *HttpServer {
	return &HttpServer{
		usecase: usecase,
		JWT:     JWT,
	}
}

func (h *HttpServer) SignUpHandler(c *gin.Context) {
	username, existUsername := c.Get("username")
	password, existPassword := c.Get("password")
	role, existRole := c.Get("role")

	if !existUsername || !existPassword || !existRole {
		c.JSON(http.StatusBadRequest, gin.H{
			"Message": "username or/and password is empty",
		})
		return
	}

	req := &dto.Request{
		Username: username.(string),
		Password: password.(string),
		Role:     role.(string),
	}

	respFromUC := h.usecase.SignUp(req)

	if respFromUC.Status == dto.Error {
		c.JSON(http.StatusConflict, gin.H{
			"Status":  respFromUC.Status,
			"Message": respFromUC.Message,
		})
		return
	}

	// Если пользователь существует ,то сообщаем об этом
	if respFromUC.Status == strconv.Itoa(http.StatusFound) {
		c.JSON(http.StatusConflict, gin.H{
			"Message": respFromUC.Message,
		})
		return

	}

	c.JSON(http.StatusCreated, gin.H{
		"Statusss": respFromUC.Status,
		"Message":  respFromUC.Message,
	})
}

func (h *HttpServer) SignInHandler(c *gin.Context) {
	username, existusername := c.Get("username")
	password, existPassword := c.Get("password")

	if !existusername || !existPassword {
		c.JSON(http.StatusBadRequest, gin.H{
			"Message": "username and password are required",
		})
		return
	}

	req := &dto.Request{
		Username: username.(string),
		Password: password.(string),
	}

	respFromUC := h.usecase.SignIn(req)

	if respFromUC.Status == dto.Error {
		log.Fatal("(SignIn) Error from UC: ", respFromUC.Message)
		return
	}

	// Check for exist
	if respFromUC.Status == strconv.Itoa(http.StatusNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"Message": respFromUC.Message})
		return
	}

	// Проверяем пароль
	if err := bcrypt.CompareHashAndPassword([]byte(respFromUC.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"Message": "Enter correct password",
		})
		return
	}

	log.Println("respFromUC.Uuid: ", respFromUC.Uuid)
	log.Println("respFromUC.Password: ", respFromUC.Password)

	// генерация jwt по uuid
	JWT, err := h.JWT.Generate(respFromUC.Uuid, respFromUC.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Message": "Error generating JWT: " + err.Error(),
		})
		return
	}

	// выдача JWT
	c.Header("Authorization", "Bearer "+JWT)
	c.JSON(200, gin.H{"status": "success"})

}

func (h *HttpServer) Validate(c *gin.Context) {
	UUID, existUUID := c.Get("UUID")
	Role, existRole := c.Get("Role")
	log.Println(UUID)
	log.Println(Role)

	if !existUUID || !existRole {
		c.JSON(http.StatusBadRequest, gin.H{
			"Message": "username and exist are required",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Message": UUID,
		"Role":    Role,
	})

}

func (h *HttpServer) Middleware(c *gin.Context) {
	var request *dto.Request
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Message": "Invalid request format",
			"Error":   err.Error(),
		})
		c.Abort()
		return
	}

	if request.Username == "" || request.Password == "" || request.Role == "" {
		c.JSON(http.StatusConflict, gin.H{
			"Message": "Can`t enter empty data"})
		c.Abort()
		return
	}

	if len(request.Username) < 3 || len(request.Username) > 25 {
		c.JSON(http.StatusConflict, gin.H{"Message": "username must have between 3 to 25 character"})
		c.Abort()
		return
	}

	if len(request.Password) < 8 || len(request.Password) > 15 {
		c.JSON(http.StatusConflict, gin.H{"Message": "Password must have between 8 to 15 character"})
		c.Abort()
		return
	}

	if len(request.Role) > 10 {
		c.JSON(http.StatusConflict, gin.H{"Message": "Role must have character limit 10"})
		c.Abort()
		return
	}

	c.Set("username", request.Username)
	c.Set("password", request.Password)
	c.Set("role", request.Role)

	c.Next()
}

func (h *HttpServer) MiddlewareJWT(c *gin.Context) {
	// Извлекаем header из запроса
	tokenAuth := c.GetHeader("Authorization")
	if tokenAuth == "" {
		log.Fatal("Empty header Authorization")
	}

	// Проверка формата
	tokenParts := strings.Split(tokenAuth, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		log.Fatal("Empty token or Bearer")
	}

	token := tokenParts[1]

	// Валидация токена
	claims, err := h.JWT.Validate(token)
	if err != nil {
		log.Fatal("Error validate JWT: ", err)
	}

	c.Set("UUID", claims.UserUUID)
	c.Set("Role", claims.UserRole)

	c.Next()

}
