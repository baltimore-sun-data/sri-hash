package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/baltimore-sun-data/sri-hash/sri"
)

func main() {
	const usage = `sri-hash creates sub-resource integrity hashes.

Usage of sri-hash:

    $ sri-hash <name>...

Name may be a file or URL. Defaults to standard input.
`
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), usage)
		flag.PrintDefaults()
	}
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		args = []string{""}
	}
	if err := run(args); err != nil {
		fmt.Fprintf(flag.CommandLine.Output(), "Problem creating SRI hash: %v\n", err)
	}
}

func run(resources []string) error {
	for _, rsc := range resources {
		if err := printResource(rsc); err != nil {
			return err
		}
	}
	return nil
}

func printResource(name string) error {
	if name == "" || name == "-" {
		s, err := sri.FromReader(os.Stdin)
		if err != nil {
			return err
		}
		fmt.Println(s)
		return nil
	}
	u, err := url.Parse(name)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") {
		f, err := os.Open(name)
		if err != nil {
			return err
		}
		defer f.Close()
		s, err := sri.FromReader(f)
		if err != nil {
			return err
		}
		fmt.Println(s)
		return nil
	} else {
		rsp, err := http.Get(name)
		if err != nil {
			return err
		}
		defer rsp.Body.Close()
		s, err := sri.FromReader(rsp.Body)
		if err != nil {
			return err
		}
		fmt.Println(s)
		return nil
	}
}
