package service

import (
	"zhashkRestApi"
	"zhashkRestApi/pkg/repository"
)

type TodoItemService struct {
	repo     repository.TodoItem
	listRepo repository.TodoList
}

func NewTodoItemService(repo repository.TodoItem, listRepo repository.TodoList) *TodoItemService {
	return &TodoItemService{
		repo:     repo,
		listRepo: listRepo,
	}
}

func (s *TodoItemService) Create(userId, listId int, item zhashkRestApi.TodoItem) (int, error) {
	_, err := s.listRepo.GetById(userId, listId)
	if err != nil {
		// list doesn't exist or doesn't belong to user
		return 0, err
	}

	return s.repo.Create(listId, item)
}

func (s *TodoItemService) GetAll(userId, listId int) ([]zhashkRestApi.TodoItem, error) {
	return s.repo.GetAll(userId, listId)
}

func (s *TodoItemService) GetById(userId, itemId int) (zhashkRestApi.TodoItem, error) {
	return s.repo.GetById(userId, itemId)
}
