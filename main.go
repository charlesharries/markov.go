package main

import (
	"bufio"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

var start string = "START"
var end string = "END"

type Markov struct {
	n       int
	isWords bool
	chain   map[string][]string
	lines   []string
}

// Parse the string into ngrams and add those ngrams to the chain.
func (m *Markov) digest(str string) {
	gram := m.ngram(str)
	gram = append(gram, end)

	if m.chain == nil {
		m.chain = make(map[string][]string)
	}

	current := start
	for _, g := range gram {
		if !(current == start && g == end) {
			m.chain[current] = append(m.chain[current], g)
		}
		current = g
	}
}

func (m *Markov) divider() string {
	if m.isWords {
		return " "
	}

	return ""
}

// Generate a new string based on m's current markov chain.
func (m *Markov) generate() string {
	current := start
	result := ""

	if m.isWords {
		return m.generateWords()
	}

	for {
		rand.Seed(time.Now().UnixNano())
		randIdx := rand.Intn(len(m.chain[current]))
		current = m.chain[current][randIdx]
		if current == end {
			return result
		}

		// Remove m.n - 1 characters from end of string before adding
		// current, so as not to duplicate, i.e. hel + ell => hell
		toRemove := m.n - 1
		if toRemove < len(result) {
			result = result[:(len(result) - toRemove)]
		}
		result = fmt.Sprintf("%s%s", result, current)
	}
}

func (m *Markov) generateWords() string {
	current := start
	result := []string{}
	rand.Seed(time.Now().UnixNano())

	for {
		randIdx := rand.Intn(len(m.chain[current]))
		current = m.chain[current][randIdx]
		if current == end {
			return strings.Join(result, m.divider())
		}

		// Remove m.n - 1 words from end of string before adding
		// current, so as not to duplicate, i.e. hel + ell => hell
		toRemove := m.n - 1
		if toRemove < len(result) {
			result = result[0:(len(result) - toRemove)]
		}

		// Split words before adding to result, so that we can count
		// them in the next iteration.
		words := strings.Split(current, m.divider())
		for _, w := range words {
			result = append(result, strings.Trim(w, "!“”()"))
		}
	}
}

// Break a string into a slice of n-grams of length determined by
// m.n, e.g.:
//   lol, 2 => [lo, ol]
//   charles, 4 => [char, harl, arle, rles]
//   my name is charles harries, 2 => [my name, name is, is charles, charles harries]
func (m *Markov) ngram(str string) []string {
	words := strings.Split(str, m.divider())
	var results []string
	for i := range words {
		if len(words) < i+m.n {
			break
		}
		result := strings.Join(words[i:i+m.n], m.divider())
		results = append(results, result)
	}

	return results
}

// Does the given string exist in the markov's input?
func (m *Markov) has(str string) bool {
	for _, l := range m.lines {
		if strings.TrimSpace(l) == str {
			return true
		}
	}
	return false
}

func main() {
	var markov Markov
	var count int

	flag.IntVar(&markov.n, "n", 2, "Length of n-gram")
	flag.BoolVar(&markov.isWords, "w", false, "Split text into words")
	flag.IntVar(&count, "c", 10, "Number of results to generate")
	flag.Parse()

	// Digest standard input, one line at a time
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		t := scanner.Text()
		markov.lines = append(markov.lines, t)
		markov.digest(t)
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}

	for i := 0; i < count; i++ {
		res := markov.generate()
		status := "new"
		if markov.has(res) {
			status = "exists"
		}

		fmt.Printf("%s\t%s\n", status, res)
	}
}
