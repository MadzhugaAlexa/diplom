package handler

import (
	"log"
	"net/http"
	"news_feed/internal/entities"
	"strconv"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	repo ItemsReader
}

type ItemsReader interface {
	ReadItems(int) ([]entities.Post, error)
}

func NewHandler(r ItemsReader) Handler {
	return Handler{
		repo: r,
	}
}

// GetItems обрабатывает запрос на /news и возвращает N последних новостей
// По умолчанию N = 10
func (h *Handler) GetItems(c echo.Context) error {
	l := c.Param("limit")

	var limit int
	var err error
	if l == "" {
		limit = 10
	} else {
		limit, err = strconv.Atoi(l)
		if err != nil {
			return err
		}
	}

	items, err := h.repo.ReadItems(limit)
	if err != nil {
		log.Printf("Ошибка %#v\n", err)

		return err
	}

	return c.JSON(http.StatusOK, items)
}
