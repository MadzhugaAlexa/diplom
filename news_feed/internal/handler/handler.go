package handler

import (
	"errors"
	"log"
	"net/http"
	"news_feed/internal/entities"
	"news_feed/internal/repo"
	"strconv"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	repo ItemsReader
}

type ItemsReader interface {
	ReadItems(perPage int, page int, s string) ([]entities.Post, error)
	ReadItem(int) (entities.Post, error)
}

func NewHandler(r ItemsReader) Handler {
	return Handler{
		repo: r,
	}
}

func (h *Handler) GetItem(c echo.Context) error {
	strId := c.Param("id")
	var id int
	var err error
	id, err = strconv.Atoi(strId)
	if err != nil {
		return err
	}
	item, err := h.repo.ReadItem(id)
	if err != nil {
		if errors.Is(err, repo.ErrNoRows) {
			return c.JSON(http.StatusNotFound, "not found")
		}

		log.Printf("Ошибка %#v\n", err)

		return err
	}

	return c.JSON(http.StatusOK, item)
}

// GetItems обрабатывает запрос на /news и возвращает N последних новостей
// По умолчанию N = 10
func (h *Handler) GetItems(c echo.Context) error {
	s := c.QueryParam("s")

	perPageStr := c.QueryParam("per_page")
	perPage, err := ParseStrQueryParam(perPageStr, 10)
	if err != nil {
		return err
	}

	pageStr := c.QueryParam("page")
	page, err := ParseStrQueryParam(pageStr, 1)
	if err != nil {
		return err
	}

	items, err := h.repo.ReadItems(perPage, page, s)
	if err != nil {
		log.Printf("Ошибка %#v\n", err)

		return err
	}

	return c.JSON(http.StatusOK, items)
}

func ParseStrQueryParam(in string, d int) (int, error) {
	if in == "" {
		return d, nil
	}
	val, err := strconv.Atoi(in)

	if err != nil {
		return 0, err
	}

	return val, nil
}
