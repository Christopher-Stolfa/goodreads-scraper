package main

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly"
)

const (
	BASE_URL          = "https://www.goodreads.com"
	SEARCH_QUERY      = "/search?q="
	SEARCH_TYPE       = "&search_type="
	SEARCH_TYPE_BOOKS = "books"
)

type Book struct {
	Link string `json:"link"`
	Name string `json:"name"`
	Img  string `json:"img"`
}

func main() {
	router := gin.Default()
	router.GET("/getBooks", getBooks)
	router.Run("localhost:8080")
}

func getBooks(context *gin.Context) {
	c := colly.NewCollector()
	query := context.Query("query")
	encodedQuery := url.QueryEscape(query)
	goodreadsLink := BASE_URL + SEARCH_QUERY + encodedQuery + SEARCH_TYPE + SEARCH_TYPE_BOOKS
	var books []Book
	c.OnHTML("tr", func(e *colly.HTMLElement) {
		count := 0
		e.ForEach("td", func(_ int, el *colly.HTMLElement) {
			count++
			if count == 2 {
				nextBook := Book{
					Link: BASE_URL + e.ChildAttr("a", "href"),
					Name: "",
					Img:  "",
				}
				fmt.Println(nextBook)
				books = append(books, nextBook)
				return
			}
		})
	})
	if err := c.Visit(goodreadsLink); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch books data"})
		return
	}
	context.JSON(http.StatusOK, books)
}
