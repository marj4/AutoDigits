package jwtt

import (
	"auth-service/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService struct {
	Config *config.Config
}

func NewJWTService(cfg *config.Config) *JWTService {
	return &JWTService{Config: cfg}
}

type Claims struct {
	UserUUID string
	UserRole string
	jwt.RegisteredClaims
}

func (j *JWTService) Generate(userUUID string, UserRole string) (string, error) {
	claims := &Claims{
		UserUUID: userUUID,
		UserRole: UserRole,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(j.Config.JWTKey))
}

func (j *JWTService) Validate(tokenString string) (*Claims, error) {
	tokenFromHandler := tokenString
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenFromHandler, claims, func(token *jwt.Token) (any, error) {
		return []byte(j.Config.JWTKey), nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	return claims, nil
}
