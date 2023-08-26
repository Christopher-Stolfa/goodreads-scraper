package main

import (
	"fmt"

	"github.com/gocolly/colly"
)

func main() {
	c := colly.NewCollector()
	c.OnHTML("tr", func(e *colly.HTMLElement) {
		count := 0
		e.ForEach("td", func(_ int, el *colly.HTMLElement) {
			count++
			if count == 2 {
				link := e.ChildAttrs("a", "href")
				fmt.Println(link)
				return
			}
		})
	})
	c.Visit("https://www.goodreads.com/search?utf8=%E2%9C%93&q=game+of+thrones&search_type=books")
}
