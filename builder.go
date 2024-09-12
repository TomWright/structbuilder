package structbuilder

import (
	"fmt"
	"go/parser"
	"go/token"
	"io"

	"github.com/TomWright/structbuilder/internal"
)

func Build(structNames []string, destPackage string, r io.Reader, w io.Writer) error {
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

	structs := make([]internal.StructData, 0)
	imports := internal.NewImports()
	imports.AddFromFile(f)

	for _, structName := range structNames {
		targetTypeSpec, err := internal.FindTargetTypeSpec(f, structName)
		if err != nil {
			return fmt.Errorf("failed to find target type: %w", err)
		}

		d, err := internal.GetStructData(src, targetTypeSpec, imports)
		if err != nil {
			return fmt.Errorf("failed to get struct data: %w", err)
		}

		structs = append(structs, d)
	}

	if err := internal.WriteHeader(w, destPackage); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	if err := internal.WriteImports(w, imports); err != nil {
		return fmt.Errorf("failed to write imports: %w", err)
	}

	for _, s := range structs {
		if err := internal.WriteStructData(w, s); err != nil {
			return fmt.Errorf("failed to write struct data: %w", err)
		}
	}

	return nil
}
