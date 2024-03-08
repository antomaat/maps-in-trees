package main

import (
	"fmt"
	"os"
	"testing"
)

func TestFileParsingWithVariables(t *testing.T) {
    testFile, err := os.ReadFile("./files/Variables.class")
    if err !=nil {
        t.Fatal("no file named Variables.class")
    }
    result := ParseFileInfo(testFile)
    fmt.Println(result.fields)
    if len(result.fields) == 0 {
        t.Fatal("no fields resulted")
    }
}

func TestFileMoreComplexExample(t *testing.T) {
    testFile, err := os.ReadFile("./files/Complex2.class")
    if err !=nil {
        t.Fatal("no file named Complex.class")
    }
    result := ParseFileInfo(testFile)
    fmt.Println(result.fields)
}
