package entities

type NewsFullDetailed struct {
	ID      int
	Title   string
	Link    string
	Content string
	PubDate int64
}

type NewsShortDetailed struct {
	ID      int
	Title   string
	Link    string
	Content string
	PubDate int64
}

type Comment struct {
	ID      int
	PostID  int
	Content string
}
