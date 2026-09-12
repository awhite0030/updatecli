package hcl

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-viper/mapstructure/v2"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/minamijoyo/hcledit/editor"
	"github.com/sirupsen/logrus"
	"github.com/updatecli/updatecli/pkg/core/result"
	"github.com/updatecli/updatecli/pkg/core/text"
	"github.com/updatecli/updatecli/pkg/plugins/utils"
)

type Hcl struct {
	spec             Spec
	contentRetriever text.TextRetriever
	files            map[string]file // map of file paths to file contents
}

type file struct {
	originalFilePath string
	filePath         string
	content          string
}

func New(spec interface{}) (*Hcl, error) {
	newSpec := Spec{}

	err := mapstructure.Decode(spec, &newSpec)
	if err != nil {
		return nil, err
	}

	newResource := &Hcl{
		spec:             newSpec,
		contentRetriever: &text.Text{},
	}

	err = newResource.spec.Validate()
	if err != nil {
		return nil, err
	}

	newResource.files = make(map[string]file)
	// File as unique element of newResource.files
	if len(newResource.spec.File) > 0 {
		filePath := strings.TrimPrefix(newResource.spec.File, "file://")
		newResource.files[filePath] = file{
			originalFilePath: filePath,
			filePath:         filePath,
		}
	}
	// Files
	for _, filePath := range newResource.spec.Files {
		filePath := strings.TrimPrefix(filePath, "file://")
		newResource.files[filePath] = file{
			originalFilePath: filePath,
			filePath:         filePath,
		}
	}

	return newResource, nil
}

func (h *Hcl) Query(resourceFile file) (string, error) {
	query := h.spec.Path

	basePath, mapKey := parseHclQuery(query)

	sink := editor.NewAttributeGetSink(basePath, false) // We don't need comments for query

	inStream := strings.NewReader(resourceFile.content)
	outStream := new(bytes.Buffer)
	err := editor.DeriveStream(inStream, outStream, resourceFile.filePath, sink)
	if err != nil {
		return "", err
	}

	if outStream.Len() == 0 {
		return "", fmt.Errorf("%s cannot find value for path %q from file %q",
			result.FAILURE,
			query,
			resourceFile.originalFilePath)
	}

	if mapKey != "" {
		attrValueStr := outStream.String()
		dummySrc := []byte("dummy = " + attrValueStr)
		f, diags := hclwrite.ParseConfig(dummySrc, "", hcl.InitialPos)
		if diags.HasErrors() {
			return "", fmt.Errorf("%s cannot parse value for path %q from file %q",
				result.FAILURE,
				basePath,
				resourceFile.originalFilePath)
		}

		attr := f.Body().GetAttribute("dummy")
		if attr == nil {
			return "", fmt.Errorf("%s cannot extract dummy attribute for path %q from file %q",
				result.FAILURE,
				basePath,
				resourceFile.originalFilePath)
		}

		tokens := attr.Expr().BuildTokens(nil)
		startIdx, endIdx := findMapValueBounds(tokens, mapKey)

		if startIdx == -1 || endIdx == -1 {
			return "", fmt.Errorf("%s cannot find key %q in map for path %q from file %q",
				result.FAILURE,
				mapKey,
				basePath,
				resourceFile.originalFilePath)
		}

		// Extract value as string
		var valBytes []byte
		for i := startIdx; i < endIdx; i++ {
			valBytes = append(valBytes, tokens[i].Bytes...)
		}
		queryOutput := strings.TrimSpace(string(valBytes))
		queryOutput = strings.Trim(queryOutput, "\"'")
		return queryOutput, nil
	}

	queryOutput := outStream.String()
	queryOutput = strings.TrimSpace(queryOutput)
	queryOutput = strings.Trim(queryOutput, "\"'")
	return queryOutput, nil
}

func (h *Hcl) Apply(filePath string, valueToWrite string) error {
	query := h.spec.Path

	basePath, mapKey := parseHclQuery(query)

	if _, err := strconv.Atoi(valueToWrite); err != nil {
		valueToWrite = fmt.Sprintf(`"%s"`, valueToWrite)
	}

	resourceFile := h.files[filePath]

	if mapKey != "" {
		// First, get the current map string
		sink := editor.NewAttributeGetSink(basePath, true)
		inStream := strings.NewReader(resourceFile.content)
		outStream := new(bytes.Buffer)
		err := editor.DeriveStream(inStream, outStream, resourceFile.filePath, sink)
		if err != nil {
			return err
		}

		if outStream.Len() == 0 {
			return fmt.Errorf("%s cannot find map for path %q from file %q", result.FAILURE, basePath, resourceFile.originalFilePath)
		}

		attrValueStr := outStream.String()
		dummySrc := []byte("dummy = " + attrValueStr)
		f, diags := hclwrite.ParseConfig(dummySrc, "", hcl.InitialPos)
		if diags.HasErrors() {
			return fmt.Errorf("%s cannot parse map for path %q from file %q", result.FAILURE, basePath, resourceFile.originalFilePath)
		}

		attr := f.Body().GetAttribute("dummy")
		if attr == nil {
			return fmt.Errorf("%s cannot extract dummy attribute for path %q from file %q", result.FAILURE, basePath, resourceFile.originalFilePath)
		}

		tokens := attr.Expr().BuildTokens(nil)
		startIdx, endIdx := findMapValueBounds(tokens, mapKey)

		if startIdx == -1 || endIdx == -1 {
			return fmt.Errorf("%s cannot find key %q in map for path %q from file %q", result.FAILURE, mapKey, basePath, resourceFile.originalFilePath)
		}

		// Parse the new value to get its tokens
		valSrc := []byte("dummy = " + valueToWrite)
		vf, diags := hclwrite.ParseConfig(valSrc, "", hcl.InitialPos)
		if diags.HasErrors() {
			return fmt.Errorf("%s cannot parse new value %q", result.FAILURE, valueToWrite)
		}

		newValTokens := vf.Body().GetAttribute("dummy").Expr().BuildTokens(nil)

		// Replace the tokens in the map
		newTokens := append(tokens[:startIdx], append(newValTokens, tokens[endIdx:]...)...)

		f.Body().SetAttributeRaw("dummy", newTokens)

		// Extract the updated map string
		newDummyBytes := f.Body().GetAttribute("dummy").Expr().BuildTokens(nil).Bytes()
		newAttrStr := string(newDummyBytes)

		// Use hcledit to set the map string back
		filter := editor.NewAttributeSetFilter(basePath, newAttrStr)
		inStream = strings.NewReader(resourceFile.content)
		outStream = new(bytes.Buffer)
		err = editor.EditStream(inStream, outStream, resourceFile.filePath, filter)
		if err != nil {
			return err
		}

		resourceFile.content = outStream.String()
		h.files[filePath] = resourceFile

		return nil
	}

	filter := editor.NewAttributeSetFilter(basePath, valueToWrite)
	inStream := strings.NewReader(resourceFile.content)
	outStream := new(bytes.Buffer)
	err := editor.EditStream(inStream, outStream, resourceFile.filePath, filter)
	if err != nil {
		return err
	}

	resourceFile.content = outStream.String()

	h.files[filePath] = resourceFile

	return nil
}

