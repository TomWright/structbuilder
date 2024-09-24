package internal

import "io"

type Renderer interface {
	Render(w io.Writer) error
}

func Render(w io.Writer, r ...Renderer) error {
	if len(r) == 0 {
		return nil
	}
	return NewMultipleRenderer(r...).Render(w)
}

type stringRenderer struct {
	Content string
}

func (s stringRenderer) Render(w io.Writer) error {
	_, err := w.Write([]byte(s.Content))
	return err
}

func NewStringRenderer(content string) Renderer {
	return stringRenderer{Content: content}
}

type multipleRenderer struct {
	Renderers []Renderer
}

func (m multipleRenderer) Render(w io.Writer) error {
	for _, r := range m.Renderers {
		if err := r.Render(w); err != nil {
			return err
		}
	}
	return nil
}

func NewMultipleRenderer(renderers ...Renderer) Renderer {
	return multipleRenderer{Renderers: renderers}
}
