package main

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"strings"
)

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
