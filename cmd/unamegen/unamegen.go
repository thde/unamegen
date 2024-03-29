package main

import (
	"bytes"
	"compress/gzip"
	_ "embed"
	"encoding/gob"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	flag "github.com/spf13/pflag"

	"thde.io/fakeword"
)

//go:generate sh -c "head -n 3000 ../../dictionaries/en.txt | go run ../calculate/calculate.go --gzip > en.gob.gz"

//go:embed en.gob.gz
var en []byte

var (
	version string
	commit  string
	date    string

	min       = flag.Int("min", 6, "min length of fake word")
	max       = flag.Int("max", 12, "max length of fake word")
	in        = flag.StringP("input", "i", "", "Input file")
	noColumns = flag.BoolP("no-columns", "1", false, "Don't print the generated usernames in columns")

	help = flag.BoolP("help", "h", false, "print help message")
)

func run() error {
	flag.Parse()

	if *help {
		fmt.Fprintf(os.Stderr, "Usage of %s [ amount ]:\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "(%s, %s %s)\n", version, commit, date)
		return nil
	}

	amount := 28
	if flag.Arg(0) != "" {
		a, err := strconv.Atoi(flag.Arg(0))
		if err != nil {
			return fmt.Errorf("amount error: %w", err)
		}
		if amount <= 0 {
			return fmt.Errorf("amount '%d' is equal or smaller than 0", amount)
		}

		amount = a
	}

	if *min > *max {
		return fmt.Errorf("min larger than max")
	}

	var w fakeword.Generator
	if *in == "" {
		reader, err := gzip.NewReader(bytes.NewBuffer(en))
		if err != nil {
			return fmt.Errorf("internal words parsing error: %w", err)
		}

		dec := gob.NewDecoder(reader)
		dec.Decode(&w)
	} else {
		file, err := os.Open(*in)
		if err != nil {
			return fmt.Errorf("file %s: %w", *in, err)
		}
		defer file.Close()

		p := fakeword.Dictionary{}
		w = p.Read(file).Generator()
	}

	words := []string{}
	for i := 0; amount < 0 || i < amount; i++ {
		words = append(words, w.WordWithDistance(*min, *max))
	}

	if *noColumns {
		fmt.Println(strings.Join(words, "\n"))
	} else {
		table(os.Stdout, words, 4)
	}

	return nil
}

func table(w io.Writer, input []string, cols int) {
	maxWidth := 0
	for _, s := range input {
		if len(s) > maxWidth {
			maxWidth = len(s)
		}
	}
	format := fmt.Sprintf("%%-%ds%%s ", maxWidth)

	rows := (len(input) + cols - 1) / cols
	for row := 0; row < rows; row++ {
		for col := 0; col < cols; col++ {
			i := col*rows + row
			if i >= len(input) {
				break // This means the last column is not "full"
			}
			padding := ""
			if i < 7 {
				padding = " "
			}
			fmt.Fprintf(w, format, input[i], padding)
		}
		fmt.Fprintln(w)
	}
}

func main() {
	if err := run(); err != nil {
		fmt.Printf("Error: %s", err.Error())
		os.Exit(1)
	}
}
