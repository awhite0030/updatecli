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

    // Instead of doing ReplaceWithNode(..., nodeToWrite), can we update the values?
    // According to updatecli guidelines in the memory prompt:
    // "When modifying YAML files and exact formatting (like blank lines) needs to be preserved,
    // use `github.com/goccy/go-yaml` (specifically its AST token manipulation) instead of `go.yaml.in/yaml/v3` which strips blank lines during encoding."

    fmt.Printf("node: %T\n", node)

    if mn, ok := node.(*ast.MappingNode); ok {
       for _, val := range mn.Values {
           fmt.Printf("key: %s\n", val.Key.String())
           fmt.Printf("value: %T, %s\n", val.Value, val.Value.String())
       }
    }
}
