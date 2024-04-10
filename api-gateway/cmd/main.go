package main

import (
	"api-gateway/internal/entities"
	"encoding/json"
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
	br, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var items []entities.NewsFullDetailed
	err = json.Unmarshal(br, &items)

	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, items)

}

func GetOneNew(c echo.Context) error {
	id := c.Param("id")

	resp, err := http.Get("http://localhost:" + RSS_PORT + "/news/" + id)
	if err != nil {
		return err
	}

	br, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var item entities.NewsFullDetailed
	err = json.Unmarshal(br, &item)

	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, item)
}
