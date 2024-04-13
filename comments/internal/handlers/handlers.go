package handlers

import (
	"comments/internal/entities"
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	repo CommentsRepo
}

type CommentsRepo interface {
	GetComments(int) ([]entities.Comment, error)
	CreateComment(*entities.Comment) error
	UpdateStatus(*entities.Comment) error
}

func NewHandler(r CommentsRepo) Handler {
	return Handler{
		repo: r,
	}
}
func (h *Handler) GetComments(c echo.Context) error {
	i := c.Param("post_id")

	var id int
	var err error

	id, err = strconv.Atoi(i)
	if err != nil {
		log.Printf("ошибка: %v\n", err)
		return err
	}

	comment, err := h.repo.GetComments(id)
	if err != nil {
		log.Printf("ошибка %#v\n", err)
		return err
	}

	return c.JSON(http.StatusOK, comment)
}

func (h *Handler) AddComment(c echo.Context) error {
	body := c.Request().Body
	comment := entities.Comment{}
	err := json.NewDecoder(body).Decode(&comment)
	if err != nil {
		log.Printf("ошибка: %#v\n", err)
		return err
	}

	err = h.repo.CreateComment(&comment)
	if err != nil {
		log.Printf("ошибка: %v\n", err)
		return err
	}

	go func() {
		if HasBadWords(comment.Content) {
			comment.Status = "bad"
		} else {
			comment.Status = "ready"
		}
		err := h.repo.UpdateStatus(&comment)
		if err != nil {
			log.Printf("failed to update status %v", comment)
		}
	}()

	return c.JSON(http.StatusOK, comment)
}

func HasBadWords(s string) bool {
	badWords := []string{"qwerty", "йцукен", "zxvbnm"}

	for _, bad := range badWords {
		match, _ := regexp.Match(bad, []byte(s))
		if match {
			return true
		}
	}

	return false
}
