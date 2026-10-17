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

    // updateNode logic in target_goyaml.go
    // tmpYAMLFile.Docs = append(tmpYAMLFile.Docs, doc)
    // if err := urlPath.ReplaceWithNode(&tmpYAMLFile, nodeToWrite); err != nil {

    tmpYAMLFile := ast.File{
		Name: file.Name,
	}
	tmpYAMLFile.Docs = append(tmpYAMLFile.Docs, file.Docs[0])

    nodeToWrite, _ := goyaml.ValueToNode(map[string]interface{}{})

    // Find all matched nodes and preserve comments?
    node, _ := path.FilterNode(tmpYAMLFile.Docs[0].Body)
    if node != nil {
        if node.GetComment() != nil {
            nodeToWrite.SetComment(node.GetComment())
        }
    }

    path.ReplaceWithNode(&tmpYAMLFile, nodeToWrite)

    fmt.Println(file.Docs[0].Body.String())
}
