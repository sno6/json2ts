package transform

import (
	"fmt"
	"io"
	"json2ts/parse"
	"os"
	"strings"
)

const DefaultBaseClassName = "BaseClass"

type Transformer struct{}

type Config struct {
	BaseClassName   string
	PrefixClassName string
	Decorators      bool
	Output          string
}

type Class struct {
	Name       string
	Attributes []*AttributeValue
}

type AttributeValue struct {
	Key   string
	Type  parse.Type
	Class *Class
}

func (t Transformer) Transform(nodes []*parse.Node, cfg *Config) error {
	tree := t.BuildClassTree(nodes, cfg)
	tmpl := t.buildTemplate(tree, cfg.Decorators)
	_, err := io.WriteString(os.Stdout, tmpl)

	return err
}

func (t Transformer) BuildClassTree(nodes []*parse.Node, cfg *Config) *Class {
	if cfg.BaseClassName == "" {
		cfg.BaseClassName = DefaultBaseClassName
	}

	// Create the root node of the class tree structure.
	tree := &Class{Name: cfg.BaseClassName}

	// Cache to keep track of nodes at N depth.
	parentCache := []*Class{tree}

	for i, node := range nodes {
		// If we are moving more shallow, pop N nodes from the cache to align with current parent.
		if i > 0 && node.ObjDepth < nodes[i-1].ObjDepth {
			nPop := nodes[i-1].ObjDepth - node.ObjDepth
			parentCache = parentCache[0 : len(parentCache)-nPop]
		}

		// The parent is always the node at the head of the cache.
		parent := parentCache[len(parentCache)-1]

		// Create a new attribute against the current parent.
		value := &AttributeValue{
			Key:  strings.Replace(node.Name, `"`, ``, -1),
			Type: node.Type,
		}

		nextNodeInArray := i < len(nodes)-1 && nodes[i+1] != nil && nodes[i+1].ArrDepth == node.ArrDepth+1
		if node.Type == parse.Object || (node.Type == parse.Array && nextNodeInArray) {
			class := &Class{Name: convertClassName(node.Name)}
			parentCache = append(parentCache, class)
			value.Class = class
		}

		parent.Attributes = append(parent.Attributes, value)
	}

	return tree
}

func (t Transformer) buildTemplate(tree *Class, addDecorators bool) string {
	var (
		s strings.Builder

		// Keep track of the classes we encounter in order to recursively build the class templates.
		classes []*Class
	)

	s.WriteString(fmt.Sprintf("export class %s {\n", tree.Name))

	for i, v := range tree.Attributes {
		var (
			isArray = false
			t       = v.Type.String()
		)

		if v.Type == parse.Array {
			isArray = true
		}

		if v.Type == parse.Array && v.Class == nil {
			// TODO ...
			t = "any"
		}

		if v.Type == parse.Null {
			t = "any"
		}

		if (v.Type == parse.Object || v.Type == parse.Array) && v.Class != nil {
			t = v.Class.Name
			classes = append(classes, v.Class)
		}

		if addDecorators {
			decorators := generateDecorators(v)
			for _, d := range decorators {
				str := fmt.Sprintf("\t %s \n", d)
				s.WriteString(str)
			}
		}

		arrayPrefix := ""
		if isArray {
			arrayPrefix = "[]"
		}

		str := fmt.Sprintf("\tpublic %s!: %s%s;\n", v.Key, t, arrayPrefix)
		s.WriteString(str)

		if i != len(tree.Attributes)-1 {
			s.WriteString("\n")
		}
	}

	s.WriteString("}\n\n")

	// Recursively generate class templates for all the classes we encountered in the current class..
	for _, v := range classes {
		s.WriteString(t.buildTemplate(v, addDecorators))
	}

	return s.String()
}

func generateDecorators(v *AttributeValue) []string {
	// Generator for class decorators.
	tGen := func(decs []string, v *AttributeValue) {
		decs = append(decs, fmt.Sprintf("@Type(() => %s)", v.Class.Name))
		decs = append(decs, "@ValidateNested({ each: true })")
	}

	var decorators []string

	switch v.Type {
	case parse.Object:
		tGen(decorators, v)
		break
	case parse.Array:
		decorators = append(decorators, "@IsArray()")
		if v.Class != nil {
			tGen(decorators, v)
		}
		break
	case parse.Bool:
		decorators = append(decorators, "@IsBoolean()")
		break
	case parse.String:
		decorators = append(decorators, "@IsString()")
		break
	case parse.Number:
		decorators = append(decorators, "@IsNumber()")
		break
	case parse.Null:
		decorators = append(decorators, "@IsOptional()")
		break

	}

	return decorators
}

func convertClassName(s string) string {
	s = strings.Replace(s, `"`, ``, -1)
	s = strings.Replace(s, `-`, ``, -1)
	s = strings.Replace(s, `_`, ``, -1)
	return strings.Title(s)
}
