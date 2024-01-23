package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
)

type Node  struct {
    Name string
    PositionX float64 
    PositionY float64 
    SizeX float64 
    SizeY float64 
    Size int
    Children []Node
    IsDir bool
    Path string
}

type Vector2 struct {
    x float64 
    y float64 
}

func main() {

    arguments := os.Args[1:]
    if (len(arguments) == 0) {
        fmt.Println("missing directory path")
        os.Exit(1) 
    }

    dir := arguments[0]
    createTreemap(dir, 1000, 500);
}

func createTreemap(dir string, areaX float64, areaY float64) {
    node := createTree(dir, dir)
    //fmt.Println(node)
    node = updateDisplay(node, 0, 0, areaX, areaY, 1)
    result, _ := json.Marshal(node)
    fmt.Println(string(result))

    _ = ioutil.WriteFile("./render/input.json", result, 0644)
}

func createTree(dirName string, pathName string) Node {
    node := Node {
        Name: dirName,
        PositionX: 0,
        PositionY: 0,
        SizeX: 1,
        SizeY: 1,
        Size: 1,
        IsDir: true,
        Path: pathName,
    }

    items, _ := os.ReadDir(dirName)

    for i := 0; i < len(items); i++ {
        child := Node {}
        if items[i].IsDir() {
            child = createTree(items[i].Name(), dirName + "/" + items[i].Name())
        } else {
            child = Node { 
                Name: items[i].Name(),
                PositionX: 0,
                PositionY: 0,
                SizeX: 1,
                SizeY: 1,
                Size: 1,
                IsDir: false,
                Path: pathName,
            }
        }

        node.Children = append(node.Children, child)
        //node.SizeX += child.SizeX
        node.Size += 1;
    }

    sort.SliceStable(node.Children, func(i, j int) bool {
        return node.Children[i].Size < node.Children[i].Size
    })
   
    return node
}

func updateDisplay(node Node, positionX float64, positionY float64, areaX float64, areaY float64, level int) Node {
    directories := []int{}
    corner := Vector2 {
        x: positionX, 
        y: positionY,
    }

    for i := 0; i < len(node.Children); i++ {

        fraction :=  node.Children[i].Size / node.Size 
        area := areaX * float64(fraction) 
        fmt.Println(area)

        //node.Children[i].SizeX = area / node.SizeX
        node.Children[i].SizeY = areaY

        node.Children[i].PositionX = corner.x
        //node.Children[i].PositionY = corner.y
        corner.x += node.Children[i].SizeX
        //corner.y += node.Children[i].SizeY

        if node.Children[i].IsDir {
            directories = append(directories, i)
        }
    }

    for j := 0; j < len(directories); j++ {
        index := directories[j]
        fmt.Println("=====================")
        fmt.Println(node.Children[index])
        fmt.Println(node.Children[index].Name)
        fmt.Printf("posx %f \n", node.Children[index].PositionX)
        fmt.Printf("posy %f \n", node.Children[index].PositionY)
        fmt.Printf("scalex %f \n", node.Children[index].SizeX)
        fmt.Printf("scaley %f \n", node.Children[index].SizeY)
        node.Children[index] = updateDisplay(node.Children[index], node.Children[index].PositionX, node.Children[index].PositionY, node.Children[index].SizeX, node.Children[index].SizeY, level + 1)
        fmt.Println(node.Children[index])
    }



    return node
}
