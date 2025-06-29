package link

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCheckLinks(t *testing.T) {
	tests := map[string]struct {
		html string
		want int
	}{
		"ok": {
			html: `
				<a href="/page">basic link</a>
				<a w-tab="tab1">webflow tab</a>
				<a class="w-lightbox">webflow lightbox</a>
			`,
			want: 0,
		},
		"no-href": {
			html: `
				<a>no href</a>
				<a href="">empty href</a>
				<a href="#">empty href #</a>
			`,
			want: 3,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", "text/html")
				w.Write([]byte(tt.html))
			}))
			defer server.Close()

			results, err := checkLinks(server.URL)
			require.NoError(t, err)
			assert.Equal(t, tt.want, len(results))
		})
	}
}
