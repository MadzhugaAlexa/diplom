package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"

	"github.com/labstack/echo/v4"
)

type Comment struct {
	Content string
}

type Check struct {
	Valid bool
}

func CheckComment(c echo.Context) error {
	body := c.Request().Body
	comment := Comment{}
	err := json.NewDecoder(body).Decode(&comment)
	if err != nil {
		log.Printf("ошибка: %#v\n", err)
		return err
	}

	if HasBadWords(comment.Content) {
		return c.String(http.StatusBadRequest, "")
	} else {
		return c.String(http.StatusOK, "")
	}
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
