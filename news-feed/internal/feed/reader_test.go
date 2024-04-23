package feed

import (
	"fmt"
	"news_feed/internal/config"
	"news_feed/internal/entities"
	"testing"
	"time"
)

type fakeRepo struct {
	Items []entities.Item
}

func (r *fakeRepo) AddItem(item entities.Item) error {
	r.Items = append(r.Items, item)
	return nil
}
func (r *fakeRepo) ReadItems(page int, per_page int, s string) ([]entities.Post, error) {
	return []entities.Post{}, nil
}

func TestFetchFeeds(t *testing.T) {
	reader := func(url string) ([]byte, error) {
		news := `
		<rss xmlns:dc="http://purl.org/dc/elements/1.1/" version="2.0">
		<channel>
			<item>
				<title><![CDATA[ [Перевод] JSON in GO ]]></title>
				<guid isPermaLink="true">https://habr.com/ru/articles/797019/</guid>
				<link>https://habr.com/ru/articles/797019</link>
				<description><![CDATA[Это перевод одноименной статьи ]]></description>
				<pubDate>Wed, 28 Feb 2024 20:37:41 GMT</pubDate>
				<category>golang</category>
			</item>
		</channel>
		</rss>
		`
		response := []byte(news)
		return response, nil
	}

	cfg := config.Config{
		RequestPeriod: 100,
		RSS:           []string{""},
		Duration:      time.Millisecond,
	}

	repo := &fakeRepo{
		Items: make([]entities.Item, 0),
	}

	go FeatchFeeds(cfg, reader, repo)

	time.Sleep(110 * time.Millisecond)
	fmt.Printf("Hовости: %#v\n", repo.Items)
	if len(repo.Items) != 1 {
		t.Fatal("Должны быть новости")
	}
}
