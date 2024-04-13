package main

import (
	"api-gateway/internal/handlers"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	e.GET("/news", handlers.GetAllNews)
	e.GET("/news/:id", handlers.GetOneNew)
	// e.GET("/comments/:post_id", h.GetComments)
	e.POST("/comments/", handlers.AddComment)
	e.Start(":8080")
}
