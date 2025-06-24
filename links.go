package main

import (
	"fmt"
	"net/http"
	"slices"
	"strings"

	"golang.org/x/net/html"

	"github.com/Lew-Lew/sherlock/style"
)

type ResultLink struct {
	URL   string
	Label string
	Issue string
	Err   error
}

func (rl ResultLink) String() string {
	if rl.Err != nil {
		return fmt.Sprintf("%v\n%v",
			style.IssueKey.Render(rl.URL),
			style.IssueAlert.Render(fmt.Sprintf("Error: %v", rl.Err)),
		)
	}
	return fmt.Sprintf("%v %v\n%v",
		style.IssueKey.Render(rl.URL),
		style.IssueAlert.Render(rl.Issue),
		style.IssueInfo.Render(rl.Label),
	)
}

func (s *Sherlock) CheckLinks() error {
	if len(s.SitemapURLs) == 0 {
		return fmt.Errorf("no sitemap URLs available to check links")
	}

	s.Links = []ResultLink{}

	for _, url := range s.SitemapURLs {
		res, err := checkLinks(url)
		if err != nil {
			s.Links = append(s.Links, ResultLink{
				URL: url,
				Err: err,
			})
			continue
		}
		s.Links = append(s.Links, res...)
	}
	return nil
}

func checkLinks(url string) ([]ResultLink, error) {
	resp, err := http.Get(url)
	if err != nil {
		return []ResultLink{}, err
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return []ResultLink{}, err
	}

	res := []ResultLink{}

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				// Check if the <a> tag has an href attribute
				if attr.Key == "href" && !slices.Contains([]string{"", "#"}, strings.TrimSpace(attr.Val)) {
					return
				}

				// skip tabs
				if strings.Contains(attr.Key, "w-tab") || strings.Contains(attr.Val, "w-tab") {
					return
				}

				// skip lightbox
				if attr.Key == "class" && strings.Contains(attr.Val, "w-lightbox") {
					return
				}
			}

			res = append(res, ResultLink{
				URL:   url,
				Label: nodeHTML(n),
				Issue: "Suspicious link",
			})
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return res, nil
}

func nodeHTML(n *html.Node) string {
	var b strings.Builder
	html.Render(&b, n)
	return b.String()
}
