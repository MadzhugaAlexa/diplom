package main

import (
	"api-gateway/internal/handlers"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	e.Use(handlers.Logger)
	e.GET("/news", handlers.GetAllNews)
	e.GET("/news/:id", handlers.GetOneNew)
	e.POST("/comments/", handlers.AddComment)
	e.Start(":8080")
}
