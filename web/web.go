//go:generate curl -s https://raw.githubusercontent.com/tinygo-org/tinygo/v0.42.0/targets/wasm_exec.js -o wasm_exec.js
package main

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"syscall/js"

	"thde.io/fakeword"
)

func main() {
	document := js.Global().Get("document")
	output := document.Call("getElementById", "output")
	words := document.Call("getElementById", "words")

	p := fakeword.Dictionary{}
	p.Read(strings.NewReader(words.Get("value").String()))
	generator := p.Generator()

	fmt.Println("main()")
	output.Set("innerText", generate(28, generator))
	document.Call("getElementById", "generate").Set("onclick", js.FuncOf(func(this js.Value, args []js.Value) any {
		output.Set("innerText", generate(28, generator))
		return nil
	}))
	select {} // keep running
}

func generate(amount int, g fakeword.Generator) string {
	if amount <= 0 {
		return ""
	}

	words := []string{}
	for i := 0; amount < 0 || i < amount; i++ {
		words = append(words, g.WordWithDistance(6, 12))
	}
	out := &bytes.Buffer{}
	table(out, words, 4)
	return out.String()
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
	for row := range rows {
		for col := range cols {
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
