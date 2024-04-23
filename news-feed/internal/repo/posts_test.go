package repo

import (
	"context"
	"log"
	"news_feed/internal/entities"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

func newTestDB() *pgxpool.Pool {
	DB_URL := "postgres://alexa:alexa@localhost:5432/rss_test"

	db, err := pgxpool.New(context.Background(), DB_URL)
	if err != nil {
		log.Fatalf("Не смогли подключиться к БД: %v\n", err)
	}

	_, err = db.Exec(context.Background(), "delete from posts;")
	if err != nil {
		log.Fatalf("Ошибка удаления из БД")
	}

	return db
}

func TestAddItem(t *testing.T) {
	db := newTestDB()
	defer db.Close()

	repo := NewRepo(db)

	item := entities.Item{
		Title:   "test title",
		Link:    "http://www.ru",
		PubDate: "Mon, 2 Jan 2006 15:04:05 -0700",
	}

	err := repo.AddItem(item)
	if err != nil {
		t.Fatal("Не смог создать новость с ошибкой:", err)
	}

	row := db.QueryRow(context.Background(), "select exists (select 1 from posts where title = 'test title');")
	exists := false
	err = row.Scan(&exists)
	if err != nil {
		t.Fatal("Ошибка сканирования: ", err)
	}

	if !exists {
		t.Fatal("Не изменилось")
	}

	err = repo.AddItem(item)
	if err != nil {
		t.Fatal("Не смог создать новость с ошибкой:", err)
	}

	row = db.QueryRow(context.Background(), "select count(*) from posts where title = 'test title'")
	cnt := 0
	if row.Scan(&cnt); cnt != 1 {
		t.Fatal("Существует дубликат новости")
	}
}

func TestReadItems(t *testing.T) {
	db := newTestDB()
	defer db.Close()

	repo := NewRepo(db)

	item := entities.Item{
		Title:   "test title",
		Link:    "http://www.ru",
		PubDate: "Mon, 2 Jan 2006 15:04:05 -0700",
	}

	_ = repo.AddItem(item)

	items, err := repo.ReadItems(10, 1, "")
	if err != nil {
		t.Fatal("Не прочитали новости", err)
	}

	if len(items) != 1 {
		t.Fatal("Прочитали не одну новость", len(items))
	}
}
