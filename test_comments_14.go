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

    // We want to see how to preserve comments
    // `matched` contains the actual matched nodes
    var matched []ast.Node

    if mn, ok := node.(*ast.MappingNode); ok {
        for _, val := range mn.Values {
            matched = append(matched, val.Value)
        }
    }

    nodeToWrite, _ := goyaml.ValueToNode(map[string]interface{}{})

    if len(matched) == 1 && matched[0].GetComment() != nil {
        nodeToWrite.SetComment(matched[0].GetComment())
    }

    path.ReplaceWithNode(file, nodeToWrite)
    fmt.Println(file.String())
}
