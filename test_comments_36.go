package main

import (
	"fmt"
	goyaml "github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/parser"
    "github.com/goccy/go-yaml/ast"
)

func cloneNode(n ast.Node) ast.Node {
    // wait we don't need to clone, since updateNode doesn't reuse nodeToWrite across DIFFERENT files/keys.
    // However, does it reuse nodeToWrite across different keys in the same file?
    return nil
}
func main() {}
