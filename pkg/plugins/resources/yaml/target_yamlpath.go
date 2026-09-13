package yaml

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/vmware-labs/yaml-jsonpath/pkg/yamlpath"
	"github.com/goccy/go-yaml/parser"
	"github.com/goccy/go-yaml/ast"

	"github.com/updatecli/updatecli/pkg/core/result"

	"go.yaml.in/yaml/v3"
)

func replaceToken(node ast.Node, line, col int, oldVal, newVal string, style yaml.Style, comment string) bool {
	if node == nil {
		return false
	}

	switch n := node.(type) {
	case *ast.StringNode:
		tok := n.GetToken()
		if tok != nil && tok.Position.Line == line && (tok.Position.Column == col || tok.Value == oldVal || strings.Trim(tok.Value, "\"") == oldVal || strings.Trim(tok.Value, "'") == oldVal) {
			n.Value = newVal
			switch style {
			case yaml.DoubleQuotedStyle:
				n.Token.Value = "\"" + newVal + "\""
			case yaml.SingleQuotedStyle:
				n.Token.Value = "'" + newVal + "'"
			default:
				n.Token.Value = newVal
			}
			if comment != "" {
				_ = setNodeComment(n, comment)
			}
			return true
		}
	case *ast.FloatNode:
		tok := n.GetToken()
		if tok != nil && tok.Position.Line == line && (tok.Position.Column == col || tok.Value == oldVal) {
			n.Value = 0
			n.Token.Value = newVal
			if comment != "" {
				_ = setNodeComment(n, comment)
			}
			return true
		}
	case *ast.IntegerNode:
		tok := n.GetToken()
		if tok != nil && tok.Position.Line == line && (tok.Position.Column == col || tok.Value == oldVal) {
			n.Value = 0
			n.Token.Value = newVal
			if comment != "" {
				_ = setNodeComment(n, comment)
			}
			return true
		}
	case *ast.BoolNode:
		tok := n.GetToken()
		if tok != nil && tok.Position.Line == line && (tok.Position.Column == col || tok.Value == oldVal) {
			n.Value = false
			n.Token.Value = newVal
			if comment != "" {
				_ = setNodeComment(n, comment)
			}
			return true
		}
	}

	switch n := node.(type) {
	case *ast.MappingNode:
		for _, v := range n.Values {
			if replaceToken(v.Key, line, col, oldVal, newVal, style, comment) { return true }
			if replaceToken(v.Value, line, col, oldVal, newVal, style, comment) { return true }
		}
	case *ast.SequenceNode:
		for _, v := range n.Values {
			if replaceToken(v, line, col, oldVal, newVal, style, comment) { return true }
		}
	case *ast.MappingValueNode:
		if replaceToken(n.Key, line, col, oldVal, newVal, style, comment) { return true }
		if replaceToken(n.Value, line, col, oldVal, newVal, style, comment) { return true }
	case *ast.DocumentNode:
		if replaceToken(n.Body, line, col, oldVal, newVal, style, comment) { return true }
	}
	return false
}

