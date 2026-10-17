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
	path, _ := goyaml.PathString("$.image")

    node, _ := path.FilterNode(file.Docs[0].Body)
    fmt.Printf("FilterNode returned: %T\n", node)

    fmt.Printf("GetComment: %+v\n", node.GetComment())
}
