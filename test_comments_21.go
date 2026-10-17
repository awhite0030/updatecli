package main

import (
	"fmt"
	goyaml "github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/parser"
    "github.com/goccy/go-yaml/ast"
)

func main() {
	yamlStr := `emptyDir: {}       # some comments
image: mariadb:11.4.3   # keep aligned
`
	file, _ := parser.ParseBytes([]byte(yamlStr), parser.ParseComments)
	path, _ := goyaml.PathString("$.emptyDir")

    node, _ := path.FilterNode(file.Docs[0].Body)

    // Instead of copying comment, what if we use ast.Node methods to replace value in place?
    fmt.Printf("node type: %T\n", node)

    if mn, ok := node.(*ast.MappingNode); ok && len(mn.Values) == 0 { // wait, no emptyDir is mapping
        // wait, FilterNode returns the MappingNode for `{}`? Yes.
    }
}
