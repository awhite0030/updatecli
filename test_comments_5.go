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

	nodeToWrite, _ := goyaml.ValueToNode(3)
	path, _ := goyaml.PathString("$..*") // match all
	path.ReplaceWithNode(file, nodeToWrite)

	fmt.Println(file.String())
}
