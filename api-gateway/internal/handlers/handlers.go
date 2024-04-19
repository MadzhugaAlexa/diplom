package handlers

import (
	"api-gateway/internal/entities"
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"sync"

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

type RssResponse struct {
	Post  entities.NewsFullDetailed
	Error error
}

type CommentsResponse struct {
	Comments []entities.Comment
	Error    error
}

func GetOneNew(c echo.Context) error {
	var wg sync.WaitGroup
	wg.Add(2)
	id := c.Param("id")
	requestID := c.QueryParam("request_id")

	var rssResponse RssResponse
	go LoadPost(&wg, &rssResponse, id, requestID)

	if rssResponse.Error != nil {
		return rssResponse.Error
	}

	commentsResponse := CommentsResponse{}
	go LoadPostComments(&wg, &commentsResponse, id, requestID)

	if commentsResponse.Error != nil {
		return commentsResponse.Error
	}

	wg.Wait()
	post := rssResponse.Post
	post.Comments = commentsResponse.Comments

	return c.JSON(http.StatusOK, post)
}

func LoadPost(wg *sync.WaitGroup, result *RssResponse, id string, requestID string) {
	defer wg.Done()

	resp, err := http.Get("http://localhost:" + RSS_PORT + "/news/" + id + "?request_id=" + requestID)
	if err != nil {
		result.Error = err
	}

	br, err := io.ReadAll(resp.Body)
	if err != nil {
		result.Error = err
	}

	err = json.Unmarshal(br, &result.Post)

	if err != nil {
		result.Error = err
	}
}

func LoadPostComments(wg *sync.WaitGroup, response *CommentsResponse, id string, requestID string) {
	defer wg.Done()

	resp, err := http.Get("http://localhost:" + COMMENTS_PORT + "/comments/" + id + "?request_id=" + requestID)
	if err != nil {
		response.Error = err
	}

	br, err := io.ReadAll(resp.Body)
	if err != nil {
		response.Error = err
	}
	err = json.Unmarshal(br, &response.Comments)
	if err != nil {
		response.Error = err
	}
}
func AddComment(c echo.Context) error {
	requestID := c.QueryParam("request_id")

	comment := entities.Comment{}

	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return err
	}
	br := bytes.NewReader(body)

	err = json.NewDecoder(br).Decode(&comment)
	if err != nil {
		log.Printf("ошибка: %#v\n", err)
		return err
	}

	valid, err := ValidateComment(c, comment)
	if err != nil {
		return err
	}

	if !valid {
		return c.String(http.StatusBadRequest, "")
	}

	url := "http://localhost:" + COMMENTS_PORT + "/comments/?request_id=" + requestID

	resp, err := http.Post(
		url,
		"application/json",
		bytes.NewReader(body),
	)

	if err != nil {
		return err
	}

	response, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return c.String(http.StatusOK, string(response))
}

type Cens struct {
	Content string
}

const CENS_PORT = "1113"

func ValidateComment(c echo.Context, comment entities.Comment) (bool, error) {
	cens := Cens{
		Content: comment.Content,
	}
	censJSON, err := json.Marshal(cens)
	if err != nil {
		return false, err
	}
	requestID := c.QueryParam("request_id")
	url := "http://localhost:" + CENS_PORT + "/check_comment?request_id=" + requestID

	resp, err := http.Post(
		url,
		"application/json",
		bytes.NewBuffer(censJSON),
	)
	if err != nil {
		return false, err
	}

	if resp.StatusCode == 200 {
		return true, nil
	}

	return false, nil
}
