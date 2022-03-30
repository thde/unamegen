package unamegen

import (
	"encoding/json"
	"math/rand"
	"strings"
	"time"

	_ "embed"
)

const maxSeq = 2

//go:embed characters.json
var charactersJSON []byte
var random = rand.New(rand.NewSource(time.Now().UnixNano()))

type (
	WordProbability map[string]map[string]float32
)

func New() (WordProbability, error) {
	w := WordProbability{}
	err := json.Unmarshal(charactersJSON, &w)

	return w, err
}

func (w WordProbability) GenerateFakeWordWithUnexpectedLength() string {
	character := "^"
	word := ""
	var characters []string
	for character != "END" {
		characters = append(characters, character)
		if len(characters) > maxSeq {
			characters = characters[1:]
		}
		nextAccumedProbs := map[string]float32{} //nolint
		n := 0
		for {
			str := strings.Join(characters[n:], "")
			nextAccumedProbs = w[str]
			n += 1
			if !(nextAccumedProbs == nil && n < len(characters)) {
				break
			}
		}
		nextCharacter := ""
		r := random.Float32()
		for ch, prob := range nextAccumedProbs {
			nextCharacterCandidate := ch
			probability := prob
			if r <= probability {
				nextCharacter = nextCharacterCandidate
				break
			}
		}
		if nextCharacter != "END" {
			word += nextCharacter
		}
		character = nextCharacter
	}
	return word
}

func (w WordProbability) GenerateFakeWordByLength(length int) string {
	fakeWord := ""
	for len(fakeWord) != length {
		fakeWord = w.GenerateFakeWordWithUnexpectedLength()
	}
	return fakeWord
}

func (w WordProbability) GenerateFakeWord(minLength int, maxLength int) string {
	fakeWord := ""
	for !(minLength <= len(fakeWord) && len(fakeWord) <= maxLength) {
		fakeWord = w.GenerateFakeWordWithUnexpectedLength()
	}
	return fakeWord
}
