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

        // Check token of comment
        cToken := comment.GetToken()
        if cToken != nil {
            fmt.Printf("origin: %q\n", cToken.Origin)
            fmt.Printf("value: %q\n", cToken.Value)
        }

        nodeToWrite.SetComment(comment)
    }

    path.ReplaceWithNode(file, nodeToWrite)
	fmt.Println(file.String())
}
