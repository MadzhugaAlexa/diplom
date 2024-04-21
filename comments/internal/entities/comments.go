package entities

type Comment struct {
	ID       int
	PostID   int `json:"post_id"`
	ParentID int `json:"parent_id"`
	Content  string
}
