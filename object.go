package dfon

import (
	"bytes"
	"strings"
)

type Object struct {
	Type     string
	Values   []string
	Children []*Object
	Traits   []*Object
}

func (o *Object) IsFlag() bool {
	return len(o.Values) == 0
}

func (o *Object) String() string {
	buf := bytes.Buffer{}
	buf.WriteString("[")
	buf.WriteString(o.Type)
	if !o.IsFlag() {
		buf.WriteString(":")
		buf.WriteString(strings.Join(o.Values, ":"))
	}
	buf.WriteString("]")
	for i := range o.Traits {
		buf.WriteString(o.Traits[i].String())
	}
	return buf.String()
}