func (y *Yaml) goYamlPathTarget(valueToWrite string, resultTarget *result.Target, dryRun bool) (notChanged int, ignoredFiles int, err error) {
	keys := y.spec.getKeys()

	resultTargetFilesMap := map[string]bool{}

	for filePath := range y.files {
		originFilePath := y.files[filePath].originalFilePath
		fileNotChanged := 0
		fileKeysProcessed := 0

		// Decode the file into one or more YAML document nodes
		var docs []*yaml.Node
		dec := yaml.NewDecoder(strings.NewReader(y.files[filePath].content))
		for {
			var doc yaml.Node
			if derr := dec.Decode(&doc); derr != nil {
				if derr == io.EOF {
					break
				}
				return 0, ignoredFiles, fmt.Errorf("parsing yaml file %q: %w", originFilePath, derr)
			}
			docs = append(docs, &doc)
		}

		// If the file has no documents, treat as ignored
		if len(docs) == 0 {
			ignoredFiles++
			continue
		}

		// Also parse with goccy/go-yaml to preserve formatting
		goccyFile, err := parser.ParseBytes([]byte(y.files[filePath].content), parser.ParseComments)
		if err != nil {
			return 0, ignoredFiles, fmt.Errorf("parsing yaml file with goccy %q: %w", originFilePath, err)
		}

		// Process each key for this file
		for _, key := range keys {
			urlPath, err := yamlpath.NewPath(key)
			if err != nil {
				return 0, 0, fmt.Errorf("crafting yamlpath query for key %q: %w", key, err)
			}

			// No DocumentIndex: search across all documents and process each doc that matches
			foundAny := false
			var docNotChangedCount int
			// docsMatched counts the documents in which the key was found, which is
			// the only meaningful denominator: len(docs) ignores both DocumentIndex
			// and the documents that simply do not carry the key.
			var docsMatched int

			for index, doc := range docs {

				if y.spec.DocumentIndex != nil {
					if index != *y.spec.DocumentIndex {
						continue
					}
				}

				nodes, err := urlPath.Find(doc)
				if err != nil {
					return 0, ignoredFiles, fmt.Errorf("searching for key %q in yaml file: %w", key, err)
				}
				if len(nodes) == 0 {
					continue
				}

				foundAny = true
				docsMatched++
				var oldVersion string
				var notChangedNode int
				for _, node := range nodes {
					oldVersion = node.Value
					resultTarget.Information = oldVersion

					if oldVersion == valueToWrite {
						resultTarget.Description = fmt.Sprintf("%s\nkey %q already set to %q, from file %q (doc %d)",
							resultTarget.Description,
							key,
							valueToWrite,
							originFilePath,
							index)
						notChangedNode++
						continue
					}

					node.Value = valueToWrite
					if y.spec.Comment != "" {
						node.LineComment = y.spec.Comment
						// Set comment in goccy tree? Wait, there is already setNodeComment, but we only have string nodes here and the comment belongs to the node.
					}

					// Update goccy AST
					if index < len(goccyFile.Docs) {
						if !replaceToken(goccyFile.Docs[index], node.Line, node.Column, oldVersion, valueToWrite, node.Style, y.spec.Comment) {
							logrus.Debugf("could not find token to replace in goccy AST for line %d col %d", node.Line, node.Column)
						}
					}

					if _, ok := resultTargetFilesMap[filePath]; !ok {
						resultTarget.Files = append(resultTarget.Files, y.files[filePath].filePath)
						resultTargetFilesMap[filePath] = true
					}

					resultTarget.Changed = true
					resultTarget.Result = result.ATTENTION

					shouldMsg := " "
					if dryRun {
						shouldMsg = " should be "
					}

					resultTarget.Description = fmt.Sprintf("%s\nkey %q%supdated from %q to %q, in file %q (doc %d)",
						resultTarget.Description,
						key,
						shouldMsg,
						oldVersion,
						valueToWrite,
						originFilePath,
						index)
				}
				if notChangedNode == len(nodes) {
					docNotChangedCount++
				}
			}

			// Whether the key was found has to be decided once every document has
			// been evaluated, not while iterating over them.
			if !foundAny {
				if y.spec.SearchPattern {
					logrus.Debugf("ignoring key %q in file %q as we couldn't find it in any document", key, originFilePath)
					continue
				}
				return 0, ignoredFiles, fmt.Errorf("couldn't find key %q from file %q", key, originFilePath)
			}

			fileKeysProcessed++

			if docNotChangedCount == docsMatched {
				// if every matching document had only unchanged nodes, consider file not changed for this key
				fileNotChanged++
			}
		} // end keys loop

		// If no keys were processed for this file (all were ignored), count as ignored
		if fileKeysProcessed == 0 {
			ignoredFiles++
			continue
		}

		// If all processed keys in this file were unchanged, count as not changed
		if fileNotChanged == fileKeysProcessed {
			notChanged++
		}

		f := y.files[filePath]
		f.content = goccyFile.String()
		// preserve leading document marker if it was present originally
		if strings.HasPrefix(y.files[filePath].content, "---\n") && !strings.HasPrefix(f.content, "---\n") {
			f.content = "---\n" + f.content
		}
		y.files[filePath] = f

		if !dryRun {
			newFile, err := os.Create(y.files[filePath].filePath)
			if err != nil {
				return 0, ignoredFiles, fmt.Errorf("creating file %q: %w", originFilePath, err)
			}
			defer newFile.Close()

			err = y.contentRetriever.WriteToFile(
				y.files[filePath].content,
				y.files[filePath].filePath)

			if err != nil {
				return 0, ignoredFiles, fmt.Errorf("saving file %q: %w", originFilePath, err)
			}
		}
	}
	return notChanged, ignoredFiles, nil
}
