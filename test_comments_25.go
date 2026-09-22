package main

import (
	"fmt"
	goyaml "github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/parser"
    "github.com/goccy/go-yaml/ast"
)

func main() {
	yamlStr := `volumes:
- emptyDir: {}       # some comments
  name: a
- emptyDir: {}       # some other comments
  name: b
`
	file, _ := parser.ParseBytes([]byte(yamlStr), parser.ParseComments)
	path, _ := goyaml.PathString("$.volumes[*].emptyDir")

    node, _ := path.FilterNode(file.Docs[0].Body)

    // Instead of doing `ReplaceWithNode`, can we just replace the contents of the nodes in `matchedNodes`?
    var matched []ast.Node
    if mn, ok := node.(*ast.SequenceNode); ok {
        for _, val := range mn.Values {
            matched = append(matched, val)
        }
    }

    // Actually in `updatecli/pkg/plugins/resources/yaml/utils.go` they have `matchedNodes` function. Let's use it.

}
