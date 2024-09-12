# Struct Builder

Simple code generator to create builder functions for structs in Go.

## Installation

```bash
go install github.com/TomWright/structbuilder/cmd/structbuilder
```

## Generating a Builder

```go
package example

import (
	"encoding/json"

	"github.com/foo/bar/abc"
)

//go:generate structbuilder -source=model.go -destination=model_builder.go -target=User

type User struct {
	ID        int
	Name      string
	Email     *string
	Something abc.Something
	Else      *abc.Else
	Numbers   []int

	iAmInternal string
}

func something(x json.Decoder) {
	_ = x.Decode(&User{})
}
```

## Usage

```go
package main

func main() {
	user := BuildUser(
		UserWithID(123),
		UserWithNumberAppend(1),
		UserWithNilEmail(),
	)
}
```