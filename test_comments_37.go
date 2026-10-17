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

    nodeToWrite, _ := goyaml.ValueToNode(3)

    // updateNode is called for pathA
    nodeA, _ := pathA.FilterNode(file.Docs[0].Body)
    if nodeToWrite.GetComment() == nil && nodeA.GetComment() != nil {
        nodeToWrite.SetComment(nodeA.GetComment())
    }
    pathA.ReplaceWithNode(file, nodeToWrite)

    // updateNode is called for pathB
    // Now nodeToWrite already HAS comment A!
    nodeB, _ := pathB.FilterNode(file.Docs[0].Body)
    // if nodeToWrite.GetComment() == nil will fail!
    if nodeB.GetComment() != nil {
        nodeToWrite.SetComment(nodeB.GetComment())
    }
    pathB.ReplaceWithNode(file, nodeToWrite)

    fmt.Print(file.String())
}
