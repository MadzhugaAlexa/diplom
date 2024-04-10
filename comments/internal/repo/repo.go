package repo

import (
	"comments/internal/entities"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo struct {
	db *pgxpool.Pool
}

func NewRepo(db *pgxpool.Pool) Repo {
	return Repo{
		db: db,
	}
}

func (r *Repo) CreateComment(c *entities.Comment) error {
	tx, err := r.db.Begin(context.Background())
	if err != nil {
		return err
	}
	defer func() {
		tx.Commit(context.Background())
	}()

	sql := `INSERT INTO comments(post_id, parent_id, content, status) values($1, $2, $3, $4) returning id`

	post := tx.QueryRow(
		context.Background(),
		sql,
		c.PostID, c.ParentID, c.Content, c.Status,
	)

	post.Scan(&c.ID)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil

}

func (r *Repo) GetComments(postId int) ([]entities.Comment, error) {
	comments := make([]entities.Comment, 0)

	sql := "select id, post_id, parent_id, content, status from comments where post_id = $1"
	rows, err := r.db.Query(context.Background(), sql, postId)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		comment := entities.Comment{}

		err = rows.Scan(
			&comment.ID, &comment.PostID, &comment.ParentID, &comment.Content, &comment.Status,
		)
		if err != nil {
			return nil, err
		}

		comments = append(comments, comment)
	}
	return comments, nil
}
