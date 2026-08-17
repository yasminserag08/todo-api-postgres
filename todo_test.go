package main

import (
	"errors"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

// POST /todos
func TestCreate_Success(t *testing.T) {
	testRepo := new(MockTodoRepository)

	testRepo.On("Create", Todo{Title: "Test Todo", Priority: "High", UserID: 1}).Return(Todo{ID: 1, Title: "Test Todo"}, nil)

	body := strings.NewReader(`{"title": "Test Todo", "priority":"high", "userID": 1}`)
	req := httptest.NewRequest("POST", "/todos", body)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Set("userID", uint(1))
	c.Request = req

	handler := NewHandler(testRepo)
	handler.Create(c)

	assert.Equal(t, 201, w.Code)
	assert.Contains(t, w.Body.String(), "Test Todo")
}

func TestCreate_EmptyTitle(t *testing.T) {
	body := strings.NewReader(`{"title": "", "priority":"high"}`)
	req := httptest.NewRequest("POST", "/todos", body)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Set("userID", uint(1))
	c.Request = req

	handler := NewHandler(new(MockTodoRepository))
	handler.Create(c)

	assert.Equal(t, 400, w.Code)
	assert.Contains(t, w.Body.String(), "title is required")
}

func TestCreate_InvalidPriority(t *testing.T) {
	body := strings.NewReader(`{"title": "hello", "priority":"x"}`)
	req := httptest.NewRequest("POST", "/todos", body)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Set("userID", uint(1))
	c.Request = req

	handler := NewHandler(new(MockTodoRepository))
	handler.Create(c)

	assert.Equal(t, 400, w.Code)
	assert.Contains(t, w.Body.String(), "priority must be Low, Medium, or High")
}

func TestCreate_PastDueDate(t *testing.T) {
	body := strings.NewReader(`{"title": "test", "priority": "High", "dueDate": "2020-01-01T00:00:00Z"}`)
	req := httptest.NewRequest("POST", "/todos", body)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Set("userID", uint(1))
	c.Request = req

	handler := NewHandler(new(MockTodoRepository))
	handler.Create(c)

	assert.Equal(t, 400, w.Code)
	assert.Contains(t, w.Body.String(), "due date must be in the future")
}

func TestCreate_InvalidJSON(t *testing.T) {
	body := strings.NewReader(`{invalid json}`)
	req := httptest.NewRequest("POST", "/todos", body)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Set("userID", uint(1))
	c.Request = req

	handler := NewHandler(new(MockTodoRepository))
	handler.Create(c)

	assert.Equal(t, 400, w.Code)
	assert.Contains(t, w.Body.String(), "error")
}

// PUT /todos/:id
func TestUpdate_Success(t *testing.T) {
	testRepo := new(MockTodoRepository)

	testRepo.On("Update", 1, Todo{Title: "Test", Priority: "High", Category: "General"}).Return(Todo{ID: 1, Title: "Test", Priority: "High", Category: "General"}, nil)

	body := strings.NewReader(`{"title": "Test", "priority":"high", "category":"General"}`)
	req := httptest.NewRequest("PUT", "/todos/1", body)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "1"}}
	c.Request = req

	handler := NewHandler(testRepo)
	handler.Update(c)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "Test")
}

func TestUpdate_NotFound(t *testing.T) {
	testRepo := new(MockTodoRepository)

	testRepo.On("Update", 2, Todo{Title: "Test", Priority: "High", Category: "General"}).Return(Todo{}, ErrNotFound)

	body := strings.NewReader(`{"title": "Test", "priority":"high", "category":"General"}`)
	req := httptest.NewRequest("PUT", "/todos/2", body)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "2"}}
	c.Request = req

	handler := NewHandler(testRepo)
	handler.Update(c)

	assert.Equal(t, 404, w.Code)
	assert.Contains(t, w.Body.String(), "not found")
}

func TestUpdate_EmptyTitle(t *testing.T) {
	body := strings.NewReader(`{"title": "", "priority":"high"}`)
	req := httptest.NewRequest("PUT", "/todos/2", body)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "2"}}
	c.Request = req

	handler := NewHandler(new(MockTodoRepository))
	handler.Update(c)

	assert.Equal(t, 400, w.Code)
	assert.Contains(t, w.Body.String(), "title is required")
}

func TestUpdate_InvalidPriority(t *testing.T) {
	body := strings.NewReader(`{"title": "hello", "priority":"x"}`)
	req := httptest.NewRequest("PUT", "/todos/2", body)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "2"}}
	c.Request = req

	handler := NewHandler(new(MockTodoRepository))
	handler.Update(c)

	assert.Equal(t, 400, w.Code)
	assert.Contains(t, w.Body.String(), "priority must be Low, Medium, or High")
}

func TestUpdate_PastDueDate(t *testing.T) {
	body := strings.NewReader(`{"title": "test", "priority": "High", "dueDate": "2020-01-01T00:00:00Z"}`)
	req := httptest.NewRequest("PUT", "/todos/2", body)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "2"}}
	c.Request = req

	handler := NewHandler(new(MockTodoRepository))
	handler.Update(c)

	assert.Equal(t, 400, w.Code)
	assert.Contains(t, w.Body.String(), "due date must be in the future")
}

// PUT /todos/category/:category
func TestUpdateByCategory_Success(t *testing.T) {
	testRepo := new(MockTodoRepository)

	testRepo.On("UpdateByCategory", "General", true, mock.AnythingOfType("*time.Time")).
		Return([]Todo{{ID: 1, Title: "Test Todo", Category: "General", Completed: true, Priority: "High"}}, nil)

	body := strings.NewReader(`{"title": "Test Todo", "priority": "High", "completed": true}`)
	req := httptest.NewRequest("PUT", "/todos/category/General", body)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "category", Value: "General"}}
	c.Request = req

	handler := NewHandler(testRepo)
	handler.UpdateByCategory(c)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "Test Todo")
}

func TestUpdateByCategory_CategoryNotFound(t *testing.T) {
	testRepo := new(MockTodoRepository)

	testRepo.On("UpdateByCategory", "School", true, mock.AnythingOfType("*time.Time")).
		Return([]Todo{}, ErrNotFound)

	body := strings.NewReader(`{"title": "Test Todo", "priority": "High", "completed": true}`)
	req := httptest.NewRequest("PUT", "/todos/category/School", body)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "category", Value: "School"}}
	c.Request = req

	handler := NewHandler(testRepo)
	handler.UpdateByCategory(c)

	assert.Equal(t, 404, w.Code)
	assert.Contains(t, w.Body.String(), "not found")
}

// delete test cases have been removed because they are tested in auth_test.go
