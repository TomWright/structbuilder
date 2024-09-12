package internal

import (
	"fmt"
	"io"
)

func WriteHeader(w io.Writer, packageName string) error {
	if _, err := w.Write([]byte(GeneratedFileHeader + "\n")); err != nil {
		return err
	}
	if _, err := w.Write([]byte("package " + packageName + "\n\n")); err != nil {
		return err
	}
	return nil
}

func WriteImports(w io.Writer, i *Imports) error {
	x := i.String()
	if x == "" {
		return nil
	}
	if _, err := w.Write([]byte(x + "\n\n")); err != nil {
		return err
	}
	return nil
}

func WriteStructData(w io.Writer, d StructData) error {
	if _, err := w.Write([]byte(fmt.Sprintf("type %s %s\n\n", d.BuildOptionName(), d.BuildOptionFunc()))); err != nil {
		return err
	}

	if _, err := w.Write([]byte(fmt.Sprintf(`func %s(opts ...%s) *%s {
	res := new(%s)
	for _, opt := range opts {
		opt(res)
	}
	return res
}
`, d.BuildFuncName(), d.BuildOptionName(), d.structName, d.structName))); err != nil {
		return err
	}

	for _, f := range d.fields {
		if err := writeGenericFieldOptions(w, d, f); err != nil {
			return err
		}

		if f.typeSpec.pointer {
			if err := writePointerFieldOptions(w, d, f); err != nil {
				return err
			}
		}

		if f.typeSpec.slice {
			if err := writeSliceFieldOptions(w, d, f); err != nil {
				return err
			}
		}
	}

	return nil
}

func writeGenericFieldOptions(w io.Writer, d StructData, f FieldData) error {
	// TypeWithField(v Type) Option
	if _, err := w.Write([]byte(fmt.Sprintf(`
func %s(v %s) %s {
	return func(u *%s) {
		u.%s = v
	}
}
`, d.OptionNameForField(f), f.typeSpec.String(), d.BuildOptionName(), d.structName, f.name))); err != nil {
		return err
	}

	return nil
}

func writePointerFieldOptions(w io.Writer, d StructData, f FieldData) error {
	// PointerTypeWithNilField() Option
	if _, err := w.Write([]byte(fmt.Sprintf(`
func %s() %s {
	return func(u *%s) {
		u.%s = nil
	}
}
`, d.OptionNameForNilField(f), d.BuildOptionName(), d.structName, f.name))); err != nil {
		return err
	}

	// PointerTypeWithFieldValue(v Type) Option
	if _, err := w.Write([]byte(fmt.Sprintf(`
func %s(v %s) %s {
	return func(u *%s) {
		u.%s = &v
	}
}
`, d.OptionNameForFieldValue(f), f.typeSpec.Format(true, false), d.BuildOptionName(), d.structName, f.name))); err != nil {
		return err
	}

	return nil
}

func writeSliceFieldOptions(w io.Writer, d StructData, f FieldData) error {
	// SliceTypeWithNilField(v []Type) Option
	if _, err := w.Write([]byte(fmt.Sprintf(`
func %s() %s {
	return func(u *%s) {
		u.%s = nil
	}
}
`, d.OptionNameForNilField(f), d.BuildOptionName(), d.structName, f.name))); err != nil {
		return err
	}

	// SliceTypeWithEmptyField(v []Type) Option
	if _, err := w.Write([]byte(fmt.Sprintf(`
func %s() %s {
	return func(u *%s) {
		u.%s = make([]%s, 0)
	}
}
`, d.OptionNameForEmptyField(f), d.BuildOptionName(), d.structName, f.name, f.typeSpec.Format(true, true)))); err != nil {
		return err
	}

	// SliceTypeWithFieldAppend(v Type) Option
	if _, err := w.Write([]byte(fmt.Sprintf(`
func %s(v %s) %s {
	return func(u *%s) {
		u.%s = append(u.%s, v)
	}
}
`, d.OptionNameForFieldAppend(f), f.typeSpec.Format(true, true), d.BuildOptionName(), d.structName, f.name, f.name))); err != nil {
		return err
	}

	return nil
}
