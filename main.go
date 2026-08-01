package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	db, err := connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	repo := NewTodoRepository(db) // to be passed to the handlers

	handler := NewHandler(repo)

	router := gin.Default()

	router.GET("/todos", handler.GetAll)

	router.GET("/todos/:id", handler.GetByID)

	router.POST("/todos", handler.Create)

	router.Run()
}
