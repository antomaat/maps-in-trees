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
}

func main() {

    arguments := os.Args[1:]
    if (len(arguments) == 0) {
        fmt.Println("missing directory path")
        os.Exit(1) 
    }

    dir := arguments[0]

    fmt.Println("Hello World")
    items, _ := os.ReadDir(dir)

    if (len(items) == 0) {
        fmt.Printf("directory %s empty or missing \n", dir)
        os.Exit(1)
    }

    createTreemap(items);
}

func createTreemap(items []os.DirEntry) {
    node := Node {
        Name: "root",
        PositionX: 0,
        PositionY: 0,
        SizeX: 25,
        SizeY: 25,
    }
    
    x := 0
    y := 0

    for i := 0; i < len(items); i++ {
        child := Node { 
            Name: items[i].Name(),
            PositionX: x + 25 * i,
            PositionY: y,
            SizeX: 25,
            SizeY: 25,
        }
        node.Children = append(node.Children, child)
        node.SizeX += child.SizeX
    }

    result, _ := json.Marshal(node)
    fmt.Println(string(result))

    _ = ioutil.WriteFile("./render/input.json", result, 0644)
}

