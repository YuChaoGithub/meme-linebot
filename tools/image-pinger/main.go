package main

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/gocolly/colly/v2"
)

const (
	targetURL = "https://meme-linebot.herokuapp.com/"
)

func main() {
	count := 0

	c := colly.NewCollector(
		colly.URLFilters(
			regexp.MustCompile(`https://i\.imgur\.com/.+`),
			regexp.MustCompile(targetURL),
		),
		colly.MaxDepth(2),
	)
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		e.Request.Visit(e.Attr("href"))
	})

	c.OnResponseHeaders(func(r *colly.Response) {
		if r.StatusCode != http.StatusOK {
			fmt.Println("Failed to visit", r.Request.URL)
		} else {
			count++
		}
	})

	c.Visit(targetURL)

	fmt.Println("Visited:", count-1, "images.")
}
