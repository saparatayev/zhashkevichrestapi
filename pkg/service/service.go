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
	Create(userId int, list zhashkRestApi.TodoList) (int, error)
	GetAll(userId int) ([]zhashkRestApi.TodoList, error)
	GetById(userId, listId int) (zhashkRestApi.TodoList, error)
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
		TodoList:      NewTodoListService(repos.TodoList),
	}
}
