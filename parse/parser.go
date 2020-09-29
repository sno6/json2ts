package parse

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"text/scanner"
)

type Parser struct {
	s scanner.Scanner
}

type Node struct {
	Name     string
	Type     Type
	ObjDepth int
	ArrDepth int
}

type Type uint

const (
	Object Type = iota
	Array
	String
	Number
	Bool
	Null
)

func (p Parser) Parse(input string) ([]*Node, error) {
	r, err := p.readInput(input)
	if err != nil {
		return nil, err
	}

	var s scanner.Scanner
	s.Init(r)

	var (
		objDepth int
		arrDepth int
	)

	var nodes []*Node
	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		if s.ErrorCount > 0 {
			return nil, errors.New("parse error")
		}

		token := s.TokenText()

		if token == "{" {
			objDepth++
		}
		if token == "}" {
			objDepth--
		}

		if token == "[" {
			arrDepth++
		}
		if token == "]" {
			arrDepth--
		}

		if isKey(s, token) {
			// Skip ":"
			s.Next()

			// Skip " "
			s.Next()

			// This is the important rune that tells us what type of node we are about to see.
			r := s.Peek()
			typ := tellNodeType(r)

			nodes = append(nodes, &Node{
				Name:     token,
				Type:     typ,
				ArrDepth: arrDepth,
				ObjDepth: objDepth - 1,
			})
		}
	}

	return nodes, nil
}

func isKey(s scanner.Scanner, tok string) bool {
	return strings.HasPrefix(tok, `"`) && strings.HasSuffix(tok, `"`) && s.Peek() == ':'
}

func tellNodeType(r rune) Type {
	switch r {
	case '{':
		return Object
	case '[':
		return Array
	case '"':
		return String
	case 't', 'f':
		return Bool
	case 'n':
		return Null
	default:
		return Number
	}
}

func (t Type) String() string {
	switch t {
	case Object:
		return "object"
	case Array:
		return "array"
	case Number:
		return "number"
	case String:
		return "string"
	case Bool:
		return "boolean"
	case Null:
		return "null"
	}

	panic("unknown node type")
}

func (p Parser) readInput(input string) (io.Reader, error) {
	var r io.Reader

	if input == "" {
		r = os.Stdin
	} else {
		f, err := os.Open(input)
		if err != nil {
			return nil, err
		}
		defer f.Close()

		r = f
	}

	return p.indentedJSONReader(r)
}

func (p Parser) indentedJSONReader(r io.Reader) (io.Reader, error) {
	j, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	b := &bytes.Buffer{}
	if err := json.Indent(b, j, "", "\t"); err != nil {
		return nil, err
	}

	return b, nil
}
