package main

import (
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// GET /todos
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

// GET /todos/:id
func TestGetByID_Success(t *testing.T) {
	testRepo := new(MockTodoRepository)

	testRepo.On("GetByID", 1).Return(Todo{ID: 1, Title: "Test Todo"}, nil)

	req := httptest.NewRequest("GET", "/todos/1", nil)

	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "1"}}
	c.Request = req

	handler := NewHandler(testRepo)
	handler.GetByID(c)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "Test Todo")
}

func TestGetByID_NotFound(t *testing.T) {
	testRepo := new(MockTodoRepository)

	testRepo.On("GetByID", 2).Return(Todo{}, ErrNotFound)

	req := httptest.NewRequest("GET", "/todos/2", nil)

	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "2"}}
	c.Request = req

	handler := NewHandler(testRepo)
	handler.GetByID(c)

	assert.Equal(t, 404, w.Code)
	assert.Contains(t, w.Body.String(), "not found")
}

func TestGetByID_InvalidID(t *testing.T) {
	// no need to set up a mock repo since error is returned from the handler before calling the repo
	req := httptest.NewRequest("GET", "/todos/abc", nil)

	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler := NewHandler(new(MockTodoRepository))
	handler.GetByID(c)

	assert.Equal(t, 400, w.Code)
	assert.Contains(t, w.Body.String(), "invalid id")
}

// Get /todos/category/:category
func TestGetTodosByCategory_Success(t *testing.T) {
	testRepo := new(MockTodoRepository)

	testRepo.On("GetTodosByCategory", "General").Return([]Todo{{ID: 1, Title: "Test Todo", Category: "General"}}, nil)

	req := httptest.NewRequest("GET", "/todos/category/General", nil)

	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "category", Value: "General"}}
	c.Request = req

	handler := NewHandler(testRepo)
	handler.GetTodosByCategory(c)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "Test Todo")
}

func TestGetTodosByCategory_EmptyArray(t *testing.T) {
	testRepo := new(MockTodoRepository)

	testRepo.On("GetTodosByCategory", "School").Return([]Todo{}, nil)

	req := httptest.NewRequest("GET", "/todos/category/School", nil)

	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "category", Value: "School"}}
	c.Request = req

	handler := NewHandler(testRepo)
	handler.GetTodosByCategory(c)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, w.Body.String(), "[]")
}

// GET /todos/status/:status
func TestGetTodosByStatus_Success(t *testing.T) {
	testRepo := new(MockTodoRepository)

	testRepo.On("GetTodosByStatus", false).Return([]Todo{{ID: 1, Title: "Test Todo", Completed: false}}, nil)

	req := httptest.NewRequest("GET", "/todos/status/false", nil)

	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "status", Value: "false"}}
	c.Request = req

	handler := NewHandler(testRepo)
	handler.GetTodosByStatus(c)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "Test Todo")
}

func TestGetTodosByStatus_InvalidStatus(t *testing.T) {
	req := httptest.NewRequest("GET", "/todos/status/x", nil)

	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "status", Value: "x"}}
	c.Request = req

	handler := NewHandler(new(MockTodoRepository))
	handler.GetTodosByStatus(c)

	assert.Equal(t, 400, w.Code)
	assert.Contains(t, w.Body.String(), "invalid status")
}

// Get /todos/search?q=...
func TestSearch_Success(t *testing.T) {
	testRepo := new(MockTodoRepository)

	testRepo.On("Search", "test").Return([]Todo{{ID: 1, Title: "test something"}}, nil)

	req := httptest.NewRequest("GET", "/todos/search?q=test", nil)

	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler := NewHandler(testRepo)
	handler.Search(c)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "test something")
}

func TestSearch_EmptyArray(t *testing.T) {
	testRepo := new(MockTodoRepository)

	testRepo.On("Search", "nonexistent").Return([]Todo{}, nil)

	req := httptest.NewRequest("GET", "/todos/search?q=nonexistent", nil)

	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler := NewHandler(testRepo)
	handler.Search(c)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, w.Body.String(), "[]")
}
