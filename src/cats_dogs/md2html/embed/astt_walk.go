package embed

import (
	"github.com/yuin/goldmark/ast"
)

type ASTTWalker func(n ast.Node) (ast.WalkStatus, error)

func ASTTWalk(n ast.Node, walker ASTTWalker) error {
	_, err := walkHelper(n, walker)
	return err
}

func walkHelper(n ast.Node, walker ASTTWalker) (ast.WalkStatus, error) {
	status, err := walker(n)
	if err != nil || status == ast.WalkStop {
		return status, err
	}
	if status != ast.WalkSkipChildren {
		for c := n.FirstChild(); c != nil; {
			next := c.NextSibling()
			if st, err := walkHelper(c, walker); err != nil || st == ast.WalkStop {
				return ast.WalkStop, err
			}
			c = next
		}
	}
	return ast.WalkContinue, nil
}
