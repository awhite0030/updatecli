package xml

import (
	"fmt"
	"strings"

	"github.com/beevik/etree"
)

// getElementAndAttribute returns the element path and the attribute name,
// if the path ends with an attribute selector (e.g., /@attrName).
func getElementAndAttribute(fullPath string) (string, string, bool) {
	if i := strings.LastIndex(fullPath, "/@"); i != -1 {
		return fullPath[:i], fullPath[i+2:], true
	}
	return fullPath, "", false
}

// queryElement returns the element text or attribute value.
// It returns a boolean indicating if the target (element or attribute) was found.
func queryElement(doc *etree.Document, fullPath string) (string, bool) {
	elemPath, attrName, hasAttr := getElementAndAttribute(fullPath)
	elem := doc.FindElement(elemPath)
	if elem == nil {
		return "", false
	}

	if hasAttr {
		attr := elem.SelectAttr(attrName)
		if attr != nil {
			return attr.Value, true
		}
		return "", false
	}

	return elem.Text(), true
}

// setElement updates the text or attribute of the element.
// Returns an error if the element doesn't exist.
func setElement(doc *etree.Document, fullPath, value string) error {
	elemPath, attrName, hasAttr := getElementAndAttribute(fullPath)
	elem := doc.FindElement(elemPath)
	if elem == nil {
		return fmt.Errorf("nothing found at path %q", elemPath)
	}

	if hasAttr {
		elem.CreateAttr(attrName, value)
	} else {
		elem.SetText(value)
	}
	return nil
}
