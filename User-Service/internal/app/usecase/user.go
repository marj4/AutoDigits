package usecase

import (
	dto "user-service/internal/app/DTO"
	"user-service/internal/app/DTO/request"
	"user-service/internal/domain"
)

type UserUseCase struct {
	repo domain.UserRepository
}

func NewUserUseCase(repo domain.UserRepository) *UserUseCase {
	return &UserUseCase{
		repo: repo,
	}
}

func (u *UserUseCase) AddUser(user *request.RequestUserInfo) *dto.Response {
	newUser := domain.NewUser(user.Username, user.Password, user.Role)
	return u.repo.Add(newUser)
}

func (u *UserUseCase) CheckExistUser(username string) *dto.Response {
	return u.repo.IsExist(username)
}

func (u *UserUseCase) UpdateBalance(req *request.RequestFromTransaction) *dto.Response {
	return u.repo.UpdateBalance(req)
}
