package main

import "fmt"
import "crypto/sha256"

func main() {
    fmt.Printf("%x", sha256.Sum256([]byte("seed")))
}
