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
	Delete(userId, listId int) error
	Update(userId, listId int, input zhashkRestApi.UpdateListInput) error
}

type TodoItem interface {
	Create(userId, listId int, item zhashkRestApi.TodoItem) (int, error)
	GetAll(userId, listId int) ([]zhashkRestApi.TodoItem, error)
	GetById(userId, itemId int) (zhashkRestApi.TodoItem, error)
	Delete(userId, itemId int) error
	Update(userId, itemId int, input zhashkRestApi.UpdateItemInput) error
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
		TodoItem:      NewTodoItemService(repos.TodoItem, repos.TodoList),
	}
}
