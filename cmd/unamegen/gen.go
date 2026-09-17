//go:build ignore

// gen produces the embedded word list for unamegen: the first N words of
// a dictionary file, gzip compressed. Run via `go generate`.
package main

import (
	"bufio"
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "gen:", err)
		os.Exit(1)
	}
}

func run() error {
	in := flag.String("in", "../../dictionaries/en.txt", "dictionary file, one word per line")
	out := flag.String("out", "en.txt.gz", "output file")
	n := flag.Int("n", 3000, "number of words to keep")
	flag.Parse()

	src, err := os.Open(*in)
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(*out)
	if err != nil {
		return err
	}
	defer dst.Close()

	gz, err := gzip.NewWriterLevel(dst, gzip.BestCompression)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(src)
	for i := 0; i < *n && scanner.Scan(); i++ {
		if _, err := io.WriteString(gz, scanner.Text()+"\n"); err != nil {
			return err
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	if err := gz.Close(); err != nil {
		return err
	}
	return dst.Close()
}
