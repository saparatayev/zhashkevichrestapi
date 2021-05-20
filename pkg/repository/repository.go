package repository

import (
	"zhashkRestApi"

	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(user zhashkRestApi.User) (int, error)
	GetUser(username, password string) (zhashkRestApi.User, error)
}

type TodoList interface {
}

type TodoItem interface {
}

type Repository struct {
	Authorization
	TodoList
	TodoItem
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: newAuthMysql(db),
	}
}
