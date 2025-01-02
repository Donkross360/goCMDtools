package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func main() {
	// Calling the count function to count the number of words
	// recieved from the standard input and printing it out
	fmt.Println(count(os.Stdin))
}

func count(r io.Reader) int {
	// A scanner is used to read text from a reader (such as files)
	scanner := bufio.NewScanner(r)

	//Define the scanner split type to words (default is split by lines)
	scanner.Split(bufio.ScanWords)

	// Define a counter
	wc := 0

	// For every word scanned, increment the counter
	for scanner.Scan() {
		wc++
	}

	// Return the total
	return wc
}
