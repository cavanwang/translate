package translate

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"

	"golang.org/x/net/html"
)

func Translate(words string) {
	if words == "" {
		return
	}

	// 网络查询请求
	http.DefaultClient.Timeout = 3 * time.Second
	u := &url.URL{
		Scheme: "http",
		Host:   "dict.youdao.com",
		Path:   "/search",
	}
	q := u.Query()
	q.Add("q", words)
	u.RawQuery = q.Encode()
	resp, err := http.Get(u.String())
	if err != nil {
		log.Fatal(err)
		return
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(b)))
	if err != nil {
		log.Fatal(err)
	}

	// 发音
	doc.Find(".pronounce").Each(func(i int, s *goquery.Selection) {
		itr := s.Nodes[0].FirstChild
		for ; itr != nil; itr = itr.NextSibling {
			if itr.Type == html.TextNode {
				data := strings.TrimSpace(itr.Data)
				if data != "" {
					fmt.Printf("%+v: ", data)
				}
			}
			if itr.Type == html.ElementNode {
				citr := itr.FirstChild
				for ; citr != nil; citr = citr.NextSibling {
					if citr.Type == html.TextNode {
						fmt.Printf("%+v\n", strings.TrimSpace(citr.Data))
					}
				}
			}
		}
	})

	// 释义
	doc.Find("#phrsListTab ul li").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the text.
		fmt.Printf("%s\n", s.Text())
	})

	// 网络翻译
	fmt.Printf("网络翻译:")
	doc.Find("#webTransToggle .title span").Each(func(i int, s *goquery.Selection) {
		fmt.Printf(" %s;", strings.TrimSpace(s.Text()))
	})
	fmt.Println("")
}
