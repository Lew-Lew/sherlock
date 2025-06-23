package main

import (
	"fmt"
	"strings"
)

type Sherlock struct {
	URL string

	SitemapURLs []string
	Links       []ResultLink
}

func NewSherlock(url string) *Sherlock {
	parsed := strings.Split(strings.TrimSpace(url), "/")
	domain := ""
	if len(parsed) >= 3 && strings.HasPrefix(parsed[0], "http") {
		domain = parsed[2]
	} else {
		domain = parsed[0]
	}
	return &Sherlock{
		URL: domain,
	}
}

func (s *Sherlock) String() string {
	res := fmt.Sprintln("SITE:", s.URL)

	if len(s.SitemapURLs) == 0 {
		res += "missing sitemap\n\n"
	}

	if len(s.Links) != 0 {
		resLinks := []string{}
		for _, url := range s.Links {
			resLinks = append(resLinks, url.String())
		}
		res += strings.Join(resLinks, "\n\n") + "\n"
	}

	return strings.TrimSpace(res)
}
