package repo

import (
	"context"
	"errors"
	"fmt"
	"log"
	"news_feed/internal/entities"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo struct {
	db *pgxpool.Pool
}

var ErrNoRows = errors.New("no rows in result set")

func NewRepo(db *pgxpool.Pool) *Repo {
	return &Repo{
		db: db,
	}
}

var errFailedToSave = errors.New("не удалось сохранить")

// AddItem проверяет, есть ли указанная новость в БД и если нет -
// создает ее
func (r *Repo) AddItem(item entities.Item) error {
	tx, err := r.db.Begin(context.Background())
	if err != nil {
		return err
	}
	defer func() {
		tx.Commit(context.Background())
	}()

	row := tx.QueryRow(
		context.Background(),
		"select exists(select 1 from posts where link = $1);",
		item.Link,
	)

	var exists bool
	err = row.Scan(&exists)
	if err != nil {
		fmt.Printf("Ошибка: %#v\n", err)
	}

	if exists {
		return nil
	}

	sql := `INSERT INTO posts(title, link, content, pubDate) 
	values($1, $2, $3, $4) `

	putDate, err := time.Parse(time.RFC1123, item.PubDate)
	if err != nil {
		putDate, err = time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			putDate, err = time.Parse("Mon, 2 Jan 2006 15:04:05 -0700", item.PubDate)
			if err != nil {
				log.Fatalf("Ошибка даты %#v\n", item)
			}

		}
	}

	t, err := tx.Exec(
		context.Background(),
		sql,
		item.Title, item.Link, item.Description, putDate.Unix(),
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	if t.RowsAffected() != 1 {
		return errFailedToSave
	}

	return nil
}

func (r *Repo) ReadItem(id int) (entities.Post, error) {
	item := entities.Post{}
	sql := "select id, title, link, content, pubdate from posts where id=$1"
	row := r.db.QueryRow(context.Background(), sql, id)

	err := row.Scan(
		&item.ID, &item.Title, &item.Link, &item.Content, &item.PubDate,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return item, ErrNoRows
		}
		return item, err
	}

	return item, nil
}

// ReadItems читает новостные посты из БД и возвращет их слайс.
func (r *Repo) ReadItems(limit int, s string) ([]entities.Post, error) {
	items := make([]entities.Post, 0)

	sql := "select id, title, link, content, pubdate from posts"

	if s != "" {
		sql = sql + " where title like '%" + s + "%'"
	}
	sql = sql + " order by pubdate desc limit $1"

	rows, err := r.db.Query(context.Background(), sql, limit)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		item := entities.Post{}

		err = rows.Scan(
			&item.ID, &item.Title, &item.Link, &item.Content, &item.PubDate,
		)
		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}
	return items, nil
}
