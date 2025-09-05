package domain

import (
	dto "user-service/internal/app/DTO"
	"user-service/internal/app/DTO/request"
)

type User struct {
	Username string
	Password string
	Role     string
}

func NewUser(username, Password, Role string) *User {
	return &User{
		Username: username,
		Password: Password,
		Role:     Role,
	}
}

type UserRepository interface {
	Add(user *User) *dto.Response
	IsExist(username string) *dto.Response
	UpdateBalance(req *request.RequestFromTransaction) *dto.Response
}
