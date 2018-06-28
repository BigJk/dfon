package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/BigJk/dfon"
)

func main() {
	dir := flag.String("dir", "./", "Directory to scan")
	search := flag.String("search", "", "Text to search")
	flag.Parse()

	*search = strings.ToLower(*search)

	if len(*search) == 0 {
		panic("Please specify a text to search\nExample: rawfinde -search=\"SIEGE_POP\"")
	}

	err := filepath.Walk(*dir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) == ".txt" {
			file, err := os.Open(path)
			if err != nil {
				return nil
			}
			defer file.Close()

			head, err := dfon.Parse(file)
			if err != nil {
				return nil
			}

			searchRecursive(head.Objects, make([]*dfon.Object, 0), regexp.MustCompile(*search), path)
		}
		return nil
	})

	if err != nil {
		fmt.Printf("walk error [%v]\n", err)
	}
}

func searchRecursive(data []*dfon.Object, parents []*dfon.Object, search *regexp.Regexp, origin string) {
	for i := range data {
		if search.MatchString(strings.ToLower(data[i].Type)) {
			fmt.Println("/=========================>")
			fmt.Println("| Path:   ", origin)
			fmt.Print("| ID:      ")
			for j := range parents {
				fmt.Print(parents[j].Type, " -> ")
			}
			fmt.Println(data[i].Type)

			if len(data[i].Traits) > 0 {
				fmt.Print("| Traits:  ")
				for j := range data[i].Traits {
					fmt.Print(data[i].Traits[j])
				}
				fmt.Print("\n")
			}

			if len(data[i].Values) > 0 {
				fmt.Println("| Values: ", strings.Join(data[i].Values, ", "))
			}
			fmt.Println("\\=========================>\n")
		}

		if data[i].Children != nil {
			searchRecursive(data[i].Children, append(parents, data[i]), search, origin)
		}
	}
}
