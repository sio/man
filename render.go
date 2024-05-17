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
	url, err := URL(query)
	if err != nil {
		return err
	}
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()
	width := getOutputWidth()
	page := addHeader(resp.Body, query, url, width)

	retry, err := showMan(page)
	if !retry {
		return err
	}
	return showGroff(page, width)
}

// Show manpage using system man viewer
func showMan(page io.Reader) (retry bool, err error) {
	exe, err := exec.LookPath("man")
	if err != nil {
		return true, err
	}
	man := exec.Command(exe, "-l", "-")
	man.Stdin = page
	man.Stdout = os.Stdout
	return false, man.Run()
}

// Render manpage with groff and show it in pager
func showGroff(page io.Reader, width int) error {
	groff, err := exec.LookPath("groff")
	if err != nil {
		return err
	}
	formatter := exec.Command(
		groff,
		"-mtty-char",
		"-Tutf8",
		"-mandoc",
		fmt.Sprintf("-rLL=%dn", width),
		fmt.Sprintf("-rLT=%dn", width),
	)
	formatter.Stdin = page
	pipe, err := formatter.StdoutPipe()
	if err != nil {
		return err
	}
	defer func() { _ = pipe.Close() }()
	pager, err := getPager()
	if err != nil {
		return err
	}
	pager.Stdin = pipe
	pager.Stdout = os.Stdout
	err = formatter.Start()
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
	err = formatter.Wait()
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

func addHeader(page io.Reader, title, url string, width int) io.Reader {
	url, ok := strings.CutSuffix(url, ".gz")
	if ok {
		url += ".html"
	}
	header := new(bytes.Buffer)
	_, err := fmt.Fprintf(
		header,
		`.lt %dn
.tl 'Debian Manpages (%s)''\fI\,%s\/\fR'

`,
		width,
		title,
		url)
	if err != nil {
		panic("writing to in-memory buffer failed")
	}
	return io.MultiReader(header, page)
}

func getOutputWidth() int {
	const defaultWidth = 80
	width := defaultWidth
	fd := int(os.Stdout.Fd())
	if term.IsTerminal(fd) {
		var err error
		width, _, err = term.GetSize(fd)
		if err != nil {
			width = defaultWidth
		}
	}
	width -= 1
	return width
}
