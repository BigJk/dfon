package dfon

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

type Object struct {
	Enabled  bool      `json:"enabled"`
	Type     string    `json:"type"`
	Values   []string  `json:"values"`
	Children []*Object `json:"children"`
	Traits   []*Object `json:"traits"`
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

func (o *Object) GetById(id string) []*Object {
	return getById(id, o.Children)
}

func (o *Object) AsBool() bool {
	if len(o.Values) > 0 && o.Values[0] == "YES" {
		return true
	}
	return false
}

func (o *Object) SetBool(value bool) {
	if len(o.Values) == 0 {
		if value {
			o.Values = []string{"YES"}
		} else {
			o.Values = []string{"NO"}
		}
	} else {
		if value {
			o.Values[0] = "YES"
		} else {
			o.Values[0] = "NO"
		}
	}
}

func (o *Object) AsInt(index int) int {
	if index >= len(o.Values) {
		return -1
	}
	n, _ := strconv.Atoi(o.Values[index])
	return n
}

func (o *Object) SetInt(index int, value int) {
	if index < len(o.Values) {
		o.Values[index] = fmt.Sprint(value)
	}
}
