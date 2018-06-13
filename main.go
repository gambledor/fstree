package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

const (
	// version is the software version
	version string = "0.1"
	// author is the software author
	author   = "Giuseppe Lo Brutto"
	maxLevel = 1<<32 - 1 // 2^32
)

type configuration struct {
	// dirOnly true to print directory names only
	dirOnly bool
	// showDotFiles true to print hidden files
	showDotFiles bool
	// level to limit the recursion depth
	level int
}

var (
	// Build is to compile passing -ldflags "-X main.Build=<sha1>"
	Build       string
	fileCounter int
	dirCounter  int
	showVersion bool
	config      configuration
)

func init() {
	flag.BoolVar(&showVersion, "version", false, "Show current program version")
	flag.BoolVar(&config.dirOnly, "d", false, "Show directories only ")
	flag.BoolVar(&config.showDotFiles, "a", false, "Show dot files ")
	flag.IntVar(&config.level, "l", maxLevel, "Limit the level of recursion depth ")
}

func printVersion() {
	if showVersion {
		fmt.Printf("fstree \033[32m%s-%s\033[0m, created by \033[96m%s\n", version, Build, author)
		os.Exit(0)
	}
}

func main() {
	flag.Parse()
	printVersion()
	args := []string{"."}
	if len(flag.Args()) > 0 {
		args = flag.Args()
	}

	for _, path := range args {
		if err := tree(path, "", 0); err != nil {
			log.Printf("tree %s: %v\n", path, err)
		}
	}

	if config.dirOnly {
		fmt.Printf("%d Directories\n", dirCounter)
	} else {
		fmt.Printf("%d Files, %d Directories\n", fileCounter, dirCounter)
	}
}

func tree(root, indent string, currLevel int) error {
	fileInfo, err := os.Stat(root)
	if err != nil {
		return fmt.Errorf("could not stat %s: %v", root, err)
	}

	fmt.Println(fileInfo.Name())
	if !fileInfo.IsDir() {
		fileCounter++
		return nil
	}

	dirCounter++

	fileInfos, err := ioutil.ReadDir(root)
	if err != nil {
		return fmt.Errorf("could not read dir %s: %v", root, err)
	}

	var names []string
	for _, fileInfo := range fileInfos {
		isDotFile := fileInfo.Name()[0] == '.'
		if config.dirOnly && fileInfo.IsDir() && !isDotFile {
			names = append(names, fileInfo.Name())
		}
		if !config.dirOnly && !isDotFile {
			names = append(names, fileInfo.Name())
		}
		if config.showDotFiles && isDotFile {
			names = append(names, fileInfo.Name())
		}
	}

	if currLevel+1 <= config.level {
		for i, name := range names {
			add := "│  "
			if i == len(names)-1 {
				fmt.Printf(indent + "└──")
				add = "   "
			} else {
				fmt.Printf(indent + "├──")
			}

			if err := tree(filepath.Join(root, name), indent+add, currLevel+1); err != nil {
				return err
			}
		}
	}

	return nil
}
