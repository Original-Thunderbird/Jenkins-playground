package main

import (
	"bytes"
	"context"
	"fmt"
	"log"

	//"io/ioutil"
	//"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/hirpc/arsenal/hihttp"
)

func main() {

	hihttp.Load(
		hihttp.WithTimeout(5*time.Second), // 设置全局超时时间
		hihttp.WithRetryCount(1),          // 设置全局重试次数
		hihttp.WithRetryWait(time.Second), // 设置全局重试等待时间
	)

	doc := getHTML("https://www.usc.edu")

	// //for offline test
	// data, _ := ioutil.ReadFile("data.html")
	// doc, _ := goquery.NewDocumentFromReader(strings.NewReader(string(data)))

	//select element "section" with id == "content-secondary"
	filterContentSecondary := func(i int, sel *goquery.Selection) bool {
		name, _ := sel.Attr("id")
		return name == "content-secondary"
	}

	//picking articles from section id == "content-secondary"
	sec := doc.Find("section").FilterFunction(filterContentSecondary)
	sec.Find("article").Each(func(i int, s *goquery.Selection) {
		ttlTg := s.Find("h1").Text()
		fmt.Println("Title:", ttlTg)
		link, _ := s.Find("h1").Find("a").Attr("href")
		fmt.Println("Link:", link)
		txt := getArticle(link)
		fmt.Println("Body:", txt)
		dscTg := s.Find("p").Text()
		fmt.Println("Description:", dscTg)
		fmt.Println()
	})

	//// picking all articles
	// doc.Find("article").Each(func(i int, s *goquery.Selection) {
	// 	ttlTg := s.Find("h1").Text()
	// 	fmt.Println("Header:", ttlTg)
	// 	dscTg := s.Find("p").Text()
	// 	fmt.Println("Description:", dscTg)
	// })
}

func getHTML(url string) *goquery.Document {
	res, err := hihttp.New().Get(context.Background(), url)
	if err != nil {
		log.Fatal(err)
	}
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(res))
	if err != nil {
		log.Fatal(err)
	}
	return doc
}

func getArticle(url string) string {

	doc := getHTML(url)

	filterArticle := func(i int, sel *goquery.Selection) bool {
		name, _ := sel.Attr("class")
		return name == "story-body"
	}

	txt := ""

	//picking articles from section id == "content-secondary"
	doc.Find("div").FilterFunction(filterArticle).Each(func(i int, s *goquery.Selection) {
		txt += s.ChildrenFiltered("p").Text()
	})
	if txt != "" {
		return txt
	}
	return "no article"
}
