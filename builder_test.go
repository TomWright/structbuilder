package structbuilder_test

import (
	"bytes"
	"testing"

	"github.com/TomWright/structbuilder"
	"github.com/TomWright/structbuilder/internal"
)

func TestBuild(t *testing.T) {

	testCases := []struct {
		name         string
		in           string
		targetStruct []string
		destPackage  string
		exp          string
	}{
		{
			name:         "Unexported fields are ignored",
			targetStruct: []string{"User"},
			in: `package models

type User struct {
	ID        		 int
	somethingPrivate string
}
`,
			exp: internal.GeneratedFileHeader + `
package models

// BuildUserOption is a function that sets the given options on a User.
type BuildUserOption func(*User)

// BuildUser creates a new User with the given options.
func BuildUser(opts ...BuildUserOption) *User {
	res := new(User)
	for _, opt := range opts {
		opt(res)
	}
	return res
}

// UserWithID sets ID to the given value.
func UserWithID(v int) BuildUserOption {
	return func(u *User) {
		u.ID = v
	}
}
`,
		},
		{
			name:         "Unused imports are ignored",
			targetStruct: []string{"User"},
			in: `package models

import (
	"encoding/json"
	"some/other/package/xyz"
)

type User struct {
	ID        		 int
	Other            xyz.Thing
	somethingPrivate string
}

func something(x json.Decoder) {
	_ = x.Decode(&User{})
}
`,
			exp: internal.GeneratedFileHeader + `
package models

import (
	"some/other/package/xyz"
)

// BuildUserOption is a function that sets the given options on a User.
type BuildUserOption func(*User)

// BuildUser creates a new User with the given options.
func BuildUser(opts ...BuildUserOption) *User {
	res := new(User)
	for _, opt := range opts {
		opt(res)
	}
	return res
}

// UserWithID sets ID to the given value.
func UserWithID(v int) BuildUserOption {
	return func(u *User) {
		u.ID = v
	}
}

// UserWithOther sets Other to the given value.
func UserWithOther(v xyz.Thing) BuildUserOption {
	return func(u *User) {
		u.Other = v
	}
}
`,
		},
		{
			name:         "Everything",
			targetStruct: []string{"User"},
			in: `package example

import (
	"encoding/json"

	"github.com/TomWright/structbuilder/abc"
	"github.com/TomWright/structbuilder/foo/v2"
)

//go:generate structbuilder -source=model.go -destination=model_builder.go -target=User

type User struct {
	ID        int
	Name      string
	Email     *string
	Something abc.Something
	Else      *abc.Else
	Numbers   []int
	Foo       foo.Foo

	iAmInternal string
}

func something(x json.Decoder) {
	_ = x.Decode(&User{})
}
`,
			exp: internal.GeneratedFileHeader + `
package example

import (
	"github.com/TomWright/structbuilder/abc"
	"github.com/TomWright/structbuilder/foo/v2"
)

// BuildUserOption is a function that sets the given options on a User.
type BuildUserOption func(*User)

// BuildUser creates a new User with the given options.
func BuildUser(opts ...BuildUserOption) *User {
	res := new(User)
	for _, opt := range opts {
		opt(res)
	}
	return res
}

// UserWithID sets ID to the given value.
func UserWithID(v int) BuildUserOption {
	return func(u *User) {
		u.ID = v
	}
}

// UserWithName sets Name to the given value.
func UserWithName(v string) BuildUserOption {
	return func(u *User) {
		u.Name = v
	}
}

// UserWithEmail sets Email to the given value.
func UserWithEmail(v *string) BuildUserOption {
	return func(u *User) {
		u.Email = v
	}
}

// UserWithNilEmail sets Email to nil.
func UserWithNilEmail() BuildUserOption {
	return func(u *User) {
		u.Email = nil
	}
}

// UserWithEmailValue sets Email to the given value.
func UserWithEmailValue(v string) BuildUserOption {
	return func(u *User) {
		u.Email = &v
	}
}

// UserWithSomething sets Something to the given value.
func UserWithSomething(v abc.Something) BuildUserOption {
	return func(u *User) {
		u.Something = v
	}
}

// UserWithElse sets Else to the given value.
func UserWithElse(v *abc.Else) BuildUserOption {
	return func(u *User) {
		u.Else = v
	}
}

// UserWithNilElse sets Else to nil.
func UserWithNilElse() BuildUserOption {
	return func(u *User) {
		u.Else = nil
	}
}

// UserWithElseValue sets Else to the given value.
func UserWithElseValue(v abc.Else) BuildUserOption {
	return func(u *User) {
		u.Else = &v
	}
}

// UserWithNumbers sets Numbers to the given value.
func UserWithNumbers(v []int) BuildUserOption {
	return func(u *User) {
		u.Numbers = v
	}
}

// UserWithNilNumbers sets Numbers to nil.
func UserWithNilNumbers() BuildUserOption {
	return func(u *User) {
		u.Numbers = nil
	}
}

// UserWithEmptyNumbers sets Numbers to an empty slice.
func UserWithEmptyNumbers() BuildUserOption {
	return func(u *User) {
		u.Numbers = make([]int, 0)
	}
}

// UserWithNumbersAppend appends the given value to Numbers.
func UserWithNumbersAppend(v int) BuildUserOption {
	return func(u *User) {
		u.Numbers = append(u.Numbers, v)
	}
}

// UserWithFoo sets Foo to the given value.
func UserWithFoo(v foo.Foo) BuildUserOption {
	return func(u *User) {
		u.Foo = v
	}
}
`,
		},
	}

	for _, testCase := range testCases {
		tc := testCase
		t.Run(testCase.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			in := &bytes.Buffer{}
			in.WriteString(tc.in)

			if err := structbuilder.Build(tc.targetStruct, tc.destPackage, "", in, out); err != nil {
				t.Fatal(err)
			}

			if exp, got := tc.exp, out.String(); exp != got {
				t.Errorf("expected:\n%s\ngot:\n%s", exp, got)
			}
		})
	}

}