// Read puts the content of the file(s) as value of the y.files map if the file(s) exist(s) or log the non existence of the file
func (h *Hcl) Read() error {
	var err error

	// Retrieve files content
	for filePath := range h.files {
		f := h.files[filePath]
		if h.contentRetriever.FileExists(f.filePath) {
			f.content, err = h.contentRetriever.ReadAll(f.filePath)
			if err != nil {
				return err
			}
			h.files[filePath] = f

		} else {
			return fmt.Errorf("%s The specified file %q does not exist", result.FAILURE, f.filePath)
		}
	}
	return nil
}

func (h *Hcl) UpdateAbsoluteFilePath(workDir string) {
	for filePath := range h.files {
		if workDir != "" {
			f := h.files[filePath]
			f.filePath = utils.JoinFilePathWithWorkingDirectoryPath(f.originalFilePath, workDir)
			logrus.Debugf("Relative path detected: changing from %q to absolute path from SCM: %q", f.originalFilePath, f.filePath)
			h.files[filePath] = f
		}
	}
}

// Changelog returns the changelog for this resource, or an empty string if not supported
func (h *Hcl) Changelog(from, to string) *result.Changelogs {
	return nil
}

// ReportConfig return a new configuration without any sensitive information
// or context specific data.
func (h *Hcl) ReportConfig() interface{} {
	return Spec{
		File:  h.spec.File,
		Files: h.spec.Files,
		Path:  h.spec.Path,
		Value: h.spec.Value,
	}
}

// parseHclQuery splits an HCL query with a trailing map indexer into a base path and a map key.
// It supports map keys with double quotes, single quotes, or unquoted identifiers.
func parseHclQuery(query string) (basePath string, mapKey string) {
	if strings.HasSuffix(query, "]") {
		idx := strings.LastIndex(query, "[")
		if idx != -1 {
			basePath = query[:idx]
			mapKey = query[idx+1 : len(query)-1]
			// Trim quotes (double or single)
			mapKey = strings.Trim(mapKey, `"'`)
			return basePath, mapKey
		}
	}
	return query, ""
}

func findMapValueBounds(tokens hclwrite.Tokens, keyToFind string) (startIndex int, endIndex int) {
	depth := 0
	foundKey := false
	startIndex = -1
	endIndex = -1

	for i := 0; i < len(tokens); i++ {
		t := tokens[i]

		if t.Type == hclsyntax.TokenOBrace || t.Type == hclsyntax.TokenOBrack || t.Type == hclsyntax.TokenOParen {
			depth++
		}
		if t.Type == hclsyntax.TokenCBrace || t.Type == hclsyntax.TokenCBrack || t.Type == hclsyntax.TokenCParen {
			depth--
		}

		if !foundKey && depth == 1 {
			// At depth 1 inside the object { ... }
			isKeyMatch := false
			if t.Type == hclsyntax.TokenQuotedLit && string(t.Bytes) == keyToFind {
				isKeyMatch = true
			} else if t.Type == hclsyntax.TokenIdent && string(t.Bytes) == keyToFind {
				isKeyMatch = true
			}

			if isKeyMatch {
				foundKey = true
				// Find the equal or colon
				for j := i + 1; j < len(tokens); j++ {
					if tokens[j].Type == hclsyntax.TokenEqual || tokens[j].Type == hclsyntax.TokenColon {
						startIndex = j + 1
						break
					}
				}
				if startIndex != -1 {
					i = startIndex - 1 // skip to start index
					continue
				}
			}
		}

		if foundKey && startIndex != -1 && endIndex == -1 {
			// Find the end of the value
			if (depth == 1 && (t.Type == hclsyntax.TokenComma || t.Type == hclsyntax.TokenNewline)) || (depth == 0 && t.Type == hclsyntax.TokenCBrace) {
				endIndex = i
				break
			}
		}
	}

	if foundKey && startIndex != -1 && endIndex != -1 {
		return startIndex, endIndex
	}

	return -1, -1
}
