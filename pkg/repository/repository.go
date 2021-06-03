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
	Create(userId int, list zhashkRestApi.TodoList) (int, error)
	GetAll(userId int) ([]zhashkRestApi.TodoList, error)
	GetById(userId, listId int) (zhashkRestApi.TodoList, error)
	Delete(userId, listId int) error
	Update(userId, listId int, input zhashkRestApi.UpdateListInput) error
}

type TodoItem interface {
	Create(listId int, item zhashkRestApi.TodoItem) (int, error)
	GetAll(userId, listId int) ([]zhashkRestApi.TodoItem, error)
}

type Repository struct {
	Authorization
	TodoList
	TodoItem
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: newAuthMysql(db),
		TodoList:      NewTodoListMysql(db),
		TodoItem:      NewTodoItemMysql(db),
	}
}
