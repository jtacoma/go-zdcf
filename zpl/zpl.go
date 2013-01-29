package zpl

import (
	"bytes"
	"fmt"
	"regexp"
)

type Section struct {
	Properties map[string][]interface{}
	Sections   map[string]*Section
}

func NewSection() *Section {
	return &Section{
		Properties: make(map[string][]interface{}),
		Sections:   make(map[string]*Section),
	}
}

var (
	reskip       = regexp.MustCompile(`^\s*(#.*)?$`)
	reskipinline = regexp.MustCompile(`\s*(#.*)?$`)
	rekey        = regexp.MustCompile(
		`^(?P<indent>(    )*)(?P<key>[a-zA-Z0-9][a-zA-Z0-9/]*)(\s*(?P<hasvalue>=)\s*(?P<value>[^\s].*)\s*)?$`)
)

func splitLines(blob []byte) [][]byte {
	return bytes.FieldsFunc(blob, func(r rune) bool {
		return r == 10 || r == 13
	})
}

func Unmarshal(src []byte, dst interface{}) error {
	switch dst.(type) {
	case *Section:
	default:
		return fmt.Errorf("unsupported destination type: %T", dst)
	}
	ancestry := []*Section{dst.(*Section)}
	for lineno, line := range splitLines(src) {
		if inline := bytes.IndexByte(line, '#'); inline >= 0 {
			line = line[:inline]
		}
		//if skip:=reskipinline.Find(line);skip!=nil{ line=line[:skip[0]] }
		line = bytes.TrimRight(line, " \t\n\r")
		if len(line) == 0 || reskip.Match(line) {
			continue
		} else if match := rekey.FindSubmatch(line); match != nil {
			depth := len(match[1]) / 4
			if depth+1 < len(ancestry) {
				ancestry = ancestry[:depth+1]
			}
			section := ancestry[len(ancestry)-1]
			key := string(match[3])
			if len(match[5]) > 0 {
				value := string(match[6])
				if array, exists := section.Properties[key]; exists {
					section.Properties[key] = append(array, value)
				} else {
					section.Properties[key] = []interface{}{value}
				}
			} else {
				if _, exists := section.Sections[key]; exists {
					// TODO: this shouldn't really be a problem...
					return fmt.Errorf("line %d: duplicate subsection %s", lineno, key)
				} else {
					section.Sections[key] = NewSection()
					ancestry = append(ancestry, section.Sections[key])
				}
			}
		} else {
			return fmt.Errorf("line %d: invalid ZPL: %v", lineno, string(line))
		}
	}
	return nil
}
