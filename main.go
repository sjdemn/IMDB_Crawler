package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"

	"github.com/gocolly/colly"
)

type star struct {
	Name      string
	Photo     string
	JobTitle  string
	BirthDay  string
	Bio       string
	TopMovies []movie
}

type movie struct {
	Title string
	Year  string
}

func main() {
	month := flag.Int("month", 1, "Month to fetch birthday for")
	day := flag.Int("day", 1, "Day to fetch birthday for")
	flag.Parse()
	crawl(*month, *day)
}

func crawl(month, day int) {
	c := colly.NewCollector(
		colly.AllowedDomains("imdb.com", "www.imdb.com"),
	)
	infoCollector := c.Clone()

	c.OnHTML(".mode-detail", func(e *colly.HTMLElement) {
		profileUrl := e.ChildAttr("div.lister-item-image > a", "href")
		profileUrl = e.Request.AbsoluteURL(profileUrl)
		infoCollector.Visit(profileUrl)
	})

	c.OnHTML("a.lister-page-next", func(e *colly.HTMLElement) {
		nextPage := e.Request.AbsoluteURL(e.Attr("href"))
		c.Visit(nextPage)
	})

	infoCollector.OnHTML("#content-2-wide", func(h *colly.HTMLElement) {
		tmpProfile := star{}
		tmpProfile.Name = h.ChildText("h1.header > span.itemprop")
		tmpProfile.Photo = h.ChildAttr("#name-poster", "src")
		tmpProfile.JobTitle = h.ChildText("#name-job-categories > a > span.itemprop")
		tmpProfile.BirthDay = h.ChildAttr("#name-born-info time", "datetime")
		tmpProfile.Bio = h.ChildText("#name-bio-text > div.name-trivia-bio-text > div.inline")

		h.ForEach("div.knownfor-title", func(_ int, h *colly.HTMLElement) {
			tmpMovie := movie{}
			tmpMovie.Title = h.ChildText("div.knownfor-title-role > a.knownfor-ellipsis")
			tmpMovie.Year = h.ChildText("div.knownfor-year > span.knownfor-ellipsis")
			tmpProfile.TopMovies = append(tmpProfile.TopMovies, tmpMovie)
		})

		js, err := json.MarshalIndent(tmpProfile, "", "   ")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(js))
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Printf("Visiting: %v\n", r.URL.String())
	})

	infoCollector.OnRequest(func(r *colly.Request) {
		fmt.Printf("visiting profile url: %v\n", r.URL.String())
	})

	startUrl := fmt.Sprintf("https://www.imdb.com/search/name/?birth_monthday=%d-%d", month, day)
	c.Visit(startUrl)
}
