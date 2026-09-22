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

    path, _ := goyaml.PathString("$.*")

    nodeToWrite, _ := goyaml.ValueToNode("NEW")

    // In updatecli they get matched using matchedNodes(node, key)
    // Here we can just see if we can copy comments from matched to nodeToWrite? No, because nodeToWrite is ONE node, but ReplaceWithNode might copy it to all.
    path.ReplaceWithNode(file, nodeToWrite)

    fmt.Println("Resulting YAML:")
    fmt.Println(file.String())
}
