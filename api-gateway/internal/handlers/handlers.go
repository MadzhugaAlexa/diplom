package handlers

import (
	"api-gateway/internal/entities"
	"encoding/json"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
)

var RSS_PORT string = "1111"
var COMMENTS_PORT string = "1112"

func GetAllNews(c echo.Context) error {
	url := "http://localhost:" + RSS_PORT + "/news"

	if qs := c.QueryString(); qs != "" {
		url = url + "?" + qs
	}

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	br, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var items []entities.NewsShortDetailed
	err = json.Unmarshal(br, &items)

	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, items)
}

func GetOneNew(c echo.Context) error {
	id := c.Param("id")
	requestID := c.QueryParam("request_id")

	resp, err := http.Get("http://localhost:" + RSS_PORT + "/news/" + id + "?request_id=" + requestID)
	if err != nil {
		return err
	}

	br, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var post entities.NewsFullDetailed
	err = json.Unmarshal(br, &post)

	if err != nil {
		return err
	}

	resp, err = http.Get("http://localhost:" + COMMENTS_PORT + "/comments/" + id + "?request_id=" + requestID)
	if err != nil {
		return err
	}

	br, err = io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(br, &post.Comments)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, post)
}

func AddComment(c echo.Context) error {
	requestID := c.QueryParam("request_id")

	url := "http://localhost:" + COMMENTS_PORT + "/comments/?request_id=" + requestID
	resp, err := http.Post(
		url,
		"application/json",
		c.Request().Body,
	)

	if err != nil {
		return err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return c.String(http.StatusOK, string(body))
}
