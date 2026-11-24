package main

import (
	"testing"
)

func TestLintGoSource(t *testing.T) {
	allowedPublic := `// SPDX-License-Identifier: MIT
// Package sample documents public API.
package sample

//go:embed data.txt
var embedded string

// ExportedType is a public API contract.
type ExportedType struct {
	// ExportedField is documented.
	ExportedField int
	unexported    int
}

// ExportedFunc is documented.
func (e *ExportedType) ExportedFunc() int {
	// external: vendor protocol requires non-zero padding
	// https://github.com/example/repo/issues/42
	return e.ExportedField
}
`
	if vs, err := lintGoSource("pkg/sample/sample.go", []byte(allowedPublic)); err != nil || len(vs) != 0 {
		t.Fatalf("expected 0 violations on valid public file, got err=%v violations=%+v", err, vs)
	}

	disallowedPublic := `package sample

// unexported type doc should fail
type hidden struct {
	x int
}

// ExportedFunc has a bad inline comment.
func ExportedFunc() {
	// narration inside function body
	_ = 1
}

//nolint:errcheck
func Suppressed() {}
`
	vs, err := lintGoSource("pkg/sample/bad.go", []byte(disallowedPublic))
	if err != nil {
		t.Fatal(err)
	}
	if len(vs) != 3 {
		t.Fatalf("expected 3 violations (unexported doc, inline narration, nolint), got %d: %+v", len(vs), vs)
	}

	internalFile := `package secret

// Doc comment inside internal package should fail.
func ExportedInInternal() {}
`
	vs, err = lintGoSource("internal/secret/secret.go", []byte(internalFile))
	if err != nil {
		t.Fatal(err)
	}
	if len(vs) != 1 {
		t.Fatalf("expected 1 violation inside internal/, got %d: %+v", len(vs), vs)
	}

	exampleTest := `package sample_test

import "fmt"

func ExampleThing() {
	fmt.Println("hi")
	// Output:
	// hi
}
`
	if vs, err := lintGoSource("pkg/sample/example_test.go", []byte(exampleTest)); err != nil || len(vs) != 0 {
		t.Fatalf("expected Example Output comment to be allowed, got err=%v vs=%+v", err, vs)
	}
}

func TestLintHTMLSource(t *testing.T) {
	clean := `<!doctype html>
<html>
<style>.a { color: red; }</style>
<script>const url = "https://example.com"; const x = 1;</script>
</html>`
	if vs := lintHTMLSource("template.html", []byte(clean)); len(vs) != 0 {
		t.Fatalf("expected 0 HTML violations, got %+v", vs)
	}

	dirty := `<!doctype html>
<!-- banner -->
<style>/* css banner */ .a { color: red; }</style>
<script>
// js banner
const x = 1;
</script>`
	if vs := lintHTMLSource("template.html", []byte(dirty)); len(vs) != 3 {
		t.Fatalf("expected 3 HTML/CSS/JS violations, got %d: %+v", len(vs), vs)
	}
}
