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

    nodeToWrite, _ := goyaml.ValueToNode(map[string]interface{}{})

    // So all we need is to copy comments from the replaced nodes in `updateNode`.
    // Wait, `ReplaceWithNode` receives `nodeToWrite`. But we might be replacing multiple matched nodes!
    // Since `ReplaceWithNode` uses `nodeToWrite` for all matches, can we have different comments per match?
    // Oh, if we use `ReplaceWithNode`, it replaces all nodes matching `urlPath` in `tmpYAMLFile.Docs[0].Body`.
    // BUT in `updatecli`, does a single `yaml.updateNode` handle multiple matches at once?
    // "matched holds every node the key selected, which is more than one for a wildcard or a recursive selector... ReplaceWithNode rewrites all of them at once."

    // If there are multiple matches, we can't preserve different comments for each match if we just use `ReplaceWithNode` with a single `nodeToWrite`.
    // BUT what happens if we replace node by node?
    // goccy doesn't provide a way to replace a specific node reference. It only has `ReplaceWithNode` which takes a `*ast.File` or `*ast.Node` and does it based on `Path`.
    // Let's check `ReplaceWithNode` in goccy.

}
