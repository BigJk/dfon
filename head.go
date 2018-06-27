package dfon

import (
	"fmt"
	"io"
	"strings"
)

type Head struct {
	Name    string
	Type    string
	Objects []*Object
}

func (h *Head) Print(writer io.Writer) {
	writer.Write([]byte(h.Name + "\n\n"))
	writer.Write([]byte("[OBJECT:" + h.Type + "]" + "\n\n"))
	h.printRecursive(writer, h.Objects, 0)
}

func (h *Head) printRecursive(writer io.Writer, data []*Object, depth int) {
	for i := range data {
		writer.Write([]byte(fmt.Sprint(strings.Repeat("\t", depth), data[i].String(), "\r\n")))
		if data[i].Children != nil {
			h.printRecursive(writer, data[i].Children, depth+1)
		}
		if depth == 0 {
			writer.Write([]byte("\n\r"))
		}
	}
}
