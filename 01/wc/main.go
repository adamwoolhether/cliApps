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
	lineCount := flag.Bool("l", false, "Count lines")
	byteCount := flag.Bool("b", false, "Count bytes")
	file := flag.Bool("f", false, "Read from file(s) instead of stdin")
	// Parse the given flags.
	flag.Parse()

	if *file {
		num := 0
		for i := range flag.Args() {
			fileData, err := os.Open(flag.Arg(i))
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			num += count(fileData, *lineCount, *byteCount)
		}
		fmt.Println(num)
		return
	}

	fmt.Println(count(os.Stdin, *lineCount, *byteCount))
}

// count returns the number of words given by io.Reader.
func count(r io.Reader, countLines, countBytes bool) int {
	// A scanner reads text from the reader.
	scanner := bufio.NewScanner(r)

	// If count lines flag isn't set, we count the words.
	if !countLines {
		scanner.Split(bufio.ScanWords)
	}
	if countBytes {
		scanner.Split(bufio.ScanBytes)
	}

	// Define a counter.
	wc := 0

	// Increment the counter for every words scanned.
	for scanner.Scan() {
		wc++
	}

	return wc
}
