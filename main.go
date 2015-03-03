package main

import (
	"fmt"
	//"log"

	"github.com/PuerkitoBio/goquery"
)

const (
	startUrl     = "http://news.ycombinator.com"
	linkSelector = "a"
)

type page_data struct {
	url   string
	count int
}

func downloader(url string, page chan<- *goquery.Document) {
	doc, _ := goquery.NewDocument(url)
	if doc != nil {
		page <- doc
	}
}

// Parse the specified __page__ for all links
//  Send the number of links on that page to be analyzed to printer
//  Send the links on the page to be downloaded
func parse(page <-chan *goquery.Document, printer chan<- *page_data, hrefs chan<- []string) {
	links := make([]string, 0)
	for doc := range page {
		doc.Find(linkSelector).Each(func(i int, s *goquery.Selection) {
			link, _ := s.Attr("href")
			links = append(links, link)
		})

		go func() {
			hrefs <- links
			result := &page_data{doc.Url.String(), len(links)}
			printer <- result
		}()
	}
}

func printer(data <-chan *page_data) {
	for page_info := range data {
		fmt.Printf("%s -> %d\n", page_info.url, page_info.count)
	}
}

func main() {
	page := make(chan *goquery.Document)
	output := make(chan *page_data)
	download_queue := make(chan []string)

	go parse(page, output, download_queue)
	go downloader(startUrl, page)

	go printer(output)

	for {
		select {
		case links := <-download_queue:
			for _, url := range links {
				go downloader(url, page)
			}
		default:
			continue
		}
	}
}
