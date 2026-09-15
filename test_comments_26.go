package main

import (
	"fmt"
	"github.com/goccy/go-yaml/parser"
)

func main() {
	yamlStr := `something: else      # comment after multiples spaces
emptyDir: {} # some comments
image: mariadb:11.4.3   # keep aligned
`
	file, _ := parser.ParseBytes([]byte(yamlStr), parser.ParseComments)
    fmt.Print(file.String())
}
