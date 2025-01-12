package main

import (
	"flag"
	"fmt"
	"io"
	"log"
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

	// log destination writer
	wlog io.Writer

	// archive directory
	archive string
}

func main() {
	// Parse command line flags
	root := flag.String("root", ".", "Root directory to start")

	logFile := flag.String("log", "", "Log delete to this file")

	// Action options
	list := flag.Bool("list", false, "List file only")

	archive := flag.String("archive", "", "Archive directory")

	del := flag.Bool("del", false, "Delete files")

	// Filter options
	ext := flag.String("ext", "", "File extension to filter out")
	size := flag.Int64("size", 0, "Minimum file size")
	flag.Parse()

	var (
		f   = os.Stdout
		err error
	)

	if *logFile != "" {
		f, err = os.OpenFile(*logFile, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		defer f.Close()
	}

	c := config{
		ext:     *ext,
		size:    *size,
		list:    *list,
		del:     *del,
		wlog:    f,
		archive: *archive,
	}

	if err := run(*root, os.Stdout, c); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

}

func run(root string, out io.Writer, cfg config) error {
	delLogger := log.New(cfg.wlog, "DELETED FILE: ", log.LstdFlags)
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

		// Archive files and continue if successful
		if cfg.archive != "" {
			if err := archiveFile(cfg.archive, root, path); err != nil {
				return err
			}
		}

		// Delete files
		if cfg.del {
			return delFile(path, delLogger)
		}

		// list is the default option if nothing else was set
		return listFile(path, out)
	})
}
