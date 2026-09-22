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

    fmt.Printf("Node type: %T\n", node)
    if mn, ok := node.(*ast.MappingNode); ok {
        fmt.Printf("MappingNode\n")
        _ = mn
    } else if mn, ok := node.(*ast.MappingValueNode); ok {
        fmt.Printf("MappingValueNode\n")
        _ = mn
    }
}
