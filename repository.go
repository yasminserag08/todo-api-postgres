package main

import (
	"database/sql"
	"time"
)

// instead of having to
type TodoRepository struct {
	db *sql.DB
}

func NewTodoRepository(db *sql.DB) *TodoRepository {
	return &TodoRepository{db: db}
}

func (r *TodoRepository) GetAll() ([]Todo, error) {
	rows, err := r.db.Query("SELECT * FROM todos")

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var todos []Todo

	for rows.Next() {
		var todo Todo
		var time time.Time
		err := rows.Scan(&todo.ID, &todo.Title, &todo.Completed, &time)

		if err != nil {
			return nil, err
		}

		todos = append(todos, todo)
	}

	return todos, nil
}

func (r *TodoRepository) Create(title string) (Todo, error) {
	var todo Todo
	return todo, nil
}

func (r *TodoRepository) GetByID(id int) (Todo, error) {
	var todo Todo
	return todo, nil
}

func (r *TodoRepository) Update(id int, title string, completed bool) (Todo, error) {
	var todo Todo
	return todo, nil
}

func (r *TodoRepository) Delete(id int) error {
	return nil
}
