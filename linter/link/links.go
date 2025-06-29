package link

import (
	"fmt"
	"net/http"
	"slices"
	"strings"

	"golang.org/x/net/html"

	"github.com/Lew-Lew/sherlock/app"
	"github.com/Lew-Lew/sherlock/internal/format"
)

const name = "suspicious links"

type Linter struct {
	urls []string

	results []app.Result
}

func NewLinter() *Linter {
	return &Linter{}
}

func (s *Linter) SetURLs(urls []string) {
	s.urls = urls
}

func (s *Linter) Run() error {
	if len(s.urls) == 0 {
		return fmt.Errorf("no sitemap URLs available to check links")
	}

	s.results = []app.Result{}

	for _, url := range s.urls {
		res, err := checkLinks(url)
		if err != nil {
			s.results = append(
				s.results,
				app.NewErrorResult(url, err),
			)
			continue
		}
		s.results = append(s.results, res...)
	}
	return nil
}

func checkLinks(url string) ([]app.Result, error) {
	resp, err := http.Get(url)
	if err != nil {
		return []app.Result{}, err
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return []app.Result{}, err
	}

	res := []app.Result{}

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

			res = append(res, app.NewResult(url, format.HTMLNode(n), name))
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return res, nil
}

func (s *Linter) Results() []app.Result {
	return s.results
}
