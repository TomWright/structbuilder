package internal

import "fmt"

type FieldData struct {
	name     string
	typeSpec TypeData
}

func (fd FieldData) Renderer(structName string, structType string, buildOptionName string) Renderer {
	returns := []FieldData{
		{
			typeSpec: TypeData{
				typeName: buildOptionName,
			},
		},
	}
	res := fd.genericRenderers(structName, structType, returns)
	if fd.typeSpec.pointer {
		res = append(res, fd.pointerRenderers(structName, structType, returns)...)
	}
	if fd.typeSpec.slice {
		res = append(res, fd.sliceRenderers(structName, structType, returns)...)
	}
	return NewMultipleRenderer(res...)
}

func (fd FieldData) genericRenderers(structName string, structType string, returns []FieldData) []Renderer {
	return []Renderer{
		funcRenderer{
			Name:    fmt.Sprintf("%sWith%s", structName, fd.name),
			Comment: fmt.Sprintf("sets %s to the given value.", fd.name),
			Args:    []FieldData{{name: "v", typeSpec: fd.typeSpec}},
			Returns: returns,
			Body: fmt.Sprintf(
				`
	return func(u *%s) {
		u.%s = v
	}`,
				structType, fd.name),
		},
	}
}

func (fd FieldData) sliceRenderers(structName string, structType string, returns []FieldData) []Renderer {
	return []Renderer{
		funcRenderer{
			Name:    fmt.Sprintf("%sWithNil%s", structName, fd.name),
			Comment: fmt.Sprintf("sets %s to nil.", fd.name),
			Returns: returns,
			Body: fmt.Sprintf(
				`
	return func(u *%s) {
		u.%s = nil
	}`,
				structType, fd.name),
		},
		funcRenderer{
			Name:    fmt.Sprintf("%sWithEmpty%s", structName, fd.name),
			Comment: fmt.Sprintf("sets %s to an empty slice.", fd.name),
			Returns: returns,
			Body: fmt.Sprintf(
				`
	return func(u *%s) {
		u.%s = make([]%s, 0)
	}`,
				structType, fd.name, fd.typeSpec.format(true, true)),
		},
		funcRenderer{
			Name:    fmt.Sprintf("%sWith%sAppend", structName, fd.name),
			Comment: fmt.Sprintf("appends the given value to %s.", fd.name),
			Args: []FieldData{{name: "v", typeSpec: TypeData{
				packageName: fd.typeSpec.packageName,
				typeName:    fd.typeSpec.typeName,
				pointer:     false,
				slice:       false,
			}}},
			Returns: returns,
			Body: fmt.Sprintf(
				`
	return func(u *%s) {
		u.%s = append(u.%s, v)
	}`,
				structType, fd.name, fd.name),
		},
	}
}

func (fd FieldData) pointerRenderers(structName string, structType string, returns []FieldData) []Renderer {
	return []Renderer{
		funcRenderer{
			Name:    fmt.Sprintf("%sWithNil%s", structName, fd.name),
			Comment: fmt.Sprintf("sets %s to nil.", fd.name),
			Returns: returns,
			Body: fmt.Sprintf(
				`
	return func(u *%s) {
		u.%s = nil
	}`,
				structType, fd.name),
		},
		funcRenderer{
			Name:    fmt.Sprintf("%sWith%sValue", structName, fd.name),
			Comment: fmt.Sprintf("sets %s to the given value.", fd.name),
			Args: []FieldData{{name: "v", typeSpec: TypeData{
				packageName: fd.typeSpec.packageName,
				typeName:    fd.typeSpec.typeName,
				pointer:     false,
				slice:       fd.typeSpec.slice,
			}}},
			Returns: returns,
			Body: fmt.Sprintf(
				`
	return func(u *%s) {
		u.%s = &v
	}`,
				structType, fd.name),
		},
	}
}
