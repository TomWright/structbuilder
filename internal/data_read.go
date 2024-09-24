package internal

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"
)

func GetStructData(src string, target *ast.TypeSpec, imports *Imports, sourceAlias string) (StructData, error) {
	res := StructData{
		structName:    target.Name.String(),
		structPackage: sourceAlias,
		fields:        make([]FieldData, 0),
	}

	sType, ok := target.Type.(*ast.StructType)
	if !ok {
		return StructData{}, fmt.Errorf("target type %q is not a struct", target.Name)
	}

	for _, f := range sType.Fields.List {
		if !f.Names[0].IsExported() {
			continue
		}
		a, b, c, d := exprTypeToString(src, f.Type)
		imports.Count(a)
		fd := FieldData{
			name: f.Names[0].Name,
			typeSpec: TypeData{
				packageName: a,
				typeName:    b,
				pointer:     c,
				slice:       d,
			},
		}
		res.fields = append(res.fields, fd)
	}

	return res, nil
}

func FindTargetTypeSpec(f *ast.File, target string) (*ast.TypeSpec, error) {
	for _, d := range f.Decls {
		genDecl, ok := d.(*ast.GenDecl)
		if !ok {
			continue
		}
		if genDecl.Tok != token.TYPE {
			continue
		}
		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}
			if typeSpec.Name.Name == target {
				return typeSpec, nil
			}
		}
	}
	return nil, fmt.Errorf("target type %q not found", target)
}

func exprTypeToString(src string, expr ast.Expr) (string, string, bool, bool) {
	n := src[expr.Pos()-1 : expr.End()-1]
	s := strings.Split(n, ".")
	p := false
	sl := false
	if strings.HasPrefix(s[0], "*") {
		p = true
		s[0] = s[0][1:]
	}
	if strings.HasPrefix(s[0], "[]") {
		sl = true
		s[0] = s[0][2:]
	}
	if len(s) == 1 {
		return "", s[0], p, sl
	}
	return s[0], s[1], p, sl
}
