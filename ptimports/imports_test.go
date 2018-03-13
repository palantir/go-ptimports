// Copyright 2016 Palantir Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ptimports_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/palantir/go-ptimports/ptimports"
)

func TestPtImports(t *testing.T) {
	for i, tc := range []struct {
		name    string
		in      string
		options *ptimports.Options
		want    string
	}{
		{
			"Does not simplify by default",
			`package foo

func Foo() {
	for i, _ := range []string{} {
		_ = i
	}
}
`,
			nil,
			`package foo

func Foo() {
	for i, _ := range []string{} {
		_ = i
	}
}
`,
		},
		{
			"Simplifies code when specified",
			`package foo

func Foo() {
	for i, _ := range []string{} {
		_ = i
	}
}
`,
			&ptimports.Options{
				Simplify: true,
			},
			`package foo

func Foo() {
	for i := range []string{} {
		_ = i
	}
}
`,
		},
		{
			"Imports not refactored if refactor is false",
			`package foo

import "github.com/palantir/go-ptimports/ptimports"
import "bytes"
import "golang.org/x/tools/imports"

func Foo() {
	_ = bytes.Buffer{}
	_ = ptimports.Process
	_ = imports.Process
}
`,

			nil,
			`package foo

import "github.com/palantir/go-ptimports/ptimports"
import "bytes"
import "golang.org/x/tools/imports"

func Foo() {
	_ = bytes.Buffer{}
	_ = ptimports.Process
	_ = imports.Process
}
`,
		},
		{
			"Refactors and groups imports based on builtin and external if refactor is true",
			`package foo

import "github.com/palantir/go-ptimports/ptimports"
import "bytes"
import "golang.org/x/tools/imports"

func Foo() {
	_ = bytes.Buffer{}
	_ = ptimports.Process
	_ = imports.Process
}
`,
			&ptimports.Options{
				Refactor: true,
			},
			`package foo

import (
	"bytes"

	"github.com/palantir/go-ptimports/ptimports"
	"golang.org/x/tools/imports"
)

func Foo() {
	_ = bytes.Buffer{}
	_ = ptimports.Process
	_ = imports.Process
}
`,
		},
		{
			"Refactors import added by goimports",
			`package foo

func Foo() {
	fmt.Println("foo")
}
`,
			&ptimports.Options{
				Refactor: true,
			},
			`package foo

import (
	"fmt"
)

func Foo() {
	fmt.Println("foo")
}
`,
		},
		{
			"Does not add import in format-only mode",
			`package foo

func Foo() {
	fmt.Println("foo")
}
`,
			&ptimports.Options{
				FormatOnly: true,
			},
			`package foo

func Foo() {
	fmt.Println("foo")
}
`,
		},
		{
			"Groups imports based on builtin, external, and project-local",
			`package foo

import "github.com/palantir/go-ptimports/ptimports"
import "bytes"
import "golang.org/x/tools/imports"

func Foo() {
	_ = bytes.Buffer{}
	_ = ptimports.Process
	_ = imports.Process
}
`,
			&ptimports.Options{
				Refactor: true,
				LocalPrefixes: []string{
					"github.com/palantir/go-ptimports/",
				},
			},
			`package foo

import (
	"bytes"

	"golang.org/x/tools/imports"

	"github.com/palantir/go-ptimports/ptimports"
)

func Foo() {
	_ = bytes.Buffer{}
	_ = ptimports.Process
	_ = imports.Process
}
`,
		},
		{
			"CGo import with multi-line comment",
			`package foo

// import "C"

import "unsafe"
import "io"

/*
#include <stdio.h>
#include <stdlib.h>

void myprint(char* s) {
	printf("%s\n", s);
}
*/
import "C"
import "archive/tar"


func Example() {
	/*
	multi-line comment
	 */
	cs := C.CString("Hello from stdio\n")
	C.myprint(cs)
	C.free(unsafe.Pointer(cs))

	// inline comment
	_ = io.Copy
	_ = tar.ErrFieldTooLong

}
`,
			&ptimports.Options{
				Refactor: true,
			},
			`package foo

// import "C"

/*
#include <stdio.h>
#include <stdlib.h>

void myprint(char* s) {
	printf("%s\n", s);
}
*/
import "C"

import (
	"archive/tar"
	"io"
	"unsafe"
)

func Example() {
	/*
		multi-line comment
	*/
	cs := C.CString("Hello from stdio\n")
	C.myprint(cs)
	C.free(unsafe.Pointer(cs))

	// inline comment
	_ = io.Copy
	_ = tar.ErrFieldTooLong

}
`,
		},
		{
			"CGo import with single-line and multi-line comments",
			`package foo

// #include <stdio.h>
// #include <stdlib.h>
import "C"
import "unsafe"

/*
#include <stdio.h>
#include <stdlib.h>

void myprint(char* s) {
	printf("%s\n", s);
}
*/
import "C"

func Print(s string) {
	cs := C.CString(s)
	C.fputs(cs, (*C.FILE)(C.stdout))
	C.free(unsafe.Pointer(cs))
}
`,
			&ptimports.Options{
				Refactor: true,
			},
			`package foo

// #include <stdio.h>
// #include <stdlib.h>
import "C"

/*
#include <stdio.h>
#include <stdlib.h>

void myprint(char* s) {
	printf("%s\n", s);
}
*/
import "C"

import (
	"unsafe"
)

func Print(s string) {
	cs := C.CString(s)
	C.fputs(cs, (*C.FILE)(C.stdout))
	C.free(unsafe.Pointer(cs))
}
`,
		},
	} {
		got, err := ptimports.Process("test.go", []byte(tc.in), tc.options)
		require.NoError(t, err, "Case %d: %s", i, tc.name)
		assert.Equal(t, tc.want, string(got), "Case %d: %s", i, tc.name)
	}
}
