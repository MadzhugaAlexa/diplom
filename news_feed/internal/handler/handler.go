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
	ReadItems(int) ([]entities.Post, error)
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
