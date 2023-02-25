package internal

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// ParseTitle fetches the title from an html document if it exists
func ParseTitle(content string) (string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(content))
	if err != nil {
		return "", err
	}

	var title string
	doc.Find("meta").EachWithBreak(func(i int, node *goquery.Selection) bool {
		if name, has := node.Attr("name"); has {
			if name == "twitter:title" {
				val, _ := node.Attr("content")
				title = val
				return false

			}
		}

		return true
	})
	return title, nil
}
