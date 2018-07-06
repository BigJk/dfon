package dfon

import (
	"bufio"
	"io"
	"regexp"

	"bytes"
	"errors"
	"strings"
)

var (
	filenameForbidden     = regexp.MustCompile("[^a-zA-Z_]")
	objectDefinitionRegex = regexp.MustCompile("\\[OBJECT:(.+?)\\]")
	objectRegex           = regexp.MustCompile("[\\[╳](.*?)[╳\\]]")

	ErrDataEndedTooSoon  = errors.New("data ended too soon")
	ErrNoCorrectObject   = errors.New("correct object definition not found")
	ErrNoObjectFound     = errors.New("no object found")
	ErrEmpty             = errors.New("empty")
	ErrBaseObjectMissing = errors.New("base object missing")
)

func ParseBytes(data []byte) (*Head, error) {
	return Parse(bytes.NewBuffer(data))
}

func ParseString(data string) (*Head, error) {
	return Parse(bytes.NewBufferString(data))
}

func Parse(stream io.Reader) (*Head, error) {
	reader := bufio.NewScanner(stream)
	skipObjectDefinition := false

	var head Head
	head.HasHeader = true

	// parse filename if existent
	if !reader.Scan() {
		return nil, ErrDataEndedTooSoon
	}

	// check if filename is present and no comment
	head.Name = strings.Trim(reader.Text(), " \t")
	if len(head.Name) == 0 || filenameForbidden.MatchString(head.Name) || objectRegex.MatchString(head.Name) {
		head.Name = ""
		head.HasHeader = false
		skipObjectDefinition = true
	}

	initial := false
	head.Objects = make([]*Object, 0)

	// parse object definition
	if !skipObjectDefinition {
		for {
			if !reader.Scan() {
				return nil, ErrDataEndedTooSoon
			}
			if objectDefinitionRegex.MatchString(reader.Text()) {
				matches := objectDefinitionRegex.FindStringSubmatch(reader.Text())
				if len(matches) != 2 {
					return nil, ErrNoCorrectObject
				}
				head.Type = matches[1]
				break
			}
		}
	} else if objectRegex.MatchString(reader.Text()) {
		initial = true
	}

	// extract sections and parse
	sections := sections(reader, initial)
	for i := range sections {
		if section, err := buildSection(sections[i]); err == nil {
			head.Objects = append(head.Objects, section)
		}
	}

	if len(head.Objects) == 0 && len(head.Name) == 0 && len(head.Type) == 0 {
		return nil, ErrEmpty
	}

	return &head, nil
}

func sections(reader *bufio.Scanner, initial bool) []string {
	var sections []string

	var cur string
	for initial || reader.Scan() && reader.Err() == nil {
		initial = false

		if len(reader.Text()) == 0 || !objectRegex.MatchString(reader.Text()) {
			continue
		}

		if len(cur) > 0 && strings.Count(reader.Text(), "\t") == 0 {
			sections = append(sections, cur)
			cur = ""
		}

		cur += reader.Text() + "\r\n"
	}
	sections = append(sections, cur)

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
				for i := 0; i < depth-1; i++ {
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
		return nil, ErrBaseObjectMissing
	}

	return base, nil
}

func buildObjectAndTraits(text string) (*Object, error) {
	objects := objectRegex.FindAllStringSubmatch(text, -1)
	if len(objects) == 0 {
		return nil, ErrNoObjectFound
	}
	object := parseObject(objects[0][1], strings.HasPrefix(objects[0][0], "[") && strings.HasSuffix(objects[0][0], "]"))
	if len(objects) > 1 {
		for i := 1; i < len(objects); i++ {
			trait := parseObject(objects[i][1], strings.HasPrefix(objects[i][0], "[") && strings.HasSuffix(objects[i][0], "]"))
			if trait == nil {
				continue
			}
			object.Traits = append(object.Traits, trait)
		}
	}
	return object, nil
}

func parseObject(content string, enabled bool) *Object {
	values := strings.Split(content, ":")
	if len(values) == 0 {
		return nil
	}

	var newObject Object
	newObject.Enabled = enabled
	newObject.Type = values[0]
	if len(values) > 1 {
		newObject.Values = values[1:]
	}

	return &newObject
}
