package man

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const repository = "https://manpages.debian.org"

var client = &http.Client{
	Timeout: time.Second * 5,
	CheckRedirect: func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse // do not follow redirects
	},
}

// Search manpages.debian.org for a manual and return its URL
func URL(query string) (*url.URL, error) {
	head, err := client.Head(repository + "/" + query)
	if err != nil {
		return nil, err
	}
	defer func() { _ = head.Body.Close() }()
	if head.StatusCode < 300 || head.StatusCode > 399 {
		return nil, fmt.Errorf("did not get a redirect for %q: %s", query, head.Status)
	}
	url, err := head.Location()
	if err != nil {
		return nil, err
	}
	prefix, found := strings.CutSuffix(url.Path, ".html")
	if !found {
		return nil, fmt.Errorf("manpage not found: %s", url)
	}
	url.Path = prefix + ".gz"
	return url, nil
}
