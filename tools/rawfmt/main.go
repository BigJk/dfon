package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"bytes"
	"io/ioutil"

	"github.com/BigJk/dfon"
)

func main() {
	dir := flag.String("dir", "./", "Directory to format")
	flag.Parse()

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

			var buf bytes.Buffer
			head.Print(&buf)
			ioutil.WriteFile(path, buf.Bytes(), 0777)

			fmt.Println(path, "formated")
		}
		return nil
	})

	if err != nil {
		fmt.Printf("walk error [%v]\n", err)
	}
}
