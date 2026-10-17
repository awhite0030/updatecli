package main

import (
	"crypto/sha256"
	"fmt"
)

func main() {
	// simulate ManifestID()
	manifestID := "seed/manifest.yaml"
	actionKey := "default"

	fmt.Printf("%x\n", sha256.Sum256([]byte(manifestID+actionKey)))
}
