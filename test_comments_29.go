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

    // Suppose we check if nodeToWrite does not have a comment, and if so, copy it from matched nodes.
    // In updatecli they only have ONE nodeToWrite but MULTIPLE matched nodes.
    // If they have multiple matched nodes with different comments... well they replace them all with ONE nodeToWrite anyway!
    // So the comment of the first matched node can be copied. That's a reasonable fallback.

    node, _ := path.FilterNode(file.Docs[0].Body)
    var matched []ast.Node
    if mn, ok := node.(*ast.MappingNode); ok {
        for _, val := range mn.Values {
            matched = append(matched, val)
        }
    } else {
        matched = append(matched, node)
    }

    if nodeToWrite.GetComment() == nil && len(matched) > 0 {
        if matched[0].GetComment() != nil {
            nodeToWrite.SetComment(matched[0].GetComment())
        }
    }

    path.ReplaceWithNode(file, nodeToWrite)

    fmt.Print(file.String())
}
