package main

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetAll(t *testing.T) {
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
