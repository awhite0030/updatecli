package main

import (
	"fmt"
	goyaml "github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/parser"
)

func main() {
	yamlStr := `emptyDir: {}       # some comments
image: mariadb:11.4.3   # keep aligned
`
	file, _ := parser.ParseBytes([]byte(yamlStr), parser.ParseComments)
	nodeToWrite, _ := goyaml.ValueToNode(map[string]interface{}{})
	path, _ := goyaml.PathString("$.emptyDir")

    node, _ := path.FilterNode(file.Docs[0].Body)
    if node != nil && node.GetComment() != nil {
        comment := node.GetComment()
        nodeToWrite.SetComment(comment)
    }

    path.ReplaceWithNode(file, nodeToWrite)

	fmt.Println(file.String())

    // Now test with ast.Node manipulation
    file2, _ := parser.ParseBytes([]byte(yamlStr), parser.ParseComments)
    node2, _ := path.FilterNode(file2.Docs[0].Body)
    fmt.Printf("node2 token: %+v\n", node2.GetToken())
    fmt.Printf("node2 comment: %+v\n", node2.GetComment())
    if node2.GetComment() != nil {
       fmt.Printf("node2 comment token: %+v\n", node2.GetComment().GetToken())
    }
}
