package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

func run(args []string, out io.Writer) error {
	if len(args) < 2 {
		return errors.New("no website provided")
	} else if len(args) > 2 {
		return errors.New("too many arguments provided")
	}

	start := time.Now()

	url := args[1]
	if _, err := fmt.Fprintf(out, "starting crawl of: %s\n", url); err != nil {
		return err
	}

	c := newCrawler(url, log.New(out, "crawler: ", log.LstdFlags))
	c.crawlPage(url)

	c.wg.Wait()

	if len(c.pages) == 0 {
		_, _ = fmt.Fprintln(out, "no urls found")
	}

	for u, count := range c.pages {
		_, _ = fmt.Fprintf(out, "%s: %d\n", u, count)
	}

	_, _ = fmt.Fprintf(out, "Execution time %s\n", time.Since(start))

	return nil
}

func main() {
	err := run(os.Args, os.Stdout)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
