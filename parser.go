package dfon

import (
	"bufio"
	"fmt"
	"io"
	"regexp"

	"bytes"
	"strings"
)

var (
	objectDefinitionRegex = regexp.MustCompile("\\[OBJECT:(.+?)\\]")
	objectRegex           = regexp.MustCompile("\\[(.*?)\\]")
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

	sections := sections(reader)
	for i := range sections {
		if section, err := buildSection(sections[i]); err == nil {
			head.Objects = append(head.Objects, section)
		}
	}

	return &head, nil
}

func sections(reader *bufio.Scanner) []string {
	var sections []string

	var cur string
	for reader.Scan() && reader.Err() == nil {
		if len(reader.Text()) == 0 || !objectRegex.MatchString(reader.Text()) {
			continue
		}
		if len(cur) > 0 && strings.Count(reader.Text(), "\t") == 0 {
			sections = append(sections, cur)
			cur = ""
		}
		cur += reader.Text() + "\r\n"
	}

	return sections
}

func buildSection(section string) (*Object, error) {
	var base *Object

	reader := bufio.NewScanner(bytes.NewBufferString(section))
	for reader.Scan() {
		if scanned, err := buildObjectAndTraits(reader.Text()); err == nil {
			if base == nil {
				base = scanned
			} else {
				depth := strings.Count(reader.Text(), "\t")

				target := base
				for i := 0; i < depth; i++ {
					if target.Children == nil || len(target.Children) == 0 {
						break
					}
					target = target.Children[len(target.Children)-1]
				}

				target.Children = append(target.Children, scanned)
			}
		}
	}

	if base == nil {
		return nil, fmt.Errorf("base object missing")
	}

	return base, nil
}

func buildObjectAndTraits(text string) (*Object, error) {
	objects := objectRegex.FindAllStringSubmatch(text, -1)
	if len(objects) == 0 {
		return nil, fmt.Errorf("no object found")
	}
	object := parseObject(objects[0][1])
	if len(objects) > 1 {
		for i := 1; i < len(objects); i++ {
			trait := parseObject(objects[i][1])
			if trait == nil {
				continue
			}
			object.Traits = append(object.Traits, trait)
		}
	}
	return object, nil
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
