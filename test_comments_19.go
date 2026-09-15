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

    path, _ := goyaml.PathString("$.emptyDir")
    node, _ := path.FilterNode(file.Docs[0].Body)

    // Modifying value in-place without replacing the whole node?
    // Wait, the issue says:
    // "actions that delete comments are not [tolerable]"
    // But they say spaces being trimmed is tolerable in some cases!
    // "actions that trims some spaces between code and comments are tolerable"
    // "actions that delete comments are not"
    // "Some comments a purely erased; it occurs when a comments is on the same line as {}"

    // Can we just set the comment to the new nodeToWrite in updateNode?

    nodeToWrite, _ := goyaml.ValueToNode(map[string]interface{}{})

    // In updatecli they replace node:
    if node != nil && node.GetComment() != nil {
        nodeToWrite.SetComment(node.GetComment())
    }

    path.ReplaceWithNode(file, nodeToWrite)

    fmt.Println("Resulting YAML:")
    fmt.Println(file.String())
}
