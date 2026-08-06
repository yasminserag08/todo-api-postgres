package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	db, err := connect()

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	defer sqlDB.Close()

	repo := NewTodoRepository(db) // to be passed to the handlers

	handler := NewHandler(repo)

	router := gin.Default()

	router.GET("/todos", handler.GetAll)

	router.GET("/todos/:id", handler.GetByID)

	router.GET("/todos/category/:category", handler.GetTodosByCategory)

	router.GET("todos/status/:status", handler.GetTodosByStatus)

	router.POST("/todos", handler.Create)

	router.PUT("/todos/:id", handler.Update)

	router.PUT("/todos/category/:category", handler.UpdateByCategory)

	router.DELETE("/todos/:id", handler.Delete)

	router.DELETE("/todos", handler.DeleteAll)

	router.Run()
}
