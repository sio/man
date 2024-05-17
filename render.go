package man

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"golang.org/x/term"
)

// Render a man page
func Render(query string) error {
	pager, err := getPager()
	if err != nil {
		return err
	}
	url, err := URL(query)
	if err != nil {
		return err
	}
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()
	renderer, err := setupRenderer(query, url, resp.Body)
	if err != nil {
		return err
	}
	pipe, err := renderer.StdoutPipe()
	if err != nil {
		return err
	}
	defer func() { _ = pipe.Close() }()
	pager.Stdin = pipe
	pager.Stdout = os.Stdout
	err = renderer.Start()
	if err != nil {
		return err
	}
	err = pager.Run()
	if err != nil {
		return err
	}
	err = pipe.Close()
	if err != nil {
		return err
	}
	err = renderer.Wait()
	if err != nil {
		return err
	}
	return nil
}

func getPager() (*exec.Cmd, error) {
	var pagers []string
	for _, env := range []string{"MANPAGER", "PAGER"} {
		value := os.Getenv(env)
		if value == "" {
			continue
		}
		// TODO: commandline arguments are not supported in $PAGER and $MANPAGER
		pagers = append(pagers, value)
	}
	pagers = append(pagers, "less", "more", "cat")
	for _, pager := range pagers {
		pager, err := exec.LookPath(pager)
		if err == nil {
			cmd := exec.Command(pager)
			cmd.Env = append(
				os.Environ(),
				fmt.Sprintf("LESS=-ix7R%s", os.Getenv("LESS")),
				"LESSCHARSET=utf-8",
			)
			return cmd, nil
		}
	}
	return nil, fmt.Errorf("pager detection failed")
}

func setupRenderer(title, url string, input io.Reader) (*exec.Cmd, error) {
	groff, err := exec.LookPath("groff")
	if err != nil {
		return nil, err
	}

	fd := int(os.Stdout.Fd())
	width := 79
	if term.IsTerminal(fd) {
		width, _, err = term.GetSize(fd)
		if err != nil {
			return nil, err
		}
		width--
	}

	url, ok := strings.CutSuffix(url, ".gz")
	if ok {
		url += ".html"
	}

	header := new(bytes.Buffer)
	_, err = fmt.Fprintf(
		header,
		`.lt %dn
.tl 'Debian Manpages (%s)''\fI\,%s\/\fR'

`,
		width,
		title,
		url)
	if err != nil {
		return nil, err
	}

	renderer := exec.Command(
		groff,
		"-Tutf8",
		"-man",
		fmt.Sprintf("-rLL=%dn", width),
		fmt.Sprintf("-rLT=%dn", width),
	)
	renderer.Stdin = io.MultiReader(header, input)
	return renderer, nil
}
