package dfon

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

type Head struct {
	HasHeader bool      `json:"hasHeader"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	Objects   []*Object `json:"objects"`
}

func (h *Head) Print(writer io.Writer) {
	if h.HasHeader {
		writer.Write([]byte(h.Name + "\n\n"))
		writer.Write([]byte("[OBJECT:" + h.Type + "]" + "\n\n"))
	}
	h.printRecursive(writer, h.Objects, 0)
}

func (h *Head) printRecursive(writer io.Writer, data []*Object, depth int) {
	for i := range data {
		writer.Write([]byte(fmt.Sprint(strings.Repeat("\t", depth), data[i].String(), "\r\n")))
		if data[i].Children != nil {
			h.printRecursive(writer, data[i].Children, depth+1)
		}
		if depth == 0 && h.HasHeader {
			writer.Write([]byte("\n\r"))
		}
	}
}

func (h *Head) String() string {
	var buffer bytes.Buffer
	h.Print(&buffer)
	return buffer.String()
}

func (h *Head) GetById(id string) []*Object {
	return getById(id, h.Objects)
}

func getById(id string, objects []*Object) []*Object {
	var found []*Object
	for i := range objects {
		if objects[i].Type == id {
			found = append(found, objects[i])
		}
		found = append(found, getById(id, objects[i].Children)...)
		found = append(found, getById(id, objects[i].Traits)...)
	}
	return found
}
