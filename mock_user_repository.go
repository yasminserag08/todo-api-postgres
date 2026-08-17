package main

import (
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) CreateUser(user User) (User, error) {
	args := m.Called(user)
	return args.Get(0).(User), args.Error(1)
}

func (m *MockUserRepository) GetUserByUsername(username string) (User, error) {
	args := m.Called(username)
	return args.Get(0).(User), args.Error(1)
}

func (m *MockUserRepository) GetUserByID(id uint) (User, error) {
	args := m.Called(id)
	return args.Get(0).(User), args.Error(1)
}
