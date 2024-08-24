package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"
)

const usage = "[usage] crawler <website_url> <max_concurrency> <limit>"

func run(args []string, out io.Writer) error {
	if len(args) < 3 {
		return errors.New("missing args. " + usage)
	} else if len(args) > 3 {
		return errors.New("too many arguments provided. " + usage)
	}

	start := time.Now()

	url := args[0]
	maxConcurrency, err := strconv.Atoi(args[1])
	if err != nil {
		return fmt.Errorf("invalid max concurrency setting. error: %s", err)
	}
	if maxConcurrency < 1 {
		return errors.New("max concurrency cannot be negative")
	}

	pageLimit, err := strconv.Atoi(args[2])
	if err != nil {
		return fmt.Errorf("invalid limit setting. error: %s", err)
	}
	if pageLimit < 0 {
		return errors.New("limit cannot be negative")
	}

	if _, err := fmt.Fprintf(out, "starting crawl of: %s\n", url); err != nil {
		return err
	}

	c := newCrawler(url, maxConcurrency, pageLimit, log.New(out, "crawler: ", log.LstdFlags))
	c.crawlPage(url)

	c.wg.Wait()

	if len(c.pages) == 0 {
		_, _ = fmt.Fprintln(out, "no urls found")
	}

	for u, count := range c.pages {
		_, _ = fmt.Fprintf(out, "%s: %d\n", u, count)
	}

	_, _ = fmt.Fprintf(out, "Execution time %s. Total Queries: %d\n", time.Since(start), c.totalQueries)

	return nil
}

func main() {
	err := run(os.Args[1:], os.Stdout)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
