package format

import (
	"strings"

	"golang.org/x/net/html"
)

func HTMLNode(n *html.Node) string {
	var b strings.Builder
	html.Render(&b, n)
	return b.String()
}
