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
		router.GET("", handler.GetAll)
		router.GET("/category/:category", handler.GetTodosByCategory)
		router.GET("/status/:status", handler.GetTodosByStatus)
		router.GET("/search", handler.Search)
		router.GET("/:id", handler.GetByID)
		router.POST("", handler.Create)
		router.PUT("/category/:category", handler.UpdateByCategory)
		router.PUT("/:id", handler.Update)
		router.DELETE("", AdminOnly(), handler.DeleteAll)
		router.DELETE("/:id", handler.Delete)
	}

	router.Run()
}
