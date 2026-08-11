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

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	defer sqlDB.Close()

	repo := NewTodoRepository(db)
	userRepo := NewUserRepository(db)
	authHandler := NewAuthHandler(userRepo)

	handler := NewHandler(repo)

	router := gin.Default()

	router.POST("/signup", authHandler.SignUp)
	router.POST("/login", authHandler.LogIn)

	todos := router.Group("/todos")
	todos.Use(AuthMiddleware())
	{
		todos.GET("", handler.GetAll)
		todos.GET("/category/:category", handler.GetTodosByCategory)
		todos.GET("/status/:status", handler.GetTodosByStatus)
		todos.GET("/search", handler.Search)
		todos.GET("/:id", handler.GetByID)
		todos.POST("", handler.Create)
		todos.PUT("/category/:category", handler.UpdateByCategory)
		todos.PUT("/:id", handler.Update)
		todos.DELETE("", AdminOnly(), handler.DeleteAll)
		todos.DELETE("/:id", handler.Delete)
	}

	router.Run()
}
