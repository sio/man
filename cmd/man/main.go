package main

import (
	"fmt"
	"os"

	"github.com/sio/man"
)

func main() {
	url, err := man.URL(query())
	if err != nil {
		fail(err)
	}
	err = man.Render(url)
	if err != nil {
		fail(err)
	}
}

func query() string {
	switch len(os.Args) {
	case 2: // MANPAGE
		return os.Args[1]
	case 3: // SECTION MANPAGE
		return fmt.Sprintf("%s.%s", os.Args[2], os.Args[1])
	default:
		usage()
		os.Exit(1)
	}
	panic("unreachable")
}

func usage() {
	_, _ = fmt.Fprintf(
		os.Stderr,
		`usage: %[1]s [SECTION] MANPAGE

Examples:
	%[1]s man
	%[1]s 2 flock
	%[1]s bookworm/ddrescue
`,
		os.Args[0],
	)
}

func fail(err error) {
	_, _ = fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
