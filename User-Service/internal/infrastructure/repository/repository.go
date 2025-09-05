package repository

import (
	"database/sql"
	"log"
	dto "user-service/internal/app/DTO"
	"user-service/internal/app/DTO/request"
	"user-service/internal/domain"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{
		db: db,
	}
}

func (r *PostgresRepository) Add(userInfo *domain.User) *dto.Response {
	query := `INSERT INTO Users (username,password,role) VALUES ($1,$2,$3) RETURNING id`
	var id string
	if err := r.db.QueryRow(query, userInfo.Username, userInfo.Password, userInfo.Role).Scan(&id); err != nil {
		log.Println("Error adding user:", err)
		return &dto.Response{Status: dto.Error}
	}

	return &dto.Response{Status: dto.Success}
}

func (r *PostgresRepository) IsExist(username string) *dto.Response {
	query := `SELECT id,password,role FROM Users WHERE username = $1`
	row := r.db.QueryRow(query, username)
	var uuidFromDB string
	var passwordFromDB string
	var roleFromDB string

	if err := row.Scan(&uuidFromDB, &passwordFromDB, &roleFromDB); err != nil {
		log.Println(username)
		if err == sql.ErrNoRows {
			return &dto.Response{Status: dto.NotExist}
		}
		log.Println("(IsExist) Error checking user existence:", err)
		return &dto.Response{Status: dto.Error, Message: err.Error()}
	}

	return &dto.Response{
		Status:   dto.Exist,
		UUID:     uuidFromDB,
		Password: passwordFromDB,
		Role:     roleFromDB,
	}
}

func (r *PostgresRepository) UpdateBalance(req *request.RequestFromTransaction) *dto.Response {
	switch req.TxType {
	case "deposit":
		// implement deposit balance
	case "withdraw":
		// implement withdraw balance
	default:
		return &dto.Response{
			Status:  dto.Error,
			Message: "invalid operation" + req.TxType,
		}
	}

	return &dto.Response{Status: dto.Success}
}
