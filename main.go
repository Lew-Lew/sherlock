package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run main.go <sitemap_url>")
		return
	}

	sherlock := NewSherlock(os.Args[1])

	if err := sherlock.FetchSitemap(); err != nil {
		fmt.Printf("Error fetching sitemap: %v\n", err)
		os.Exit(1)
	}

	if err := sherlock.CheckLinks(); err != nil {
		fmt.Printf("Error checking links: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(sherlock)
}
