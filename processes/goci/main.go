package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
)

func run(proj string, out io.Writer) error {

	if proj == "" {
		return fmt.Errorf("project directory is required")
	}

	args := []string{"build", ".", "errors"}

	cmd := exec.Command("go", args...)

	cmd.Dir = proj

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("'go build' failed: %s", err)
	}

	_, err := fmt.Fprintln(out, "GO BUILD: SUCCESS")
	return err
}

func main() {
	proj := flag.String("p", "", "Project directory")
	flag.Parse()

	if err := run(*proj, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
