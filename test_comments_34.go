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

    // According to target_goyaml.go:
    // matched, missing, err = matchedNodes(node, key)
    var matched []ast.Node
    matched = append(matched, node) // just for this test

    nodeToWrite, _ := goyaml.ValueToNode(map[string]interface{}{})

    // Our fix in updateNode:
    if nodeToWrite.GetComment() == nil && len(matched) > 0 && matched[0].GetComment() != nil {
		nodeToWrite.SetComment(matched[0].GetComment())
	}

    tmpYAMLFile := ast.File{
		Name: file.Name,
	}
	tmpYAMLFile.Docs = append(tmpYAMLFile.Docs, file.Docs[0])

    path.ReplaceWithNode(&tmpYAMLFile, nodeToWrite)

    fmt.Print(file.String())
}
