package main

import (
	"context"
	"log"
	"news_feed/internal/config"
	"news_feed/internal/feed"
	"news_feed/internal/handler"
	"news_feed/internal/repo"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

func main() {
	DB_URL := os.Getenv("DB")
	if DB_URL == "" {
		DB_URL = "postgres://alexa:alexa@localhost:5432/rss"
	}

	db, err := pgxpool.New(context.Background(), DB_URL)
	if err != nil {
		log.Fatalf("Не смогли подключиться к БД: %v\n", err)
	}
	defer db.Close()

	repo := repo.NewRepo(db)
	cfg := config.LoadConfig("./config.json")
	feed.FeatchFeeds(cfg, feed.LoadFeed, repo)

	e := echo.New()

	e.Use(handler.Logger)
	h := handler.NewHandler(repo)
	e.GET("/news", h.GetItems)
	e.GET("/news/:id", h.GetItem)
	e.Logger.Fatal(e.Start(":1111"))
}
