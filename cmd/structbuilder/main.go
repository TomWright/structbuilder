package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/TomWright/structbuilder"
)

var (
	destination   = flag.String("destination", "", "The destination file")
	packageName   = flag.String("package", "", "The destination package name")
	target        = flag.String("target", "", "The target struct in the source file")
	source        = flag.String("source", "", "The source file")
	sourcePackage = flag.String("source-package", "", "The source package")
)

func main() {
	flag.Parse()

	if *destination == "" {
		fmt.Println("destination flag is required")
		os.Exit(1)
	}
	if *target == "" {
		fmt.Println("target flag is required")
		os.Exit(1)
	}
	if *source == "" {
		fmt.Println("source flag is required")
		os.Exit(1)
	}

	f, err := os.Open(*source)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer f.Close()

	outF, err := os.Create(*destination)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer outF.Close()

	if err := structbuilder.Build(strings.Split(*target, ","), *packageName, *sourcePackage, f, outF); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	os.Exit(0)
}
