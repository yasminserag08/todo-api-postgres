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

	var h Handler
	h.repo = repo

	router := gin.Default()

	router.GET("/todos", h.GetAll)

	router.Run()
}
