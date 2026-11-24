package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type violation struct {
	file string
	line int
	col  int
	text string
}

var skipDirs = map[string]bool{
	".git":         true,
	".agents":      true,
	"node_modules": true,
	"dist":         true,
	"vendor":       true,
	"reports":      true,
}

var (
	urlRe = regexp.MustCompile(`https?://\S+`)
	rfcRe = regexp.MustCompile(`\bRFC\s+\d+\b`)
)

func main() {
	root := flag.String("root", ".", "root directory to scan")
	flag.Parse()

	violations, err := lintTree(*root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "lintcomments: %v\n", err)
		os.Exit(2)
	}
	if len(violations) > 0 {
		for _, v := range violations {
			fmt.Fprintf(os.Stderr, "%s:%d:%d: disallowed comment: %s\n", v.file, v.line, v.col, v.text)
		}
		fmt.Fprintf(os.Stderr, "lintcomments: found %d disallowed comment(s)\n", len(violations))
		os.Exit(1)
	}
}

func lintTree(root string) ([]violation, error) {
	var all []violation
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if path != root && skipDirs[d.Name()] {
				return filepath.SkipDir
			}
			return nil
		}
		switch filepath.Ext(path) {
		case ".go":
			vs, err := lintGoFile(path)
			if err != nil {
				return err
			}
			all = append(all, vs...)
		case ".html":
			vs, err := lintHTMLFile(path)
			if err != nil {
				return err
			}
			all = append(all, vs...)
		}
		return nil
	})
	return all, err
}

func isInternalPath(path string) bool {
	cleaned := filepath.ToSlash(filepath.Clean(path))
	return cleaned == "internal" ||
		strings.HasPrefix(cleaned, "internal/") ||
		strings.Contains(cleaned, "/internal/")
}

func lintGoFile(path string) ([]violation, error) {
	src, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return lintGoSource(path, src)
}

func lintGoSource(path string, src []byte) ([]violation, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, path, src, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}

	isTest := strings.HasSuffix(path, "_test.go")
	isPublicPkg := !isTest && file.Name.Name != "main" && !isInternalPath(path)

	allowed := make(map[*ast.CommentGroup]bool)
	if isPublicPkg {
		collectPublicDocComments(file, allowed)
	}
	if isTest {
		collectExampleOutputs(file, allowed)
	}

	var violations []violation
	for _, cg := range file.Comments {
		if isLicenseHeader(cg, file.Package) {
			continue
		}
		if isAllowedExceptionGroup(cg) {
			continue
		}
		if allowed[cg] {
			for _, c := range cg.List {
				if isSuppressionComment(c.Text) {
					pos := fset.Position(c.Pos())
					violations = append(violations, violation{
						file: path,
						line: pos.Line,
						col:  pos.Column,
						text: truncateComment(c.Text),
					})
				}
			}
			continue
		}
		for _, c := range cg.List {
			if isGoDirective(c.Text) {
				continue
			}
			pos := fset.Position(c.Pos())
			violations = append(violations, violation{
				file: path,
				line: pos.Line,
				col:  pos.Column,
				text: truncateComment(c.Text),
			})
		}
	}
	return violations, nil
}

func collectPublicDocComments(file *ast.File, allowed map[*ast.CommentGroup]bool) {
	if file.Doc != nil {
		allowed[file.Doc] = true
	}
	for _, decl := range file.Decls {
		switch d := decl.(type) {
		case *ast.FuncDecl:
			if d.Doc != nil && ast.IsExported(d.Name.Name) && receiverExported(d.Recv) {
				allowed[d.Doc] = true
			}
		case *ast.GenDecl:
			hasExported := false
			for _, spec := range d.Specs {
				switch s := spec.(type) {
				case *ast.TypeSpec:
					if ast.IsExported(s.Name.Name) {
						hasExported = true
						if s.Doc != nil {
							allowed[s.Doc] = true
						}
						if s.Comment != nil {
							allowed[s.Comment] = true
						}
						collectExportedTypeMembers(s.Type, allowed)
					}
				case *ast.ValueSpec:
					specExported := false
					for _, name := range s.Names {
						if ast.IsExported(name.Name) {
							specExported = true
							hasExported = true
							break
						}
					}
					if specExported {
						if s.Doc != nil {
							allowed[s.Doc] = true
						}
						if s.Comment != nil {
							allowed[s.Comment] = true
						}
					}
				}
			}
			if hasExported && d.Doc != nil {
				allowed[d.Doc] = true
			}
		}
	}
}

