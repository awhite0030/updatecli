package main

import (
	"fmt"

	goyaml "github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/ast"
	"github.com/goccy/go-yaml/parser"
)

func cloneNode(n ast.Node) ast.Node {
	// Let's just create a new one to be safe, or just use ValueToNode again
	return nil
}

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

	node1, _ := path1.FilterNode(file.Docs[0].Body)
	if node1 != nil && node1.GetComment() != nil {
        comment := node1.GetComment()
		nodeToWrite1.SetComment(comment)
        // Can we preserve spacing?
	}
	path1.ReplaceWithNode(file, nodeToWrite1)

	fmt.Println(file.String())
}
