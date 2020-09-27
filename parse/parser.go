package parse

import (
	"errors"
	"io"
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

func (p Parser) Parse(r io.Reader) ([]*Node, error) {
	var s scanner.Scanner
	s.Init(os.Stdin)

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
