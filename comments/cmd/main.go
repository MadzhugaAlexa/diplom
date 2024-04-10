package main

import (
	"comments/internal/handlers"
	"comments/internal/repo"
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

func main() {
	DB_URL := "postgres://alexa:alexa@localhost:5432/comments"
	db, err := pgxpool.New(context.Background(), DB_URL)

	if err != nil {
		log.Fatalf("Не смогли подключиться к БД: %v\n", err)
	}
	defer db.Close()
	repo := repo.NewRepo(db)

	// comment := entities.Comment{
	// 	PostID:  1,
	// 	Content: "content",
	// 	Status:  "status",
	// }

	// err = repo.CreateComment(&comment)
	// log.Printf("comment: %#v\n", comment)

	// if err != nil {
	// 	log.Fatal(err)
	// }

	// comments, err := repo.GetComments(1)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// log.Printf("comments: %#v\n", comments)
	e := echo.New()
	h := handlers.NewHandler(&repo)
	e.GET("/comments/:post_id", h.GetComments)
	e.POST("/comments/", h.AddComment)
	e.Logger.Fatal(e.Start(":1112"))
}
