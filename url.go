package man

import (
	"fmt"
	"net/http"
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
func URL(query string) (string, error) {
	head, err := client.Head(repository + "/" + query)
	if err != nil {
		return "", err
	}
	defer func() { _ = head.Body.Close() }()
	if head.StatusCode < 300 || head.StatusCode > 399 {
		return "", fmt.Errorf("did not get a redirect for %q: %s", query, head.Status)
	}
	url, err := head.Location()
	if err != nil {
		return "", err
	}
	prefix, found := strings.CutSuffix(url.Path, ".html")
	if !found {
		return "", fmt.Errorf("manpage not found: %s", url)
	}
	url.Path = prefix + ".gz"
	return url.String(), nil
}
