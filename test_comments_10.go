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

    // To preserve comments on matched nodes when replacing:
    matchedNodes := []ast.Node{node} // in updatecli they have matchedNodes func

    nodeToWrite, _ := goyaml.ValueToNode(map[string]interface{}{"foo": "bar"})

    // we set comment on nodeToWrite?
    if len(matchedNodes) == 1 {
        if matchedNodes[0].GetComment() != nil {
            nodeToWrite.SetComment(matchedNodes[0].GetComment())
        }
    }

    path.ReplaceWithNode(file, nodeToWrite)
    fmt.Println(file.String())
}
