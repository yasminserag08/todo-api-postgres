package main

import (
	"errors"

	"gorm.io/gorm"
)

var ErrNotFound = errors.New("not found")

// instead of having to pass the connection to every function
type TodoRepository struct {
	db *gorm.DB
}

func NewTodoRepository(db *gorm.DB) *TodoRepository {
	return &TodoRepository{db: db}
}

func (r *TodoRepository) GetAll() ([]Todo, error) {
	todos := []Todo{}

	if err := r.db.Find(&todos).Error; err != nil {
		return nil, err
	}

	return todos, nil
}

func (r *TodoRepository) Create(todo Todo) (Todo, error) {
	result := r.db.Create(&todo)

	return todo, result.Error
}

func (r *TodoRepository) GetByID(id int) (Todo, error) {
	todo := Todo{}

	result := r.db.First(&todo, id)

	if result.Error == gorm.ErrRecordNotFound {
		return Todo{}, ErrNotFound
	}

	return todo, result.Error
}

func (r *TodoRepository) GetTodosByCategory(category string) ([]Todo, error) {
	todos := []Todo{}

	result := r.db.Where("category = ?", category).Find(&todos)

	return todos, result.Error
}

func (r *TodoRepository) GetTodosByStatus(status bool) ([]Todo, error) {
	todos := []Todo{}

	result := r.db.Where("completed = ?", status).Find(&todos)

	return todos, result.Error
}

func (r *TodoRepository) Search(q string) ([]Todo, error) {
	todos := []Todo{}

	result := r.db.Where("title ILIKE ?", "%"+q+"%").Find(&todos) // ilike because case insensitive

	return todos, result.Error
}

func (r *TodoRepository) Update(id int, title string, completed bool) (Todo, error) {
	result := r.db.Model(&Todo{}).Where("id = ?", id).Updates(map[string]interface{}{"title": title, "completed": completed})

	if result.Error != nil {
		return Todo{}, result.Error
	}
	if result.RowsAffected == 0 {
		return Todo{}, ErrNotFound
	}

	return r.GetByID(id)
}

func (r *TodoRepository) UpdateByCategory(category string, completed bool) ([]Todo, error) {
	var todos []Todo

	result := r.db.Model(&Todo{}).Where("category = ?", category).Update("completed", completed)

	if result.Error != nil {
		return []Todo{}, result.Error
	}
	if result.RowsAffected == 0 {
		return []Todo{}, ErrNotFound
	}

	r.db.Where("category = ?", category).Find(&todos)

	return todos, nil
}

func (r *TodoRepository) Delete(id int) error {
	result := r.db.Delete(&Todo{}, id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *TodoRepository) DeleteAll() error {
	result := r.db.Exec("TRUNCATE TABLE todos RESTART IDENTITY")
	return result.Error
}
