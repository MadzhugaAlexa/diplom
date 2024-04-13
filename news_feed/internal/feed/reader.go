package feed

import (
	"encoding/xml"
	"io"
	"log"
	"net/http"
	"news_feed/internal/config"
	"news_feed/internal/entities"
	"time"
)

type FeedReader func(string) ([]byte, error)

// Read загружает RSS-новости для заданного URL-а и пишет их в out канал.
// Если возникли ошибки - они пишутся в ech канал
func read(reader FeedReader, url string, out chan<- entities.Item, ech chan<- entities.Error) {
	body, err := reader(url)

	if err != nil {
		ech <- entities.Error{Error: err}
		return
	}

	rss := entities.RSS{}
	err = xml.Unmarshal(body, &rss)
	if err != nil {
		ech <- entities.Error{Error: err}
		return
	}

	for _, item := range rss.Channel.Items {
		out <- item
	}
}

type RepoInterface interface {
	AddItem(item entities.Item) error
	ReadItems(limit int, s string) ([]entities.Post, error)
}

// Функция FeatchFeeds читает новости из RSS каналов и пишет их в БД.
func FeatchFeeds(cfg config.Config, reader FeedReader, repo RepoInterface) {

	ch := make(chan entities.Item)
	ech := make(chan entities.Error)

	for _, rss := range cfg.RSS {
		ticker := time.NewTicker(time.Duration(cfg.RequestPeriod) * cfg.Duration)
		go func(rss string) {
			for range ticker.C {
				read(reader, rss, ch, ech)
			}
		}(rss)
	}

	go func(ch chan entities.Item) {
		for {
			item := <-ch
			repo.AddItem(item)
		}
	}(ch)

	go func(ech chan entities.Error) {
		for {
			err := <-ech
			log.Printf("Ошибка: %v\n", err)
		}
	}(ech)
}

// Функция LoadFeed отправляет запрос по сети по заданному URL и возвращает ответ или ошибку
func LoadFeed(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		return []byte{}, err
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return []byte{}, err
	}

	return body, nil
}
