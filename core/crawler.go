package core

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"strings"
)

var VisitedLinks = make(map[string]bool)

func CrawlWebsite(baseURL string, currentURL string, depth int) {
	if depth == 0 || VisitedLinks[currentURL] {
		return
	}

	VisitedLinks[currentURL] = true
	resp, err := http.Get(currentURL)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	doc.Find("a").Each(func(index int, item *goquery.Selection) {
		link, exists := item.Attr("href")
		if !exists {
			return
		}

		if strings.HasPrefix(link, "/") {
			fullURL := baseURL + link
			CrawlWebsite(baseURL, fullURL, depth-1)
		}
	})
}
