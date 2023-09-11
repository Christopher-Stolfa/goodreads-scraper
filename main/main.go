package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/chromedp/chromedp"
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
	router.GET("/getBookByLink", getBookByLink)
	router.Run("localhost:8080")
}

func getBooks(ginContext *gin.Context) {
	c := colly.NewCollector()
	query := ginContext.Query("query")
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
		ginContext.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch books data"})
		return
	}
	ginContext.JSON(http.StatusOK, books)
}

func getBookByLink(ginContext *gin.Context) {
	link := ginContext.Query("link")
	// var title string
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var res string
	err := chromedp.Run(ctx,
		chromedp.Navigate(link),
		chromedp.OuterHTML("html", &res),
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)
	ginContext.JSON(http.StatusOK, res)
}
