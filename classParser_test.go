package main

import (
	"os"
	"testing"
)

func TestFileParsingWithVariables(t *testing.T) {
    testFile, err := os.ReadFile("./files/Variables.class")
    if err !=nil {
        t.Fatal("no file named Main.class")
    }
    ParseFileInfo(testFile)
}

func TestFileMoreComplexExample(t *testing.T) {
    testFile, err := os.ReadFile("./files/Complex2.class")
    if err !=nil {
        t.Fatal("no file named Complex.class")
    }
    ParseFileInfo(testFile)
}