func collectExportedTypeMembers(expr ast.Expr, allowed map[*ast.CommentGroup]bool) {
	switch t := expr.(type) {
	case *ast.StructType:
		if t.Fields == nil {
			return
		}
		for _, f := range t.Fields.List {
			exported := false
			for _, name := range f.Names {
				if ast.IsExported(name.Name) {
					exported = true
					break
				}
			}
			if exported {
				if f.Doc != nil {
					allowed[f.Doc] = true
				}
				if f.Comment != nil {
					allowed[f.Comment] = true
				}
			}
		}
	case *ast.InterfaceType:
		if t.Methods == nil {
			return
		}
		for _, m := range t.Methods.List {
			exported := len(m.Names) == 0
			for _, name := range m.Names {
				if ast.IsExported(name.Name) {
					exported = true
					break
				}
			}
			if exported {
				if m.Doc != nil {
					allowed[m.Doc] = true
				}
				if m.Comment != nil {
					allowed[m.Comment] = true
				}
			}
		}
	}
}

func receiverExported(recv *ast.FieldList) bool {
	if recv == nil || len(recv.List) == 0 {
		return true
	}
	t := recv.List[0].Type
	if star, ok := t.(*ast.StarExpr); ok {
		t = star.X
	}
	if idx, ok := t.(*ast.IndexExpr); ok {
		t = idx.X
	}
	if idxList, ok := t.(*ast.IndexListExpr); ok {
		t = idxList.X
	}
	ident, ok := t.(*ast.Ident)
	return ok && ast.IsExported(ident.Name)
}

func collectExampleOutputs(file *ast.File, allowed map[*ast.CommentGroup]bool) {
	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Recv != nil || fn.Body == nil || !strings.HasPrefix(fn.Name.Name, "Example") {
			continue
		}
		for _, cg := range file.Comments {
			if cg.Pos() < fn.Body.Pos() || cg.End() > fn.Body.End() || len(cg.List) == 0 {
				continue
			}
			first := strings.TrimSpace(cg.List[0].Text)
			if first == "// Output:" || first == "// Unordered output:" {
				allowed[cg] = true
			}
		}
	}
}

func isLicenseHeader(cg *ast.CommentGroup, pkgPos token.Pos) bool {
	if cg.End() >= pkgPos {
		return false
	}
	text := cg.Text()
	return strings.Contains(text, "SPDX-License-Identifier") ||
		strings.Contains(text, "Copyright ")
}

func isGoDirective(raw string) bool {
	return strings.HasPrefix(raw, "//go:") || strings.HasPrefix(raw, "//+build")
}

func isSuppressionComment(raw string) bool {
	trimmed := strings.TrimSpace(strings.TrimPrefix(strings.TrimPrefix(raw, "//"), "/*"))
	return strings.HasPrefix(trimmed, "nolint")
}

func isAllowedExceptionText(body string) bool {
	trimmed := strings.TrimSpace(body)
	if isSuppressionComment("//" + trimmed) {
		return false
	}
	lower := strings.ToLower(trimmed)
	if strings.HasPrefix(lower, "external:") || strings.HasPrefix(lower, "prettier-ignore") {
		return true
	}
	if strings.Contains(trimmed, "SPDX-License-Identifier") || strings.Contains(trimmed, "Copyright ") {
		return true
	}
	return urlRe.MatchString(trimmed) || rfcRe.MatchString(trimmed)
}

func isAllowedExceptionGroup(cg *ast.CommentGroup) bool {
	return isAllowedExceptionText(cg.Text())
}

func truncateComment(raw string) string {
	s := strings.TrimSpace(strings.ReplaceAll(raw, "\n", " "))
	if len(s) > 80 {
		return s[:77] + "..."
	}
	return s
}

func lintHTMLFile(path string) ([]violation, error) {
	src, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return lintHTMLSource(path, src), nil
}

