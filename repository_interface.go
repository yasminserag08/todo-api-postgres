package main

import "time"

type TodoRepositoryInterface interface {
	GetAll() ([]Todo, error)
	Create(todo Todo) (Todo, error)
	GetByID(id int) (Todo, error)
	GetTodosByCategory(category string) ([]Todo, error)
	GetTodosByStatus(status bool) ([]Todo, error)
	Search(q string) ([]Todo, error)
	Update(id int, newTodo Todo) (Todo, error)
	UpdateByCategory(category string, completed bool, completedAt *time.Time) ([]Todo, error)
	Delete(id int) error
	DeleteAll() error
}
