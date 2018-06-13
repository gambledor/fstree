package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func main() {
	args := []string{"."}
	if len(os.Args) > 1 {
		args = os.Args[1:]
	}
	for _, path := range args {
		if err := tree(path, ""); err != nil {
			log.Printf("tree %s: %v\n", path, err)
		}
	}
}

func tree(root, indent string) error {
	fileInfo, err := os.Stat(root)
	if err != nil {
		return fmt.Errorf("could not stat %s: %v", root, err)
	}

	fmt.Println(fileInfo.Name())
	if !fileInfo.IsDir() {
		return nil
	}

	fileInfos, err := ioutil.ReadDir(root)
	if err != nil {
		return fmt.Errorf("could not read dir %s: %v", root, err)
	}

	var names []string
	for _, fileInfo := range fileInfos {
		if fileInfo.Name()[0] != '.' {
			names = append(names, fileInfo.Name())
		}
	}

	for i, name := range names {
		add := "│  "
		if i == len(names)-1 {
			fmt.Printf(indent + "└──")
			add = "   "
		} else {
			fmt.Printf(indent + "├──")
		}

		if err := tree(filepath.Join(root, name), indent+add); err != nil {
			return err
		}
	}

	return nil
}
