package src

import (
	"fmt"
	"log"
	"os"
	"strings"
)

type WriteOption struct {
	FileName string
	TypeName string
	Fields   []Field
}

type Field struct {
	OmitEmpty bool
	Repeated  bool
	Type      string
	Name      string
	Serial    int
}

type Protobuf struct {
	opt      *WriteOption
	builder  *strings.Builder
	filepath string
}

func NewProtobufWriter(filepath string) *Protobuf {
	p := Protobuf{
		filepath: ReplaceFileExt(filepath, ".proto"),
	}
	header := "syntax = \"proto3\";\n\n"
	var builder strings.Builder
	builder.WriteString(header)
	p.builder = &builder
	return &p
}

// Write uses strings.Builder to WriteString from given WriteOption
func (p *Protobuf) Write(w *WriteOption) {
	if w == nil {
		return
	}
	p.opt = w
	message := fmt.Sprintf("message %s {\n", p.opt.TypeName)
	end := "}\n\n"
	p.builder.WriteString(message)
	for i := range p.opt.Fields {
		protoType := ConvertTypeToProto(p.opt.Fields[i].Type)
		p.builder.WriteString("  ")
		if p.opt.Fields[i].Repeated {
			p.builder.WriteString("repeated ")
		}
		p.builder.WriteString(fmt.Sprintf("%s %s = %d;\n", protoType, p.opt.Fields[i].Name, i+1))
	}
	p.builder.WriteString(end)
}

// WriteFile saves data in strings.Builder to file.
func (p *Protobuf) WriteFile() {
	err := os.WriteFile(p.filepath, []byte(p.builder.String()), 0644)
	if err != nil {
		log.Printf("Failed to write file: %v", err)
	}
}

type Graphql struct {
	opt      *WriteOption
	builder  *strings.Builder
	filepath string
}

func NewGraphqlWriter(filepath string) *Graphql {
	g := Graphql{
		filepath: ReplaceFileExt(filepath, ".graphqls"),
	}
	var builder strings.Builder
	builder.WriteString("")
	g.builder = &builder
	return &g
}

// Write uses strings.Builder to WriteString from given WriteOption
func (g *Graphql) Write(w *WriteOption) {
	if w == nil {
		return
	}
	g.opt = w
	message := fmt.Sprintf("type %s {\n", g.opt.TypeName)
	end := "}\n\n"
	g.builder.WriteString(message)
	for i := range g.opt.Fields {
		protoType := ConvertTypeToGraphQL(g.opt.Fields[i].Type)
		g.builder.WriteString(fmt.Sprintf("  %s: ", w.Fields[i].Name))
		if g.opt.Fields[i].Repeated {
			g.builder.WriteString("[")
		}
		g.builder.WriteString(protoType)
		if g.opt.Fields[i].Repeated {
			g.builder.WriteString("]")
		}
		if w.Fields[i].OmitEmpty {
			g.builder.WriteString("\n")
			continue
		}
		g.builder.WriteString("!\n")
	}
	g.builder.WriteString(end)
}

// WriteFile saves data in strings.Builder to file.
func (g *Graphql) WriteFile() {
	err := os.WriteFile(g.filepath, []byte(g.builder.String()), 0644)
	if err != nil {
		log.Printf("Failed to write file: %v", err)
	}
}
