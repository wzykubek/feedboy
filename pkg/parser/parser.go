package parser

import (
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gorilla/feeds"
)

type Parser struct {
	Scheme *Scheme
	doc    *goquery.Document
	Feed   *feeds.Feed
}

func fetchContent(url string) (io.ReadCloser, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	return res.Body, nil
}

func (p *Parser) parseDoc(html string) error {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return err
	}
	p.doc = doc

	return nil
}

func (p *Parser) MapTextToValue(selector string, value *string) {
	*value = strings.TrimSpace(p.doc.Find(selector).Text())
}

func (p *Parser) MapLinkToValue(selector string, value **feeds.Link) {
	href, exists := p.doc.Find(selector).Attr("href")
	if exists {
		*value = &feeds.Link{
			Href: href,
		}
	}
}

func (p *Parser) MapDateToValue(selector string, value *time.Time) {
	dateStr := strings.TrimSpace(p.doc.Find(selector).Text())
	if dateStr != "" {
		date, err := time.Parse(p.Scheme.DateFormat, dateStr)
		if err != nil {
			panic(err)
		}

		*value = date
	}
}

func (p *Parser) Do() {
	body, err := fetchContent(p.Scheme.Url)
	if err != nil {
		panic(err)
	}
	defer body.Close()

	p.doc, err = goquery.NewDocumentFromReader(body)
	if err != nil {
		panic(err)
	}

	p.Feed = &feeds.Feed{
		Title:       p.Scheme.Title,
		Description: p.Scheme.Description,
		Link:        &feeds.Link{Href: p.Scheme.Url},
	}
	p.Feed.Items = []*feeds.Item{}

	p.doc.Find(p.Scheme.Selectors.Item).Each(func(i int, s *goquery.Selection) {
		html, err := s.Html()
		if err != nil {
			panic(err)
		}

		err = p.parseDoc(html)
		if err != nil {
			log.Printf("parsing error: %v", err)
			return
		}

		item := &feeds.Item{}
		p.MapTextToValue(p.Scheme.Selectors.ItemTitle, &item.Title)
		p.MapTextToValue(p.Scheme.Selectors.ItemDescription, &item.Description)
		p.MapTextToValue(p.Scheme.Selectors.ItemContent, &item.Content)
		p.MapLinkToValue(p.Scheme.Selectors.ItemUrl, &item.Link)
		p.MapDateToValue(p.Scheme.Selectors.ItemDate, &item.Created)

		p.Feed.Items = append(p.Feed.Items, item)
	})

	println(p.Feed.ToRss())

}
