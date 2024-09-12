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

type BuildUserOption func(*User)

func BuildUser(opts ...BuildUserOption) *User {
	res := new(User)
	for _, opt := range opts {
		opt(res)
	}
	return res
}

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

type BuildUserOption func(*User)

func BuildUser(opts ...BuildUserOption) *User {
	res := new(User)
	for _, opt := range opts {
		opt(res)
	}
	return res
}

func UserWithID(v int) BuildUserOption {
	return func(u *User) {
		u.ID = v
	}
}

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
`,
			exp: internal.GeneratedFileHeader + `
package example

import (
	"github.com/TomWright/structbuilder/abc"
)

type BuildUserOption func(*User)

func BuildUser(opts ...BuildUserOption) *User {
	res := new(User)
	for _, opt := range opts {
		opt(res)
	}
	return res
}

func UserWithID(v int) BuildUserOption {
	return func(u *User) {
		u.ID = v
	}
}

func UserWithName(v string) BuildUserOption {
	return func(u *User) {
		u.Name = v
	}
}

func UserWithEmail(v *string) BuildUserOption {
	return func(u *User) {
		u.Email = v
	}
}

func UserWithNilEmail() BuildUserOption {
	return func(u *User) {
		u.Email = nil
	}
}

func UserWithEmailValue(v string) BuildUserOption {
	return func(u *User) {
		u.Email = &v
	}
}

func UserWithSomething(v abc.Something) BuildUserOption {
	return func(u *User) {
		u.Something = v
	}
}

func UserWithElse(v *abc.Else) BuildUserOption {
	return func(u *User) {
		u.Else = v
	}
}

func UserWithNilElse() BuildUserOption {
	return func(u *User) {
		u.Else = nil
	}
}

func UserWithElseValue(v abc.Else) BuildUserOption {
	return func(u *User) {
		u.Else = &v
	}
}

func UserWithNumbers(v []int) BuildUserOption {
	return func(u *User) {
		u.Numbers = v
	}
}

func UserWithNilNumbers() BuildUserOption {
	return func(u *User) {
		u.Numbers = nil
	}
}

func UserWithEmptyNumbers() BuildUserOption {
	return func(u *User) {
		u.Numbers = make([]int, 0)
	}
}

func UserWithNumbersAppend(v int) BuildUserOption {
	return func(u *User) {
		u.Numbers = append(u.Numbers, v)
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

			if err := structbuilder.Build(tc.targetStruct, tc.destPackage, in, out); err != nil {
				t.Fatal(err)
			}

			if exp, got := tc.exp, out.String(); exp != got {
				t.Errorf("expected:\n%s\ngot:\n%s", exp, got)
			}
		})
	}

}
