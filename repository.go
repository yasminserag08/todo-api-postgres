package main

import (
	"database/sql"
	"fmt"
)

var ErrNotFound = fmt.Errorf("not found")

// instead of having to pass the connection to every function
type TodoRepository struct {
	db *sql.DB
}

func NewTodoRepository(db *sql.DB) *TodoRepository {
	return &TodoRepository{db: db}
}

func (r *TodoRepository) GetAll() ([]Todo, error) {
	rows, err := r.db.Query("SELECT id, title, completed FROM todos")

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var todos []Todo

	for rows.Next() {
		var todo Todo
		err := rows.Scan(&todo.ID, &todo.Title, &todo.Completed)

		if err != nil {
			return nil, err
		}

		todos = append(todos, todo)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return todos, nil
}

func (r *TodoRepository) Create(title string) (Todo, error) {
	var todo Todo
	return todo, nil
}

func (r *TodoRepository) GetByID(id int) (Todo, error) {
	var todo Todo

	row := r.db.QueryRow("SELECT id, title, completed FROM todos WHERE id = $1", id)
	err := row.Scan(&todo.ID, &todo.Title, &todo.Completed)

	if err == sql.ErrNoRows {
		return Todo{}, ErrNotFound
	} else if err != nil {
		return Todo{}, err
	}

	return todo, nil
}

func (r *TodoRepository) Update(id int, title string, completed bool) (Todo, error) {
	var todo Todo
	return todo, nil
}

func (r *TodoRepository) Delete(id int) error {
	return nil
}
