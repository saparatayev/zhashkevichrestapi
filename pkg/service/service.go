package service

import (
	"zhashkRestApi"
	"zhashkRestApi/pkg/repository"
)

type Authorization interface {
	CreateUser(user zhashkRestApi.User) (int, error)
	GenerateToken(username, pasword string) (string, error)
	ParseToken(token string) (int, error)
}

type TodoList interface {
}

type TodoItem interface {
}

type Service struct {
	Authorization
	TodoList
	TodoItem
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
	}
}
