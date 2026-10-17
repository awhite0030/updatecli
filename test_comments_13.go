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

    // In goccy/go-yaml, FilterNode returns matched nodes wrapped.
    // updatecli uses matchedNodes(node, key) to extract.
    // Let's just use ValueToNode and preserve comment.

    nodeToWrite, _ := goyaml.ValueToNode(map[string]interface{}{})

    // We get the matched node from FilterNode:
    // It is `node`

    fmt.Printf("nodeToWrite: %T\n", nodeToWrite)

    // We want to apply comments of ALL matched nodes to the newly replaced node.
    // Since ReplaceWithNode replaces ALL matched nodes by `nodeToWrite`,
    // maybe we can extract comments from original nodes and apply them after ReplaceWithNode?
    // Wait, ReplaceWithNode will just inject `nodeToWrite`. If we modify `nodeToWrite`'s comment, it applies to all matches?
    // That's what currently happens! If we don't set comment, we lose it.

    // Let's look at `updateNode` in target_goyaml.go
}
