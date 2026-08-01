package main

import (
	"database/sql"
	"errors"
)

var ErrNotFound = errors.New("not found")

// instead of having to pass the connection to every function
type TodoRepository struct {
	db *sql.DB
}

func NewTodoRepository(db *sql.DB) *TodoRepository {
	return &TodoRepository{db: db}
}

func (r *TodoRepository) GetAll() ([]Todo, error) {
	todos := []Todo{}
	rows, err := r.db.Query("SELECT id, title, completed FROM todos")

	if err != nil {
		return todos, err
	}
	defer rows.Close()

	for rows.Next() {
		var todo Todo
		err := rows.Scan(&todo.ID, &todo.Title, &todo.Completed)

		if err != nil {
			return todos, err
		}

		todos = append(todos, todo)
	}

	if rows.Err() != nil {
		return todos, rows.Err()
	}

	return todos, nil
}

func (r *TodoRepository) Create(title string) (Todo, error) {
	var todo Todo
	err := r.db.QueryRow(
		"INSERT INTO todos (title) VALUES ($1) RETURNING id, title, completed",
		title,
	).Scan(&todo.ID, &todo.Title, &todo.Completed) // insert and also return the new todo

	if err != nil {
		return Todo{}, err
	}

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
	err := r.db.QueryRow("UPDATE todos SET title = $1, completed = $2 WHERE id = $3 RETURNING id, title, completed", title, completed, id).Scan(&todo.ID, &todo.Title, &todo.Completed)

	if err == sql.ErrNoRows {
		return Todo{}, ErrNotFound
	}
	if err != nil {
		return Todo{}, err
	}

	return todo, nil
}

func (r *TodoRepository) Delete(id int) error {
	result, err := r.db.Exec("DELETE FROM todos WHERE id = $1", id)
	if err != nil {
		return err
	}

	// Exec doesn't return ErrNoRows like QueryRow, so we need to check the number of rows affected
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}
