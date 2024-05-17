package man

import (
	"testing"
)

func TestURL(t *testing.T) {
	tests := map[string]string{
		"bookworm/ddrescue": "https://manpages.debian.org/bookworm/gddrescue/ddrescue.1.en.gz",
		"bookworm/flock.2":  "https://manpages.debian.org/bookworm/manpages-dev/flock.2.en.gz",
	}
	for query, result := range tests {
		t.Run(query, func(t *testing.T) {
			url, err := URL(query)
			if err != nil {
				t.Fatal(err)
			}
			if url != result {
				t.Errorf("unexpected result:\nwant %s\n got %s", result, url)
			}
		})
	}
}
