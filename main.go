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
	if _, err := fmt.Fprintf(out, "starting crawl of: %s", url); err != nil {
		return err
	}

	rawHTML, err := getHTML(url)
	if err != nil {
		return err
	}

	if _, err = fmt.Fprintln(out, rawHTML); err != nil {
		return err
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
