package replace

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/printer"
	"go/token"
	"strings"
)

func removeType(f *ast.File, t string) error {
	var isFound bool
	for i, d := range f.Decls {
		g, ok := d.(*ast.GenDecl)
		if !ok {
			continue
		}
		if len(g.Specs) != 1 {
			continue
		}
		ts, ok := g.Specs[0].(*ast.TypeSpec)
		if !ok {
			continue
		}
		if ts.Name.Name == t {
			isFound = true
			g.Doc.List = []*ast.Comment{}
			f.Decls = append(f.Decls[:i], f.Decls[i+1:]...)
		}
	}
	if !isFound {
		return ErrType
	}
	return nil
}

func replaceTypeInNode(node ast.Node, old, new string) {
	ast.Inspect(node, func(n ast.Node) bool {
		switch v := n.(type) {
		case *ast.Ellipsis:
			replaceTypeInEllipsis(v, old, new)
		case *ast.CompositeLit:
			replaceTypeInCompositeLit(v, old, new)
		case *ast.Comment:
			replaceTypeInComment(v, old, new)
		case *ast.CallExpr:
			replaceTypeInCallExpr(v, old, new)
		case *ast.ChanType:
			replaceTypeInChanType(v, old, new)
		case *ast.MapType:
			replaceTypeInMapType(v, old, new)
		case *ast.ArrayType:
			replaceTypeInArrayType(v, old, new)
		case *ast.ValueSpec:
			replaceTypeInValueSpec(v, old, new)
		case *ast.Field:
			replaceTypeInField(v, old, new)
		}
		return true
	})
}

func replaceTypeInEllipsis(n *ast.Ellipsis, old, new string) {
	t, ok := n.Elt.(*ast.Ident)
	if !ok || t.Name != old {
		return
	}
	n.Elt = ast.NewIdent(new)
}

func replaceTypeInCompositeLit(n *ast.CompositeLit, old, new string) {
	t, ok := n.Type.(*ast.Ident)
	if !ok || t.Name != old {
		return
	}
	n.Type = ast.NewIdent(new)
}

func replaceTypeInComment(n *ast.Comment, old, new string) {
	n.Text = strings.Replace(n.Text, old, new, -1)
}

func replaceTypeInCallExpr(n *ast.CallExpr, old, new string) {
	t, ok := n.Fun.(*ast.Ident)
	if !ok || t.Name != old {
		return
	}
	n.Fun = ast.NewIdent(new)
}

func replaceTypeInChanType(n *ast.ChanType, old, new string) {
	t, ok := n.Value.(*ast.Ident)
	if !ok || t.Name != old {
		return
	}
	n.Value = ast.NewIdent(new)
}

func replaceTypeInMapType(n *ast.MapType, old, new string) {
	replaceTypeInMapKey(n, old, new)
	replaceTypeInMapValue(n, old, new)
}

func replaceTypeInMapKey(n *ast.MapType, old, new string) {
	t, ok := n.Key.(*ast.Ident)
	if !ok || t.Name != old {
		return
	}
	n.Key = ast.NewIdent(new)
}

func replaceTypeInMapValue(n *ast.MapType, old, new string) {
	t, ok := n.Value.(*ast.Ident)
	if !ok || t.Name != old {
		return
	}
	n.Value = ast.NewIdent(new)
}

func replaceTypeInArrayType(n *ast.ArrayType, old, new string) {
	t, ok := n.Elt.(*ast.Ident)
	if !ok || t.Name != old {
		return
	}
	n.Elt = ast.NewIdent(new)
}

func replaceTypeInValueSpec(n *ast.ValueSpec, old, new string) {
	t, ok := n.Type.(*ast.Ident)
	if !ok || t.Name != old {
		return
	}
	n.Type = ast.NewIdent(new)
}

func replaceTypeInField(n *ast.Field, old, new string) {
	t, ok := n.Type.(*ast.Ident)
	if !ok || t.Name != old {
		return
	}
	n.Type = ast.NewIdent(new)
}

func fileToBytes(fset *token.FileSet, f *ast.File) ([]byte, error) {
	var buf bytes.Buffer
	if err := printer.Fprint(&buf, fset, f); err != nil {
		return nil, err
	}
	return format.Source(buf.Bytes())
}
