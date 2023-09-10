package main

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly"
)

func main() {
	router := gin.Default()
	router.GET("/getBooks", getBooks)
	router.Run("localhost:8080")
}

func getBooks(context *gin.Context) {
	c := colly.NewCollector()
	title := context.Query("title")
	encodedTitle := url.QueryEscape(title)
	var goodreadsLink = "https://www.goodreads.com/search?utf8=%E2%9C%93&q=" + encodedTitle + "&search_type=books"
	var links []string
	fmt.Println(goodreadsLink)
	c.OnHTML("tr", func(e *colly.HTMLElement) {
		count := 0
		e.ForEach("td", func(_ int, el *colly.HTMLElement) {
			count++
			if count == 2 {
				link := e.ChildAttr("a", "href")
				links = append(links, link)
				return
			}
		})
	})
	if err := c.Visit(goodreadsLink); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch data"})
		return
	}
	context.IndentedJSON(http.StatusOK, links)
}
