package main

import (
	"fmt"

	"github.com/PuerkitoBio/goquery"
)

const START_URL string = "http://news.ycombinator.com"
const SELECTOR string = "a"

type page_data struct {
	url   string
	count int
}

func download_manager(links []string, parser chan *goquery.Document) {
	for _, url := range links {
		go downloader(url, parser)
	}
}

func downloader(url string, page chan *goquery.Document) {
	doc, _ := goquery.NewDocument(url)

	page <- doc
}

//
// Parse the specified __page__ for all links
//  Send the number of links on that page to be analyzed to printer
//  Send the links on the page to be downloaded
func parse(doc *goquery.Document, printer chan *page_data, hrefs chan []string) {
	links := make([]string, 0)
	doc.Find(SELECTOR).Each(func(i int, s *goquery.Selection) {
		link, _ := s.Attr("href")
		links = append(links, link)
	})

	hrefs <- links

	result := &page_data{doc.Url.String(), len(links)}
	printer <- result
}

func printer(page_info *page_data) {
	fmt.Printf("%s -> %d\n", page_info.url, page_info.count)
}

func main() {
	downloaded_page := make(chan *goquery.Document)
	result := make(chan *page_data)
	download_queue := make(chan []string)

	go downloader(START_URL, downloaded_page)

	for {
		select {
		case page := <-downloaded_page:
			go parse(page, result, download_queue)

		case urls := <-download_queue:
			go download_manager(urls, downloaded_page)

		case report := <-result:
			printer(report)
		}
	}
}
