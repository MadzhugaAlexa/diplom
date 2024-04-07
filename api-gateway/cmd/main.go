package main

import (
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
)

var RSS_PORT string = "1111"

func main() {
	e := echo.New()

	e.GET("/news/", GetAllNews)
	e.GET("/news/:id", GetOneNew)
	e.Start(":8080")
}
func GetAllNews(c echo.Context) error {
	resp, err := http.Get("http://localhost:" + RSS_PORT + "/news/")
	if err != nil {
		return err
	}
	var br []byte
	br, err = io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return c.String(http.StatusOK, string(br))

}

func GetOneNew(c echo.Context) error {
	id := c.Param("id")

	resp, err := http.Get("http://localhost:" + RSS_PORT + "/news/" + id)
	if err != nil {
		return err
	}
	var br []byte
	br, err = io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return c.String(http.StatusOK, string(br))
}
