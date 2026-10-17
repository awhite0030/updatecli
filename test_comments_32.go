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
	path, _ := goyaml.PathString("$.*")
    nodeToWrite, _ := goyaml.ValueToNode("NEW")

    path.ReplaceWithNode(file, nodeToWrite)

    fmt.Print(file.String())
}
