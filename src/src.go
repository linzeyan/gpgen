package src

import (
	"go/ast"
	"go/parser"
	"go/token"
	"log"
)

type Options struct {
	Graphql    bool
	Protobuf   bool
	SrcFile    string
	DstFile    string
	StructName string
}

func Analysis(opt *Options) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, opt.SrcFile, nil, 0)
	// f, err := parser.ParseFile(fset, fileName, nil, 0)
	if err != nil {
		log.Fatalf("Failed to parse: %v", err)
	}
	// ast.Print(fset, f)

	dst := opt.SrcFile
	if opt.DstFile != "" {
		dst = opt.DstFile
	}
	p := NewProtobufWriter(dst)
	g := NewGraphqlWriter(dst)

	// store self define type, and is basic type
	tempBasicType := make(map[string]string)

	ast.Inspect(f, func(n ast.Node) bool {
		switch genDecl := n.(type) {
		case *ast.GenDecl:
			if genDecl.Tok != token.TYPE {
				return true
			}
			for i := range genDecl.Specs {
				switch spec := genDecl.Specs[i].(type) {
				case *ast.TypeSpec:
					switch field := spec.Type.(type) {
					case *ast.Ident:
						tempBasicType[spec.Name.Name] = field.String()
					case *ast.ArrayType:
						tempBasicType[spec.Name.Name] = field.Elt.(*ast.Ident).String()
					case *ast.StructType:
						if opt.StructName != "" && spec.Name.Name != opt.StructName {
							continue
						}
						writer := WriteOption{
							TypeName: spec.Name.Name,
							Fields:   make([]Field, len(field.Fields.List)),
						}
						for i, f := range field.Fields.List {
							// field name: f.Names[0].Name
							// field type: f.Type.(*ast.Ident).Name
							// field tag: f.Tag.Value
							// log.Println(f.Type.(*ast.Ident).Name, f.Tag.Value)
							var typ string
							var repeated bool
							switch t := f.Type.(type) {
							case *ast.Ident:
								typ = t.Name
								if basicType, ok := tempBasicType[typ]; ok {
									typ = basicType
								}
							case *ast.ArrayType:
								repeated = true
								typ = t.Elt.(*ast.Ident).Name
								if basicType, ok := tempBasicType[typ]; ok {
									typ = basicType
								}
							}
							tag, empty := GetStructTag(f.Tag.Value)
							writer.Fields[i] = Field{
								Type:      typ,
								Name:      tag,
								OmitEmpty: empty,
								Repeated:  repeated,
								Serial:    i,
							}
						}
						p.Write(&writer)
						g.Write(&writer)
					}
				}
			}
		}
		return true
	})

	if opt.Graphql {
		g.WriteFile()
	}
	if opt.Protobuf {
		p.WriteFile()
	}
}
