package structbuilder

import (
	"fmt"
	"go/parser"
	"go/token"
	"io"

	"github.com/TomWright/structbuilder/internal"
)

func Build(structNames []string, destPackage string, sourcePackage string, r io.Reader, w io.Writer) error {
	bytes, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("failed to read source: %w", err)
	}
	src := string(bytes)
	fileSet := token.NewFileSet()
	f, err := parser.ParseFile(fileSet, "", src, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("failed to parse source: %w", err)
	}

	if destPackage == "" {
		destPackage = f.Name.String()
	}

	imports := internal.NewImports()
	imports.AddFromFile(f)

	var sourcePackageAlias string
	if sourcePackage != "" {
		sourcePackageAlias = imports.Add(sourcePackage)
	}

	renderers := make([]internal.Renderer, 0)

	renderers = append(renderers, internal.NewHeaderRenderer(destPackage, internal.GeneratedFileHeader))
	renderers = append(renderers, imports)

	for _, structName := range structNames {
		targetTypeSpec, err := internal.FindTargetTypeSpec(f, structName)
		if err != nil {
			return fmt.Errorf("failed to find target type: %w", err)
		}

		d, err := internal.GetStructData(src, targetTypeSpec, imports, sourcePackageAlias)
		if err != nil {
			return fmt.Errorf("failed to get struct data: %w", err)
		}

		renderers = append(renderers, d)
	}

	if err := internal.Render(w, renderers...); err != nil {
		return fmt.Errorf("failed to render: %w", err)
	}

	return nil
}
