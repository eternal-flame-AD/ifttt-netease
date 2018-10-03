package api

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Song struct {
	Name string
	ID   string
}

func (c *Song) URL() string {
	return "https://music.163.com/#/song?id=" + c.ID
}

func GetSongList(id string) ([]Song, error) {
	req, err := http.NewRequest("GET", "https://music.163.com/playlist?id="+id, bytes.NewBuffer([]byte{}))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64…) Gecko/20100101 Firefox/62.0")
	req.Header.Set("Referer", "https://music.163.com/")

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Unexpected code: %d", resp.StatusCode)
	}
	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return nil, err
	}

	res := make([]Song, 0)
	doc.Find("ul.f-hide").Find("li").Each(func(_ int, elem *goquery.Selection) {
		song := Song{
			Name: elem.Text(),
			ID:   strings.TrimPrefix(elem.Find("a").AttrOr("href", ""), "/song?id="),
		}
		res = append(res, song)
	})
	return res, nil
}
