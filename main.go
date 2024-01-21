package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Node  struct {
    Name string
    PositionX int
    PositionY int
    SizeX int
    SizeY int
    Children []Node
    IsDir bool
}

func main() {

    arguments := os.Args[1:]
    if (len(arguments) == 0) {
        fmt.Println("missing directory path")
        os.Exit(1) 
    }

    dir := arguments[0]
    createTreemap(dir);
}

func createTreemap(dir string) {
    node := createTree(dir, 0, 0)
    result, _ := json.Marshal(node)
    fmt.Println(string(result))

    _ = ioutil.WriteFile("./render/input.json", result, 0644)
}

func createTree(dirName string, startX int, startY int) Node {
    node := Node {
        Name: dirName,
        PositionX: startX,
        PositionY: startY,
        SizeX: 25,
        SizeY: 25,
        IsDir: true,
    }

    items, _ := os.ReadDir(dirName)

    for i := 0; i < len(items); i++ {
        child := Node {}
        if items[i].IsDir() {
            child = createTree(items[i].Name(), startX, startY)
        } else {
            child = Node { 
                Name: items[i].Name(),
                PositionX: startX + 25 * i,
                PositionY: startY,
                SizeX: 25,
                SizeY: 25,
                IsDir: false,
            }
        }

        node.Children = append(node.Children, child)
        node.SizeX += child.SizeX
    }

    return node
}

