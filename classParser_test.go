package main

import (
	"os"
	"testing"
)

func TestFileParsing(t *testing.T) {
    testFile, err := os.ReadFile("Main.class")
    if err !=nil {
        t.Fatal("node height is messed up")
    }
    ParseFileInfo(testFile)
}
