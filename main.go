package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
)

type ByLength []string

func (s ByLength) Len() int {
	return len(s)
}
func (s ByLength) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s ByLength) Less(i, j int) bool {
	return len(s[i]) < len(s[j])
}

func IsUsefulExtension(extension string) bool {
	switch extension {
	case
		".html",
		".htm":
		return true
	}
	return false
}

func IsDirEmpty(name string) (bool, error) {
	f, err := os.Open(name)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdir(1)

	if err == io.EOF {
		return true, nil
	}
	return false, err
}

func main() {
	searchDir := flag.String("s", "./rendered", "a string with the path")
	flag.Parse()
	fileList := []string{}

	if _, err := os.Stat(*searchDir); err == nil {
		if err != nil {
			log.Fatalln(err)
		}
		// Create a list of all files
		err := filepath.Walk(*searchDir, func(path string, f os.FileInfo, err error) error {
			if !f.IsDir() && f.Name() != ".DS_Store" {
				if IsUsefulExtension(strings.ToLower(filepath.Ext(path))) {
					fileList = append(fileList, path)
				}
			} else if f.Name() == ".DS_Store" {
				// Remove ds_store
				os.Remove(path)
			}
			return nil
		})
		//Sort slice (longest path first)
		sort.Sort(sort.Reverse(ByLength(fileList)))

		if err != nil {
			log.Fatalln(err)
			return
		} else {
			for _, file := range fileList {
				input, err := ioutil.ReadFile(file)
				if err != nil {
					log.Fatalln(err)
				}
				inputAsString := string(input)
				if err != nil {
					log.Fatalln(err)
				} else {
					if len(inputAsString) > 0 {
						// search and replace
						lines := strings.Split(inputAsString, "\n")

						for i, line := range lines {
							if strings.Contains(line, "<p><section") {
								lines[i] = strings.Replace(line, "<p><section", "<section", 1)
							}
							if strings.Contains(line, "section></p>") {
								lines[i] = strings.Replace(line, "section></p>", "section>", 1)
							}
						}
						output := strings.Join(lines, "\n")
						fileInfo, _ := os.Stat(file)
						mode := fileInfo.Mode()
						err = ioutil.WriteFile(file, []byte(output), mode)
						if err != nil {
							log.Fatalln(err)
						}
					} else {
						// Delete empty file
						os.Remove(file)
						dirPath := path.Dir(file)
						checkedDir, _ := IsDirEmpty(dirPath)
						// if necessary delete the empty directory
						if checkedDir {
							os.Remove(dirPath)
						}
					}
				}
			}
		}
	} else {
		fmt.Println("Error: The given directory doesn't exist")
	}
}
