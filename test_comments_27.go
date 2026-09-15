package main

import (
	"fmt"
	goyaml "github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/parser"
)

func main() {
	yamlStr := `something: else      # comment after multiples spaces
emptyDir: {} # some comments
image: mariadb:11.4.3   # keep aligned
`
	file, _ := parser.ParseBytes([]byte(yamlStr), parser.ParseComments)

    path1, _ := goyaml.PathString("$.emptyDir")
    nodeToWrite1, _ := goyaml.ValueToNode(map[string]interface{}{})

    node1, _ := path1.FilterNode(file.Docs[0].Body)
    if node1.GetComment() != nil {
        nodeToWrite1.SetComment(node1.GetComment())
    }
    path1.ReplaceWithNode(file, nodeToWrite1)

    fmt.Print(file.String())
}
