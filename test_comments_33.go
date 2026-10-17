package main

import (
	"fmt"
	goyaml "github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/parser"
    "github.com/goccy/go-yaml/ast"
)

func main() {
	yamlStr := `volumes:
- emptyDir: {}       # some comments
  name: a
- emptyDir: {}       # some other comments
  name: b
`
	file, _ := parser.ParseBytes([]byte(yamlStr), parser.ParseComments)
	path, _ := goyaml.PathString("$.volumes[*].emptyDir")

    node, _ := path.FilterNode(file.Docs[0].Body)

    var matched []ast.Node
    if mn, ok := node.(*ast.SequenceNode); ok {
        for _, val := range mn.Values {
            matched = append(matched, val)
        }
    }

    nodeToWrite, _ := goyaml.ValueToNode(map[string]interface{}{"foo": "bar"})

    // Copy comment of the first match if it exists
    if nodeToWrite.GetComment() == nil && len(matched) > 0 {
        if matched[0].GetComment() != nil {
            nodeToWrite.SetComment(matched[0].GetComment())
        }
    }

    path.ReplaceWithNode(file, nodeToWrite)

    fmt.Print(file.String())
}
