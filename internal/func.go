package internal

import (
	"io"
)

type funcRenderer struct {
	Name    string
	Comment string
	Args    []FieldData
	Returns []FieldData
	Body    string
}

func (f funcRenderer) Render(w io.Writer) error {
	if _, err := w.Write([]byte("\n")); err != nil {
		return err
	}
	if f.Comment != "" {
		if _, err := w.Write([]byte("// " + f.Name + " " + f.Comment + "\n")); err != nil {
			return err
		}
	}
	if _, err := w.Write([]byte("func " + f.Name + "(")); err != nil {
		return err
	}
	for i, a := range f.Args {
		if i > 0 {
			if _, err := w.Write([]byte(", ")); err != nil {
				return err
			}
		}
		if _, err := w.Write([]byte(a.name + " " + a.typeSpec.String())); err != nil {
			return err
		}
	}
	if _, err := w.Write([]byte(")")); err != nil {
		return err
	}
	if len(f.Returns) > 0 {
		if _, err := w.Write([]byte(" ")); err != nil {
			return err
		}
		if len(f.Returns) > 1 {
			if _, err := w.Write([]byte("(")); err != nil {
				return err
			}
		}
		for i, r := range f.Returns {
			if i > 0 {
				if _, err := w.Write([]byte(", ")); err != nil {
					return err
				}
			}
			if r.name != "" {
				if _, err := w.Write([]byte(r.name + " ")); err != nil {
					return err
				}
			}
			if _, err := w.Write([]byte(r.typeSpec.String())); err != nil {
				return err
			}
		}
		if len(f.Returns) > 1 {
			if _, err := w.Write([]byte(")")); err != nil {
				return err
			}
		}
	}
	if _, err := w.Write([]byte(" {")); err != nil {
		return err
	}
	if _, err := w.Write([]byte(f.Body)); err != nil {
		return err
	}
	if _, err := w.Write([]byte("\n}\n")); err != nil {
		return err
	}
	return nil
}
