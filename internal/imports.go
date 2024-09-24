package internal

import (
	"fmt"
	"go/ast"
	"io"
	"regexp"
	"strings"
)

var versionRegexp = regexp.MustCompile(`^v\d+$`)

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

func (i *Imports) Add(path string) string {
	im := Import{
		Path: strings.Trim(path, "\""),
	}
	i.List = append(i.List, im)
	return im.Name()
}

func (i *Imports) Render(w io.Writer) error {
	used := i.usedImports()
	if len(used) == 0 {
		return nil
	}

	if _, err := w.Write([]byte("\nimport (\n")); err != nil {
		return err
	}

	for _, imp := range used {
		if _, err := w.Write([]byte("\t" + imp.String() + "\n")); err != nil {
			return err
		}
	}

	if _, err := w.Write([]byte(")\n")); err != nil {
		return err
	}

	return nil
}

func (i *Imports) usedImports() []Import {
	var used []Import
	for _, imp := range i.List {
		if i.Uses[imp.Name()] == 0 {
			continue
		}
		used = append(used, imp)
	}
	return used
}

func (i *Imports) Count(name string) {
	imp, ok := i.find(name)
	if ok {
		i.Uses[imp.Name()]++
	}
}

func (i *Imports) find(name string) (Import, bool) {
	for _, imp := range i.List {
		if imp.Name() == name {
			return imp, true
		}
	}
	return Import{}, false
}

func (i *Imports) contains(imp Import) bool {
	for _, i2 := range i.List {
		if i2.Path == imp.Path {
			return true
		}
	}
	return false
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
	sLen := len(s)
	if sLen == 0 {
		return i.Path
	}
	if sLen > 1 && versionRegexp.MatchString(s[sLen-1]) {
		// If the last part of the path is a version number, use the second to last part
		return s[sLen-2]
	}
	return s[sLen-1]
}
