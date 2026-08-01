package main

import (
	"net/http"
	"strconv"

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
	// TODO
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

func (h *Handler) Update(c *gin.Context) {
	// TODO
}

func (h *Handler) Delete(c *gin.Context) {
	// TODO
}

func parseID(c *gin.Context) (int, error) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return 0, err
	}
	return id, nil
}
