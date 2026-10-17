package main

import (
	"fmt"
	goyaml "github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/parser"
    "github.com/goccy/go-yaml/ast"
)

func main() {
	yamlStr := `emptyDir: {}       # some comments
image: mariadb:11.4.3   # keep aligned
`
	file, _ := parser.ParseBytes([]byte(yamlStr), parser.ParseComments)
	path, _ := goyaml.PathString("$.emptyDir")

    node, _ := path.FilterNode(file.Docs[0].Body)

    // Instead of ReplaceWithNode, can we modify the matched node directly?
    // No, we need to rewrite node with new value.

    // In updatecli they replace the whole node. But maybe we can reuse the ast package.
    // What about just getting the token from the old node and copying it to the new node?
    // Or copy comment tokens to the replaced nodes?

    nodeToWrite, _ := goyaml.ValueToNode(map[string]interface{}{})

    if node != nil && node.GetComment() != nil {
        // Can we do anything about the comment spacing?
        // Wait, issue says "The spaces are changed between the "code" and the "comments left by author" ...
        // Some comments a purely erased; it occurs when a comments is on the same line as {}"

        // Actually, in updatecli, they do NOT currently transfer the comments to the nodeToWrite.
        // Let's check updatecli code.
    }
}
