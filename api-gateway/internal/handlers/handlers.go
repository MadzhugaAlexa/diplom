package handlers

import (
	"api-gateway/internal/entities"
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

var RSS_PORT string = "1111"
var COMMENTS_PORT string = "1112"

func GetAllNews(c echo.Context) error {
	url := "http://localhost:" + RSS_PORT + "/news"
	ip := GetIP(c)

	if qs := c.QueryString(); qs != "" {
		url = url + "?" + qs
	}

	client := http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("X-Forwarded-For", ip)
	resp, err := client.Do(req)
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

func GetIP(c echo.Context) string {
	ips := c.Request().Header["X-Forwarded-For"]
	var ip string

	if len(ips) > 0 {
		ip = ips[0]
	}

	return ip
}

func GetOneNew(c echo.Context) error {
	id := c.Param("id")
	ip := GetIP(c)

	requestID := c.QueryParam("request_id")

	rssResponse := RssResponse{}
	commentsResponse := CommentsResponse{}

	postCh := make(chan RssResponse)
	go LoadPost(postCh, id, requestID, ip)

	commentsCh := make(chan CommentsResponse)
	go LoadPostComments(commentsCh, id, requestID, ip)

	for i := 0; i < 2; i++ {
		select {
		case resp := <-postCh:
			rssResponse = resp
			if rssResponse.Error != nil {
				return rssResponse.Error
			}
		case resp := <-commentsCh:
			commentsResponse = resp
			if commentsResponse.Error != nil {
				return commentsResponse.Error
			}
		}
	}

	if rssResponse.Error != nil {
		return rssResponse.Error
	}

	if commentsResponse.Error != nil {
		return commentsResponse.Error
	}

	post := rssResponse.Post
	post.Comments = commentsResponse.Comments

	return c.JSON(http.StatusOK, post)
}

func LoadPost(out chan RssResponse, id string, requestID string, ip string) {
	result := RssResponse{}
	url := "http://localhost:" + RSS_PORT + "/news/" + id + "?request_id=" + requestID

	client := http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		result.Error = err
		out <- result
		return
	}

	req.Header.Set("X-Forwarded-For", ip)
	resp, err := client.Do(req)

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

	out <- result
}

func LoadPostComments(out chan CommentsResponse, id string, requestID string, ip string) {
	response := CommentsResponse{}

	url := "http://localhost:" + COMMENTS_PORT + "/comments/" + id + "?request_id=" + requestID

	client := http.Client{}
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		response.Error = err
	}

	req.Header.Set("X-Forwarded-For", ip)
	resp, err := client.Do(req)
	if err != nil {
		return
	}

	br, err := io.ReadAll(resp.Body)
	if err != nil {
		response.Error = err
	}
	err = json.Unmarshal(br, &response.Comments)
	if err != nil {
		response.Error = err
	}
	out <- response
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

	client := http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewReader(body))

	if err != nil {
		return err
	}

	ip := GetIP(c)
	req.Header.Set("X-Forwarded-For", ip)
	resp, err := client.Do(req)
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

	client := http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(censJSON))

	if err != nil {
		return false, err
	}

	ip := GetIP(c)
	req.Header.Set("X-Forwarded-For", ip)
	resp, err := client.Do(req)

	if err != nil {
		return false, err
	}

	if resp.StatusCode == 200 {
		return true, nil
	}

	return false, nil
}
