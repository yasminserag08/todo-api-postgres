package main

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSignUpSuccess(t *testing.T) {
	testRepo := new(MockUserRepository)

	testRepo.On("CreateUser", mock.AnythingOfType("User")).Return(User{ID: 1, Username: "yasmin", Role: "user"}, nil)

	body := strings.NewReader(`{"username": "yasmin", "role":"user", "password":"test123"}`)
	req := httptest.NewRequest("POST", "/signup", body)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Set("userID", uint(1))
	c.Request = req

	authHandler := NewAuthHandler(testRepo)
	authHandler.SignUp(c)
	assert.Equal(t, 201, w.Code)
	assert.Contains(t, w.Body.String(), "yasmin")
}

// success when user owns the todo
func TestDelete_OwnTodoSuccess(t *testing.T) {
	testRepo := new(MockTodoRepository)

	testRepo.On("GetByID", 1).Return(Todo{ID: 1, Title: "Test Todo", UserID: 1}, nil)
	testRepo.On("Delete", 1).Return(nil)

	req := httptest.NewRequest("DELETE", "/todos/1", nil)
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "1"}}
	c.Set("userID", uint(1))
	c.Set("role", "user")
	c.Request = req

	handler := NewHandler(testRepo)
	handler.Delete(c)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "Todo deleted")
}

// normal user trying to delete other user's todo
func TestDelete_ForbiddenOtherUserTodo(t *testing.T) {
	testRepo := new(MockTodoRepository)

	testRepo.On("GetByID", 1).Return(Todo{ID: 1, Title: "Someone else's Todo", UserID: 2}, nil)

	req := httptest.NewRequest("DELETE", "/todos/1", nil)
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "1"}}
	c.Set("userID", uint(1))
	c.Set("role", "user")
	c.Request = req

	handler := NewHandler(testRepo)
	handler.Delete(c)

	assert.Equal(t, 403, w.Code)
	assert.Contains(t, w.Body.String(), "you can only delete your own todos")
}

// admin can delete any todo
func TestDelete_AdminSuccess(t *testing.T) {
	testRepo := new(MockTodoRepository)

	// Todo belongs to User ID 2, but admin is trying to delete it
	testRepo.On("GetByID", 1).Return(Todo{ID: 1, Title: "Someone else's Todo", UserID: 2}, nil)
	testRepo.On("Delete", 1).Return(nil)

	req := httptest.NewRequest("DELETE", "/todos/1", nil)
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "1"}}
	c.Set("userID", uint(99))
	c.Set("role", "admin")
	c.Request = req

	handler := NewHandler(testRepo)
	handler.Delete(c)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "Todo deleted")
}
