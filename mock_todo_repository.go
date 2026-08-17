package main

import (
	"time"

	"github.com/stretchr/testify/mock"
)

type MockTodoRepository struct {
	mock.Mock
}

func (m *MockTodoRepository) GetAll() ([]Todo, error) {
	args := m.Called()
	return args.Get(0).([]Todo), args.Error(1)
}

func (m *MockTodoRepository) Create(todo Todo) (Todo, error) {
	args := m.Called(todo)
	return args.Get(0).(Todo), args.Error(1)
}

func (m *MockTodoRepository) GetByID(id int) (Todo, error) {
	args := m.Called(id)
	return args.Get(0).(Todo), args.Error(1)
}

func (m *MockTodoRepository) GetTodosByCategory(category string) ([]Todo, error) {
	args := m.Called(category)
	return args.Get(0).([]Todo), args.Error(1)
}

func (m *MockTodoRepository) GetTodosByStatus(status bool) ([]Todo, error) {
	args := m.Called(status)
	return args.Get(0).([]Todo), args.Error(1)
}

func (m *MockTodoRepository) Search(q string) ([]Todo, error) {
	args := m.Called(q)
	return args.Get(0).([]Todo), args.Error(1)
}

func (m *MockTodoRepository) Update(id int, newTodo Todo) (Todo, error) {
	args := m.Called(id, newTodo)
	return args.Get(0).(Todo), args.Error(1)
}

func (m *MockTodoRepository) UpdateByCategory(category string, completed bool, completedAt *time.Time) ([]Todo, error) {
	args := m.Called(category, completed, completedAt)
	return args.Get(0).([]Todo), args.Error(1)
}

func (m *MockTodoRepository) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockTodoRepository) DeleteAll() error {
	args := m.Called()
	return args.Error(0)
}
