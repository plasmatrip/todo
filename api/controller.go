package api

import (
	"todo/repository"
)

type TodoHandlers struct {
	Repo *repository.Repository
}

func NewTodoHandlers(todo *repository.Repository) *TodoHandlers {
	return &TodoHandlers{Repo: todo}
}
