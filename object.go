package dfon

import (
	"bytes"
	"strings"
)

type Object struct {
	Enabled  bool
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
	if o.Enabled {
		buf.WriteString("[")
	} else {
		buf.WriteString("╳")
	}
	buf.WriteString(o.Type)
	if !o.IsFlag() {
		buf.WriteString(":")
		buf.WriteString(strings.Join(o.Values, ":"))
	}
	if o.Enabled {
		buf.WriteString("]")
	} else {
		buf.WriteString("╳")
	}
	for i := range o.Traits {
		buf.WriteString(o.Traits[i].String())
	}
	return buf.String()
}

func (o *Object) EnableState(state, traits bool) {
	o.Enabled = state
	if traits {
		for i := range o.Traits {
			o.Traits[i].Enabled = state
		}
	}
}
