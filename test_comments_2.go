package main

import (
	"fmt"

	goyaml "github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/parser"
)

func main() {
	yamlStr := `emptyDir: {}     # some comments
image: mariadb:11.4.3   # keep aligned
`
	file, err := parser.ParseBytes([]byte(yamlStr), parser.ParseComments)
	if err != nil {
		panic(err)
	}

	nodeToWrite1, _ := goyaml.ValueToNode(map[string]interface{}{})
	path1, _ := goyaml.PathString("$.emptyDir")
	path1.ReplaceWithNode(file, nodeToWrite1)

	nodeToWrite2, _ := goyaml.ValueToNode("mariadb:11.8.2@sha256:123")
	path2, _ := goyaml.PathString("$.image")
	path2.ReplaceWithNode(file, nodeToWrite2)

	fmt.Println(file.String())
}
