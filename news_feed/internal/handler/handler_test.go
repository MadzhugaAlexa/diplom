package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"news_feed/internal/entities"
	"testing"

	"github.com/labstack/echo/v4"
)

func TestNewHandler(t *testing.T) {
	r := TestHandlerReader{
		Posts: []entities.Post{
			{
				Title:   "test 1",
				Content: "content",
			},
			{
				Title:   "test 2",
				Content: "content",
			},
		},
	}

	h := NewHandler(&r)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.GetItems(c)
	if err != nil {
		t.Fatalf("Ошибка: %v\n", err)
	}

	if rec.Code != http.StatusOK {
		t.Fatalf("Плохой ответ: %v\n", rec.Code)
	}

	body := rec.Body
	data := []map[string]interface{}{}
	json.NewDecoder(body).Decode(&data)

	if len(data) != 2 {
		t.Fatalf("Ожидаем 2 элемента, получили: %v\n", len(data))
	}

	item1 := data[0]
	item2 := data[1]

	if item1["Title"] != "test 1" {
		t.Fatal("Неправильный первый Title")
	}
	if item1["Content"] != "content" {
		t.Fatal("Неправильный первый Сontent")
	}
	if item2["Title"] != "test 2" {
		t.Fatal("Неправильный второй Title")
	}
	if item2["Content"] != "content" {
		t.Fatal("Неправильный второй Сontent")
	}

}

type TestHandlerReader struct {
	Posts []entities.Post
}

func (t *TestHandlerReader) ReadItems(int, int, string) ([]entities.Post, error) {
	return t.Posts, nil
}

func (t *TestHandlerReader) ReadItem(int) (entities.Post, error) {
	return entities.Post{}, nil
}
