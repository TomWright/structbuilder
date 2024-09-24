package internal

import (
	"fmt"
	"io"
)

type StructData struct {
	structName    string
	structPackage string
	fields        []FieldData
}

func (sd StructData) Render(w io.Writer) error {
	buildOptionName := fmt.Sprintf("Build%sOption", sd.structName)
	structType := sd.structName
	if sd.structPackage != "" {
		structType = sd.structPackage + "." + structType
	}
	res := []Renderer{
		builderTypeRenderer(sd.structName, structType, buildOptionName),
	}
	for _, f := range sd.fields {
		res = append(res, f.Renderer(sd.structName, structType, buildOptionName))
	}
	return NewMultipleRenderer(res...).Render(w)
}

type TypeData struct {
	packageName string
	typeName    string
	pointer     bool
	slice       bool
}

func (td TypeData) String() string {
	return td.format(false, false)
}

func (td TypeData) format(ignorePointer bool, ignoreSlice bool) string {
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

func builderTypeRenderer(structName string, structType string, buildOptionName string) Renderer {
	return NewMultipleRenderer(
		NewStringRenderer(fmt.Sprintf(`
// %s is a function that sets the given options on a %s.
type %s func(*%s)
`, buildOptionName, structName, buildOptionName, structType)),
		NewStringRenderer(fmt.Sprintf(`
// Build%s creates a new %s with the given options.
func Build%s(opts ...%s) *%s {
	res := new(%s)
	for _, opt := range opts {
		opt(res)
	}
	return res
}
`, structName, structName, structName, buildOptionName, structType, structType)),
	)
}
