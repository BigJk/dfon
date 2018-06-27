package dfon

type Object struct {
	Type     string
	Values   []string
	Children []*Object
}

func (o *Object) IsFlag() bool {
	return len(o.Values) == 0
}
