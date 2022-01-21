package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {
	// Define a boolean flag -l to count lines instead of words.
	lines := flag.Bool("l", false, "Count lines")
	// Parse the given flags.
	flag.Parse()

	fmt.Println(count(os.Stdin, *lines))
}

// count returns the number of words given by io.Reader.
func count(r io.Reader, countLines bool) int {
	// A scanner reads text from a reader.
	scanner := bufio.NewScanner(r)

	// If count lines flag isn't set, we count the words.
	if !countLines {
		scanner.Split(bufio.ScanWords)
	}

	// Define a counter.
	wc := 0

	// Increment the counter for every words scanned.
	for scanner.Scan() {
		wc++
	}

	return wc
}
