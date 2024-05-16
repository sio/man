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
	rendererCmd, err := getRenderer()
	if err != nil {
		return err
	}
	pagerCmd, err := getPager()
	if err != nil {
		return err
	}

	from, err := URL(query)
	if err != nil {
		return err
	}
	resp, err := client.Get(from.String())
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	header := new(bytes.Buffer)
	_, err = fmt.Fprintf(header, `.TH "manpages.debian.org " "%s"%s`, strings.ReplaceAll(query, `"`, ""), "\n")
	if err != nil {
		return err
	}

	renderer := exec.Command(rendererCmd[0], rendererCmd[1:]...)
	renderer.Stdin = io.MultiReader(header, resp.Body)
	pipe, err := renderer.StdoutPipe()
	if err != nil {
		return err
	}
	defer func() { _ = pipe.Close() }()

	pager := exec.Command(pagerCmd)
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
	err = renderer.Wait()
	if err != nil {
		return err
	}
	return nil
}

func getPager() (string, error) {
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
			return pager, nil
		}
	}
	return "", fmt.Errorf("pager detection failed")
}

func getRenderer() ([]string, error) {
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
	return []string{groff, "-Tutf8", "-man", fmt.Sprintf("-rLL=%dn", width), fmt.Sprintf("-rLT=%dn", width)}, nil
}
