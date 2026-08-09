package main

import (
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetAll_Success(t *testing.T) {
	testRepo := new(MockTodoRepository)

	testRepo.On("GetAll").Return([]Todo{{ID: 1, Title: "Test Todo"}}, nil)

	req := httptest.NewRequest("GET", "/todos", nil)

	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler := NewHandler(testRepo)
	handler.GetAll(c)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "Test Todo")

}

func TestGetAll_EmptyArray(t *testing.T) {
	// create an instance of repository mock
	testRepo := new(MockTodoRepository)

	// control what mock returns
	testRepo.On("GetAll").Return([]Todo{}, nil)

	req := httptest.NewRequest("GET", "/todos", nil)

	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler := NewHandler(testRepo)
	handler.GetAll(c)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, "[]", w.Body.String())
}

func TestGetAll_DatabaseError(t *testing.T) {
	// create an instance of repository mock
	testRepo := new(MockTodoRepository)

	// control what mock returns
	testRepo.On("GetAll").Return([]Todo{}, errors.New("db error")) // simulate that repo returns an error

	req := httptest.NewRequest("GET", "/todos", nil)

	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler := NewHandler(testRepo)
	handler.GetAll(c)

	assert.Equal(t, 500, w.Code)
	assert.Contains(t, w.Body.String(), "db error")

}
