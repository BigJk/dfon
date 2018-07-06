package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"bufio"

	"github.com/BigJk/dfon"
)

func main() {
	dir := flag.String("dir", "./", "Directory to scan")
	filemask := flag.String("file", "*", "Filemask to scan for")
	object := flag.String("obj", "", "Object to modify")
	disable := flag.Bool("disable", false, "If object should be disabled")
	enable := flag.Bool("enable", false, "If object should be enabled")
	newValue := flag.String("val", "", "New object value")
	noprompt := flag.Bool("no-prompt", false, "Skip prompt before making changes")
	flag.Parse()

	if *object == "" {
		flag.Usage()
		return
	}

	changes := make(map[string]*dfon.Head)

	// walk files
	filepath.Walk(*dir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if match, err := filepath.Match(*filemask, filepath.Base(path)); match && err == nil && filepath.Ext(path) == ".txt" {
			file, err := os.Open(path)
			if err != nil {
				return nil
			}
			defer file.Close()

			head, err := dfon.Parse(file)
			if err != nil {
				return nil
			}

			found := searchRecursive(head.Objects, *object)
			if len(found) > 0 {
				fmt.Println(">", path)
				for i := range found {
					fmt.Print(i+1, "# ", found[i].String(), " -> ")

					if *disable {
						found[i].Enabled = false
					} else if *enable {
						found[i].Enabled = true
					}

					if len(*newValue) > 0 {
						found[i].Values = strings.Split(*newValue, ",")
					}

					fmt.Print(found[i].String(), "\n")
				}
				fmt.Print("\n")

				changes[path] = head
			}
		}
		return nil
	})

	// check if something was found
	if len(changes) == 0 {
		fmt.Println("nothing found")
		return
	}

	// prompt
	if !*noprompt {
		fmt.Print("Want to commit changes? (Y/n) ")
		reader := bufio.NewScanner(os.Stdin)

	inputLoop:
		for {
			if reader.Scan() {
				switch strings.ToLower(reader.Text()) {
				case "y":
					break inputLoop
				case "n":
					fmt.Println("Aborted...")
					return
				default:
					continue
				}
			} else {
				panic("failed to read input")
			}
		}

		fmt.Println()
	}

	// write changes to files
	for path, head := range changes {
		file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
		if err != nil {
			fmt.Errorf("error while writing to file %s %v", path, err)
			continue
		}
		head.Print(file)
		file.Close()
		fmt.Println("Wrote to", path)
	}

	fmt.Println("\nFinished")
}

func searchRecursive(data []*dfon.Object, object string) []*dfon.Object {
	var results []*dfon.Object
	for i := range data {
		if match, err := filepath.Match(object, data[i].Type); match && err == nil {
			results = append(results, data[i])
		}

		if data[i].Traits != nil {
			results = append(results, searchRecursive(data[i].Traits, object)...)
		}

		if data[i].Children != nil {
			results = append(results, searchRecursive(data[i].Children, object)...)
		}
	}
	return results
}
