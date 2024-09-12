package internal

import (
	"fmt"
	"go/ast"
	"strings"
)

const GeneratedFileHeader = `// Code generated by structbuilder. DO NOT EDIT.`

type Imports struct {
	List []Import
	Uses map[string]int
}

func NewImports() *Imports {
	return &Imports{
		List: []Import{},
		Uses: map[string]int{},
	}
}

func (i *Imports) String() string {
	used := 0
	res := "import (\n"
	for _, imp := range i.List {
		if i.Uses[imp.Name()] == 0 {
			continue
		}
		used++
		res += "\t" + imp.String() + "\n"
	}
	res += ")"
	if used == 0 {
		return ""
	}
	return res
}

func (i *Imports) Count(name string) {
	imp, ok := i.Find(name)
	if ok {
		i.Uses[imp.Name()]++
	}
}

func (i *Imports) Find(name string) (Import, bool) {
	for _, imp := range i.List {
		if imp.Name() == name {
			return imp, true
		}
	}
	return Import{}, false
}

func (i *Imports) Contains(imp Import) bool {
	for _, i2 := range i.List {
		if i2.Path == imp.Path {
			return true
		}
	}
	return false
}

func (i *Imports) AddFromFile(f *ast.File) {
	for _, importSpec := range f.Imports {
		im := Import{
			Path: strings.Trim(importSpec.Path.Value, "\""),
		}
		if importSpec.Name != nil {
			im.Alias = importSpec.Name.String()
		}
		i.List = append(i.List, im)
	}
}

type Import struct {
	Path  string
	Alias string
}

func (i Import) String() string {
	if i.Alias != "" {
		return fmt.Sprintf("%s %q", i.Alias, i.Path)
	}
	return fmt.Sprintf("%q", i.Path)
}

func (i Import) Name() string {
	if i.Alias != "" {
		return i.Alias
	}
	s := strings.Split(i.Path, "/")
	return s[len(s)-1]
}

type StructData struct {
	structName string
	fields     []FieldData
}

type FieldData struct {
	name     string
	typeSpec TypeData
}

type TypeData struct {
	packageName string
	typeName    string
	pointer     bool
	slice       bool
}

func (td TypeData) String() string {
	return td.Format(false, false)
}

func (td TypeData) Format(ignorePointer bool, ignoreSlice bool) string {
	res := td.typeName
	if td.packageName != "" {
		res = td.packageName + "." + res
	}
	if !ignoreSlice && td.slice {
		res = "[]" + res
	}
	if !ignorePointer && td.pointer {
		res = "*" + res
	}
	return res
}

func (sd StructData) BuildOptionName() string {
	return "Build" + sd.structName + "Option"
}

func (sd StructData) BuildOptionFunc() string {
	return "func(*" + sd.structName + ")"
}

func (sd StructData) BuildFuncName() string {
	return "Build" + sd.structName
}

func (sd StructData) OptionNameForField(f FieldData) string {
	return fmt.Sprintf("%sWith%s", sd.structName, f.name)
}

func (sd StructData) OptionNameForNilField(f FieldData) string {
	return fmt.Sprintf("%sWithNil%s", sd.structName, f.name)
}

func (sd StructData) OptionNameForEmptyField(f FieldData) string {
	return fmt.Sprintf("%sWithEmpty%s", sd.structName, f.name)
}

func (sd StructData) OptionNameForFieldValue(f FieldData) string {
	return fmt.Sprintf("%sWith%sValue", sd.structName, f.name)
}

func (sd StructData) OptionNameForFieldAppend(f FieldData) string {
	return fmt.Sprintf("%sWith%sAppend", sd.structName, f.name)
}