func lintHTMLSource(path string, src []byte) []violation {
	var violations []violation
	lower := bytes.ToLower(src)
	i := 0
	for i < len(src) {
		if bytes.HasPrefix(src[i:], []byte("<!--")) {
			end := bytes.Index(src[i+4:], []byte("-->"))
			var raw []byte
			if end < 0 {
				raw = src[i:]
				i = len(src)
			} else {
				raw = src[i : i+4+end+3]
				endIdx := i + 4 + end + 3
				body := string(src[i+4 : i+4+end])
				if !isAllowedExceptionText(body) {
					line, col := lineCol(src, i)
					violations = append(violations, violation{file: path, line: line, col: col, text: truncateComment(string(raw))})
				}
				i = endIdx
				continue
			}
			line, col := lineCol(src, i)
			violations = append(violations, violation{file: path, line: line, col: col, text: truncateComment(string(raw))})
			continue
		}
		if bytes.HasPrefix(lower[i:], []byte("<style")) && isTagBoundary(lower, i+6) {
			openEnd := bytes.IndexByte(src[i:], '>')
			if openEnd < 0 {
				break
			}
			contentStart := i + openEnd + 1
			closeRel := bytes.Index(lower[contentStart:], []byte("</style"))
			if closeRel < 0 {
				break
			}
			violations = append(violations, scanBlockComments(path, src, contentStart, contentStart+closeRel, false)...)
			i = contentStart + closeRel + 7
			continue
		}
		if bytes.HasPrefix(lower[i:], []byte("<script")) && isTagBoundary(lower, i+7) {
			openEnd := bytes.IndexByte(src[i:], '>')
			if openEnd < 0 {
				break
			}
			contentStart := i + openEnd + 1
			closeRel := bytes.Index(lower[contentStart:], []byte("</script"))
			if closeRel < 0 {
				break
			}
			violations = append(violations, scanBlockComments(path, src, contentStart, contentStart+closeRel, true)...)
			i = contentStart + closeRel + 8
			continue
		}
		i++
	}
	return violations
}

func isTagBoundary(b []byte, idx int) bool {
	if idx >= len(b) {
		return true
	}
	c := b[idx]
	return c == '>' || c == ' ' || c == '\t' || c == '\n' || c == '\r' || c == '/'
}

func scanBlockComments(path string, full []byte, start, end int, allowLineComments bool) []violation {
	var violations []violation
	i := start
	var quote byte
	for i < end {
		c := full[i]
		if quote != 0 {
			if c == '\\' && i+1 < end {
				i += 2
				continue
			}
			if c == quote {
				quote = 0
			}
			i++
			continue
		}
		if c == '\'' || c == '"' || (allowLineComments && c == '`') {
			quote = c
			i++
			continue
		}
		if allowLineComments && c == '/' && i+1 < end && full[i+1] == '/' {
			lineEnd := bytes.IndexByte(full[i:end], '\n')
			var stop int
			if lineEnd < 0 {
				stop = end
			} else {
				stop = i + lineEnd
			}
			body := string(full[i+2 : stop])
			if !isAllowedExceptionText(body) {
				line, col := lineCol(full, i)
				violations = append(violations, violation{file: path, line: line, col: col, text: truncateComment(string(full[i:stop]))})
			}
			i = stop
			continue
		}
		if c == '/' && i+1 < end && full[i+1] == '*' {
			closeRel := bytes.Index(full[i+2:end], []byte("*/"))
			var stop int
			var body string
			if closeRel < 0 {
				stop = end
				body = string(full[i+2 : end])
			} else {
				stop = i + 2 + closeRel + 2
				body = string(full[i+2 : i+2+closeRel])
			}
			if !isAllowedExceptionText(body) {
				line, col := lineCol(full, i)
				violations = append(violations, violation{file: path, line: line, col: col, text: truncateComment(string(full[i:stop]))})
			}
			i = stop
			continue
		}
		i++
	}
	return violations
}

func lineCol(src []byte, offset int) (int, int) {
	line := 1 + bytes.Count(src[:offset], []byte{'\n'})
	lastNL := bytes.LastIndexByte(src[:offset], '\n')
	col := offset - lastNL
	return line, col
}
