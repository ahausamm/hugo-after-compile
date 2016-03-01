package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	searchDir := "./rendered"

	fileList := []string{}
	filesToMinify := ".html,.htm,.css,.js"
	// Create a list of all files
	err := filepath.Walk(searchDir, func(path string, f os.FileInfo, err error) error {
		if !f.IsDir() && f.Name() != ".DS_Store" {
			fmt.Println(strings.ToLower(filepath.Ext(path)))
			//fileList = append(fileList, path)
		} else if f.Name() == ".DS_Store" {
			// Remove ds_store
			os.Remove(path)
		}
		return nil
	})

	if err != nil {
		fmt.Println(err)
		return
	} else {
		for _, file := range fileList {
			fmt.Println(file)
		}
	}
}
