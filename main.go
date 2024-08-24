package main

import (
	"errors"
	"fmt"
	"io"
	"os"
)

func run(args []string, out io.Writer) error {
	if len(args) < 2 {
		return errors.New("no website provided")
	} else if len(args) > 2 {
		return errors.New("too many arguments provided")
	}

	url := args[1]
	if _, err := fmt.Fprintf(out, "starting crawl of: %s\n", url); err != nil {
		return err
	}

	pages := map[string]int{}
	crawlPage(url, url, pages)

	if len(pages) == 0 {
		_, _ = fmt.Fprintf(out, "Nothing to show for %q\n", url)
	}

	for u, count := range pages {
		_, _ = fmt.Fprintf(out, "%q: %d\n", u, count)
	}

	return nil
}

func main() {
	err := run(os.Args, os.Stdout)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
