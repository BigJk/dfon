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
		if data[i].IsFlag() {
			writer.Write([]byte(fmt.Sprint(strings.Repeat("\t", depth), "[", data[i].Type, "]\r\n")))
		} else {
			writer.Write([]byte(fmt.Sprint(strings.Repeat("\t", depth), "[", data[i].Type, ":", strings.Join(data[i].Values, ":"), "]\r\n")))
		}
		if data[i].Children != nil {
			h.printRecursive(writer, data[i].Children, depth+1)
		}
	}
}
