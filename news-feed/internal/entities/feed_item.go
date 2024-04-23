package entities

type RSSChannel struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	Language    string `xml:"language"`
	Items       []Item `xml:"item"`
}

type Item struct {
	GUID        string `xml:"guid"`
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	PdaLink     string `xml:"pdalink"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
	Category    string `xml:"category"`
	Author      string `xml:"author"`
}

type Post struct {
	ID      int
	Title   string
	Link    string
	Content string
	PubDate int64
}

type RSS struct {
	Channel RSSChannel `xml:"channel"`
}

type Error struct {
	Error error
}
