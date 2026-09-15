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

    tmpYAMLFile := ast.File{
		Name: file.Name,
	}
	tmpYAMLFile.Docs = append(tmpYAMLFile.Docs, file.Docs[0])

    // updateNode logic from target_goyaml.go
    // it gets matched nodes using matchedNodes
    // `matched` contains all nodes that match.
    // If ReplaceWithNode uses `nodeToWrite` for all matches, can we have multiple different comments?

    // We can't use a single `nodeToWrite` if we want to preserve different comments for each match!
    // But `ReplaceWithNode` only takes a single `ast.Node`.
    // And it will replace ALL matches with that single `ast.Node`.
    // That means if we have multiple matches with DIFFERENT comments, they will all get the SAME comment!

    nodeToWrite, _ := goyaml.ValueToNode(map[string]interface{}{"foo": "bar"})
    path.ReplaceWithNode(&tmpYAMLFile, nodeToWrite)

    fmt.Println(file.Docs[0].Body.String())
}
