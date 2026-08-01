package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	repo *TodoRepository
}

func (h *Handler) GetAll(c *gin.Context) {
	todos, err := h.repo.GetAll()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, todos)
}

func (h *Handler) Create(*gin.Context) {
	// TODO
}

func (h *Handler) GetByID(*gin.Context) {
	// TODO
}

func (h *Handler) Update(*gin.Context) {
	// TODO
}

func (h *Handler) Delete(*gin.Context) {
	// TODO
}
