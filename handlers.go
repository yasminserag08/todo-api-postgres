package main

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	repo *TodoRepository
}

func NewHandler(repo *TodoRepository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) GetAll(c *gin.Context) {
	todos, err := h.repo.GetAll()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, todos)
}

func (h *Handler) Create(c *gin.Context) {
	var newTodo Todo
	err := c.ShouldBindJSON(&newTodo)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if newTodo.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title is required"})
		return
	}

	switch strings.ToLower(newTodo.Priority) {
	case "low":
		newTodo.Priority = "Low"
	case "medium":
		newTodo.Priority = "Medium"
	case "high":
		newTodo.Priority = "High"
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "priority must be Low, Medium, or High"})
		return
	}
	newTodo, err = h.repo.Create(newTodo)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, newTodo)
}

func (h *Handler) GetByID(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		return // parseID already returns 400 response
	}

	todo, err := h.repo.GetByID(id)

	if err == ErrNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, todo)
}

func (h *Handler) GetTodosByCategory(c *gin.Context) {
	category := c.Param("category")

	todos, err := h.repo.GetTodosByCategory(category)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, todos)
}

func (h *Handler) GetTodosByStatus(c *gin.Context) {
	status, err := parseStatus(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	todos, err := h.repo.GetTodosByStatus(status)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, todos)
}

func (h *Handler) Search(c *gin.Context) {
	q := c.Query("q")

	todos, err := h.repo.Search(q)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, todos)
}

func (h *Handler) Update(c *gin.Context) {
	id, err := parseID(c)

	if err != nil {
		return
	}

	var newTodo Todo

	err = c.ShouldBindJSON(&newTodo)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if newTodo.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title is required"})
		return
	}

	switch strings.ToLower(newTodo.Priority) {
	case "low":
		newTodo.Priority = "Low"
	case "medium":
		newTodo.Priority = "Medium"
	case "high":
		newTodo.Priority = "High"
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "priority must be Low, Medium, or High"})
		return
	}

	newTodo, err = h.repo.Update(id, newTodo.Title, newTodo.Completed)

	if err == ErrNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, newTodo)
}

func (h *Handler) UpdateByCategory(c *gin.Context) {
	category := c.Param("category")

	var newTodo Todo

	err := c.ShouldBindJSON(&newTodo)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// only take completed since the spec says so
	newTodos, err := h.repo.UpdateByCategory(category, newTodo.Completed)

	if err == ErrNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, newTodos)
}

func (h *Handler) Delete(c *gin.Context) {
	id, err := parseID(c)

	if err != nil {
		return
	}

	err = h.repo.Delete(id)

	if err == ErrNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Todo deleted"})
}

func (h *Handler) DeleteAll(c *gin.Context) {
	err := h.repo.DeleteAll()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "All todos deleted"})
}

func parseID(c *gin.Context) (int, error) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return 0, err
	}
	return id, nil
}

func parseStatus(c *gin.Context) (bool, error) {
	status, err := strconv.ParseBool(c.Param("status"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status"})
		return false, err
	}
	return status, nil
}
