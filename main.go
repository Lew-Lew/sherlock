package main

import (
	"fmt"
	"os"

	"github.com/Lew-Lew/sherlock/app"
	"github.com/Lew-Lew/sherlock/linter/link"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run main.go <sitemap_url>")
		return
	}

	sherlock := app.NewSherlock(
		os.Args[1],
		[]app.Linter{
			link.NewLinter(),
		},
	)

	if err := sherlock.Run(); err != nil {
		fmt.Printf("Error checking links: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(sherlock)
}
