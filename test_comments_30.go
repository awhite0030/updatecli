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
    nodeToWrite, _ := goyaml.ValueToNode(map[string]interface{}{})

    node, _ := path.FilterNode(file.Docs[0].Body)

    // In updateNode we have `matched []ast.Node` already!
    var matched []ast.Node
    matched = append(matched, node)

    if nodeToWrite.GetComment() == nil && len(matched) > 0 {
        if matched[0].GetComment() != nil {
            nodeToWrite.SetComment(matched[0].GetComment())
        }
    }

    path.ReplaceWithNode(file, nodeToWrite)

    fmt.Print(file.String())
}
