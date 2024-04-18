package handlers

import (
	"bytes"
	"comments/internal/entities"
	"encoding/json"
	"log"
	"net/http"
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

const CENS_PORT = "1113"

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

type Cens struct {
	Content string
}

type CensResponse struct {
	Valid bool
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

	go h.ValidateComment(c, comment)

	return c.JSON(http.StatusOK, comment)
}

func (h *Handler) ValidateComment(c echo.Context, comment entities.Comment) {
	cens := Cens{
		Content: comment.Content,
	}
	censJSON, err := json.Marshal(cens)
	if err != nil {
		return
	}
	requestID := c.QueryParam("request_id")
	url := "http://localhost:" + CENS_PORT + "/check_comment?request_id=" + requestID

	resp, err := http.Post(
		url,
		"application/json",
		bytes.NewBuffer(censJSON),
	)
	if err != nil {
		return
	}

	valid := CensResponse{}
	err = json.NewDecoder(resp.Body).Decode(&valid)
	if err != nil {
		log.Printf("ошибка: %#v\n", err)
		return
	}

	if valid.Valid {
		comment.Status = "ready"
	} else {
		comment.Status = "bad"
	}

	err = h.repo.UpdateStatus(&comment)
	if err != nil {
		log.Printf("failed to update status %v", comment)
	}
}
