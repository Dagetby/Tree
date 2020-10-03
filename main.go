package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type treeState struct {
}

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}

func dirTree(out io.Writer, filePath string, printFiles bool) error {
	printDirTree(out, filePath, 0, printFiles, "")

	return nil
}

func printDirTree(out io.Writer, filePath string, level int, flag bool, prevTab string) {
	files := readDir(filePath) // возвращает []os.FileInfo
	var endOfFile bool
	folders := make([]os.FileInfo, 0)

	for _, f := range files {
		if f.IsDir() {
			folders = append(folders, f)
		}
	}
	var items []os.FileInfo
	if flag {
		items = files
	} else {
		items = folders
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].Name() < items[j].Name()
	})
	for i, value := range items {
		if flag {
			if i == (len(items) - 1) {
				if value.Name() != ".DS_Store" && value.Name() != ".idea" {
					endOfFile = true
					print(out, value, level, endOfFile, prevTab)
				}
			} else {
				if value.Name() != ".DS_Store" && value.Name() != ".idea" {
					print(out, value, level, endOfFile, prevTab)
				}
			}
		} else {
			if value.IsDir() {
				if i == (len(items) - 1) {
					if value.Name() != ".idea" {
						endOfFile = true
						print(out, value, level, endOfFile, prevTab)
					}
				} else {
					if value.Name() != ".idea" {
						print(out, value, level, endOfFile, prevTab)
					}
				}
			}

		}
		if !strings.Contains(value.Name(), ".idea") && value.IsDir() {
			var tempTab string
			if endOfFile {
				tempTab = prevTab + "\t"
			} else {
				tempTab = prevTab + "│\t"
			}

			path := filepath.Join(filePath, value.Name())
			printDirTree(out, path, level+1, flag, tempTab)
		}
	}
}
func readDir(filePath string) []os.FileInfo {
	file, err := os.Open(filePath) //Чтение директории
	if err != nil {
		log.Fatal(err)
	}
	fileInfo, err := file.Readdir(0)
	if err != nil {
		log.Fatalf("[ERROR]: %s/n", err)
	}

	return fileInfo
}

func print(out io.Writer, name os.FileInfo, level int, endOfFile bool, previosTab string) {
	open_c := "├───"
	close_c := "└───"
	var size string
	if !name.IsDir() {
		if name.Size() != 0 {
			size = fmt.Sprintf(" (%db)\n", name.Size())
		} else {
			size = " (empty)\n"
		}
	} else {
		size = "\n"
	}

	if endOfFile {
		fmt.Fprint(out, previosTab+close_c+name.Name()+size)

	} else {
		fmt.Fprint(out, previosTab+open_c+name.Name()+size)
	}

}
