package main

import (
	"cens/internal/handlers"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	e.Use(handlers.Logger)
	e.POST("/check_comment", handlers.CheckComment)

	e.Start(":1113")
}
