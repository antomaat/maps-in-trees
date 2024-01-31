package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"
)

func TestSquarify(t *testing.T) {
    node := createNodes()
    rect := SquarifyDisplay(node)
    if rect.height != 4   {
        t.Fatal("node height is messed up")
    }
    if rect.width != 6 * (6 /  24)   {
        t.Fatal("node height is messed up")
    }
}


func TestSquarifyNew(t *testing.T) {
    fmt.Println("-------------------")
    fmt.Println("test squarified new")
    node := createNodes()
    result := NewSquarifyDisplay(node)
    fmt.Println(result.Children)
    saveResultToJson("./render/input.json", node)
}

func createNodes() Node {
    node := Node {
        Name: "root",
        PositionX: 0,
        PositionY: 0,
        SizeX: 1000,
        SizeY: 500,
        Size: 24,
        IsDir: true,
        Path: "path",
    }
    node.Children = append(node.Children, createNode(6, "first"))
    node.Children = append(node.Children, createNode(6, "second"))
    node.Children = append(node.Children, createNode(4, "third"))
    node.Children = append(node.Children, createNode(4, "fourth"))
    return node
}

func saveResultToJson(output string, node Node) {
    result, _ := json.Marshal(node)
    fmt.Println(string(result))

    _ = ioutil.WriteFile(output, result, 0644)
}
func createNode(size int64, name string) Node {
    node := Node {
        Name: name,
        PositionX: 0,
        PositionY: 0,
        SizeX: 1,
        SizeY: 1,
        Size: size,
        IsDir: false,
        Path: "file",
    }
    return node
}
