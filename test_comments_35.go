package main

import (
	"fmt"
	goyaml "github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/parser"
)

func main() {
	yamlStr := `a: 1
b: 2
`
	file, _ := parser.ParseBytes([]byte(yamlStr), parser.ParseComments)
	pathA, _ := goyaml.PathString("$.a")
    pathB, _ := goyaml.PathString("$.b")

    nodeToWrite, _ := goyaml.ValueToNode(3)

    pathA.ReplaceWithNode(file, nodeToWrite)
    pathB.ReplaceWithNode(file, nodeToWrite)

    fmt.Print(file.String())
}
