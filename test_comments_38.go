package main

import (
	"fmt"
	goyaml "github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/parser"
)

func main() {
	yamlStr := `a: 1 # comment A
b: 2 # comment B
`
	file, _ := parser.ParseBytes([]byte(yamlStr), parser.ParseComments)
	pathA, _ := goyaml.PathString("$.a")
    pathB, _ := goyaml.PathString("$.b")

    nodeToWriteA, _ := goyaml.ValueToNode(3)
    nodeA, _ := pathA.FilterNode(file.Docs[0].Body)
    if nodeToWriteA.GetComment() == nil && nodeA.GetComment() != nil {
        nodeToWriteA.SetComment(nodeA.GetComment())
    }
    pathA.ReplaceWithNode(file, nodeToWriteA)

    nodeToWriteB, _ := goyaml.ValueToNode(3)
    nodeB, _ := pathB.FilterNode(file.Docs[0].Body)
    if nodeToWriteB.GetComment() == nil && nodeB.GetComment() != nil {
        nodeToWriteB.SetComment(nodeB.GetComment())
    }
    pathB.ReplaceWithNode(file, nodeToWriteB)

    fmt.Print(file.String())
}
