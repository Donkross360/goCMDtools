package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"sync"
)

func main() {
	// verify and parse argument
	op := flag.String("op", "sum", "Operation to be executed")
	column := flag.Int("col", 1, "CSV column on which to execute operation")
	flag.Parse()

	if err := run(flag.Args(), *op, *column, os.Stderr); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(filenames []string, op string, column int, out io.Writer) error {
	var opFunc statsFunc

	if len(filenames) == 0 {
		return ErrNoFiles
	}

	if column < 1 {
		return fmt.Errorf("%w: %d", ErrInvalidColumn, column)
	}

	// Validate the operation and define the opFunc accordingly
	switch op {
	case "sum":
		opFunc = sum
	case "avg":
		opFunc = avg
	default:
		return fmt.Errorf("%w: %s", ErrInvalidOperation, op)
	}

	consolidate := make([]float64, 0)

	// Create the channel to recieve result or error of operation
	resCh := make(chan []float64)
	errCh := make(chan error)
	doneCh := make(chan struct{})

	wg := sync.WaitGroup{}

	// Loop through all files and create a goroutine to process
	// each one concurrently
	for _, fname := range filenames {
		wg.Add(1)
		go func(fname string) {
			defer wg.Done()

			// Open the file for reading
			f, err := os.Open(fname)
			if err != nil {
				errCh <- fmt.Errorf("connot open file: %w", err)
				return
			}

			// Parse the csv into a slice os float64 numbers
			data, err := csv2float(f, column)
			if err != nil {
				errCh <- err
			}

			if err := f.Close(); err != nil {
				errCh <- err
			}

			resCh <- data
		}(fname)
	}

	go func() {
		wg.Wait()
		close(doneCh)
	}()

	for {
		select {
		case err := <-errCh:
			return err
		case data := <-resCh:
			consolidate = append(consolidate, data...)
		case <-doneCh:
			_, err := fmt.Fprintln(out, opFunc(consolidate))
			return err
		}
	}

}
