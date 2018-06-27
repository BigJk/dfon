package dfon

import (
	"bufio"
	"fmt"
	"io"
	"regexp"

	"strings"
)

var (
	objectDefinitionRegex = regexp.MustCompile("\\[OBJECT:(.+)\\]")
	objectRegex           = regexp.MustCompile("\\[(.*)\\]")
)

func Parse(stream io.Reader) (*Head, error) {
	reader := bufio.NewScanner(stream)

	var head Head

	// parse filename
	if !reader.Scan() {
		return nil, fmt.Errorf("filename missing")
	}

	head.Name = reader.Text()
	head.Objects = make([]*Object, 0)

	// parse object definition
	for {
		if !reader.Scan() {
			return nil, fmt.Errorf("data ended too soon")
		}
		if objectDefinitionRegex.MatchString(reader.Text()) {
			matches := objectDefinitionRegex.FindStringSubmatch(reader.Text())
			if len(matches) != 2 {
				return nil, fmt.Errorf("correct object defintion not found")
			}
			head.Type = matches[1]
			break
		}
	}

	// parse objects
	fill(reader, 0, &head.Objects, nil)

	return &head, nil
}

func fill(reader *bufio.Scanner, depth int, target *[]*Object, lastTarget *[]*Object) {
	for reader.Scan() {
		var lastObject *Object
		if target != nil && *target != nil && len(*target) > 0 {
			lastObject = (*target)[len(*target)-1]
		}

		if objectRegex.MatchString(reader.Text()) {
			tabs := strings.Count(reader.Text(), "\t")
			if tabs < depth {
				if lastTarget != nil {
					objects := objectRegex.FindAllStringSubmatch(reader.Text(), -1)
					for i := range objects {
						object := parseObject(objects[i][1])
						if object == nil {
							continue
						}

						*lastTarget = append(*lastTarget, object)
					}
				}
				return
			} else if tabs > depth && target != nil {
				if lastObject.Children == nil {
					lastObject.Children = make([]*Object, 0)
				}

				nextTarget := &lastObject.Children
				objects := objectRegex.FindAllStringSubmatch(reader.Text(), -1)
				for i := range objects {
					object := parseObject(objects[i][1])
					if object == nil {
						continue
					}

					*nextTarget = append(*nextTarget, object)
				}

				fill(reader, depth+1, nextTarget, target)
			} else {
				objects := objectRegex.FindAllStringSubmatch(reader.Text(), -1)
				for i := range objects {
					object := parseObject(objects[i][1])
					if object == nil {
						continue
					}

					*target = append(*target, object)
				}
			}
		}
	}
}

func parseObject(content string) *Object {
	values := strings.Split(content, ":")
	if len(values) == 0 {
		return nil
	}

	var newObject Object
	newObject.Type = values[0]
	if len(values) > 1 {
		newObject.Values = values[1:]
	}

	return &newObject
}
