package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type config struct {
	// extention to filter out
	ext string

	// min file size
	size int64

	// list files
	list bool

	// delete files
	del bool
}

func main() {
	// Parse command line flags
	root := flag.String("root", "", "Root directory to start")

	// Action options
	list := flag.Bool("list", false, "List file only")

	del := flag.Bool("del", false, "Delete files")

	// Filter options
	ext := flag.String("ext", "", "File extension to filter out")
	size := flag.Int64("size", 0, "Minimum file size")
	flag.Parse()

	c := config{
		ext:  *ext,
		size: *size,
		list: *list,
		del:  *del,
	}

	if err := run(*root, os.Stdout, c); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

}

func run(root string, out io.Writer, cfg config) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if filterOut(path, cfg.ext, cfg.size, info) {
			return nil
		}

		// if list was explitcity set, don't do anything else
		if cfg.list {
			return listFile(path, out)
		}

		// Delete files
		if cfg.del {
			return delFile(path)
		}

		// list is the default option if nothing else was set
		return listFile(path, out)
	})
}
