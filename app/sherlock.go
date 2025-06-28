package app

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"strings"
)

type Linter interface {
	SetURLs(urls []string)
	Run() error
	Results() []Result
}

type Sherlock struct {
	URL         string
	SitemapURLs []string

	Linters []Linter
}

func NewSherlock(url string, linters []Linter) *Sherlock {
	parsed := strings.Split(strings.TrimSpace(url), "/")
	domain := ""
	if len(parsed) >= 3 && strings.HasPrefix(parsed[0], "http") {
		domain = parsed[2]
	} else {
		domain = parsed[0]
	}
	return &Sherlock{
		URL:     domain,
		Linters: linters,
	}
}

func (s *Sherlock) Run() error {
	if err := s.FetchSitemap(); err != nil {
		return fmt.Errorf("fetching sitemap: %w", err)
	}

	if len(s.SitemapURLs) == 0 {
		return fmt.Errorf("no sitemap URLs found")
	}

	for _, linter := range s.Linters {
		linter.SetURLs(s.SitemapURLs)
		if err := linter.Run(); err != nil {
			return fmt.Errorf("running linter %T: %w", linter, err)
		}
	}

	return nil
}

func (s *Sherlock) FetchSitemap() error {
	candidates := []string{
		"https://" + s.URL + "/sitemap.xml",
		"http://" + s.URL + "/sitemap.xml",
		"https://www." + s.URL + "/sitemap.xml",
		"http://www." + s.URL + "/sitemap.xml",
	}

	for _, sitemapURL := range candidates {
		resp, err := http.Get(sitemapURL)
		if err != nil || resp.StatusCode != http.StatusOK {
			if resp != nil {
				resp.Body.Close()
			}
			continue
		}
		defer resp.Body.Close()

		var sitemap struct {
			URLs []struct {
				Loc string `xml:"loc"`
			} `xml:"url"`
		}
		if err := xml.NewDecoder(resp.Body).Decode(&sitemap); err != nil {
			return err
		}
		var urls []string
		for _, u := range sitemap.URLs {
			urls = append(urls, strings.TrimSpace(u.Loc))
		}
		s.SitemapURLs = urls
		return nil
	}
	return fmt.Errorf("missing sitemap")
}

func (s *Sherlock) String() string {
	str := fmt.Sprintln("SITE:", s.URL)

	if len(s.SitemapURLs) == 0 {
		str += "missing sitemap\n\n"
		return strings.TrimSpace(str)
	}

	for _, linter := range s.Linters {
		results := linter.Results()
		if len(results) == 0 {
			continue
		}
		for _, res := range results {
			str += res.String() + "\n\n"
		}
	}

	return strings.TrimSpace(str)
}
