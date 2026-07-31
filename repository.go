package main

import "database/sql"

// instead of having to
type TodoRepository struct {
	db *sql.DB
}

func NewTodoRepository(db *sql.DB) *TodoRepository {
	return &TodoRepository{db: db}
}

func (r *TodoRepository) GetAll() ([]Todo, error) {
	// TODO
}

func (r *TodoRepository) Create(title string) (Todo, error) {
	// TODO
}

func (r *TodoRepository) GetByID(id int) (Todo, error) {
	// TODO
}

func (r *TodoRepository) Update(id int, title string, completed bool) (Todo, error) {
	// TODO
}

func (r *TodoRepository) Delete(id int) error {
	// TODO
}
