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
    Size int64
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
    output := "./input.json"
    if len(arguments) > 1 {
        output = arguments[1]
    }
    createTreemap(dir, output, 2000, 500);
}

func createTreemap(dir string, output string, areaX float64, areaY float64) {
    node := createTree(dir, dir)
    node.SizeX = areaX
    node.SizeY = areaY
    //fmt.Println(node)
    node = updateDisplay(node, 0)
    result, _ := json.Marshal(node)
    fmt.Println(string(result))

    _ = ioutil.WriteFile(output, result, 0644)
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

    items, _ := os.ReadDir(pathName)

    for i := 0; i < len(items); i++ {
        child := Node {}
        info, _ := items[i].Info()
        if items[i].IsDir() {
            fmt.Print(items[i].Name())
            fmt.Println(": is dir")
            child = createTree(items[i].Name(), dirName + "/" + items[i].Name())
        } else {
            child = Node { 
                Name: items[i].Name(),
                PositionX: 0,
                PositionY: 0,
                SizeX: 1,
                SizeY: 1,
                Size: info.Size(), 
                IsDir: false,
                Path: pathName,
            }
        }

        node.Children = append(node.Children, child)
        //node.SizeX += child.SizeX
        node.Size += child.Size;
    }

    sort.SliceStable(node.Children, func(i, j int) bool {
        return node.Children[i].Size > node.Children[j].Size
    })
    fmt.Println("sorted array")
    fmt.Println(node.Children)
   
    return node
}

func squarifyDisplay(node Node) {

}

func updateDisplay(node Node, level int) Node {
    directories := []int{}
    corner := Vector2 {
        x: node.PositionX, 
        y: node.PositionY,
    }
    scale := Vector2 {
        x: node.SizeX,
        y: node.SizeY,
    }

    fmt.Println("-------------------")
    fmt.Print("parent size: ")
    fmt.Println(node.Size)
    fmt.Print("parent scale: ")
    fmt.Println(scale)


    for i := 0; i < len(node.Children); i++ {

        child := node.Children[i]

        if level % 2 == 0 {
            fraction :=  float64(child.Size) / float64(node.Size)
            area := scale.x * float64(fraction) 

            child.SizeX = area
            child.SizeY = scale.y 

            child.PositionX = corner.x
            corner.x += child.SizeX
            fmt.Println(" new child ------")
            fmt.Print(" size: ")
            fmt.Println(child.Size)
            fmt.Print(" fraction: ")
            fmt.Println(fraction)
            fmt.Print(" area: ")
            fmt.Println(area)
        } else {
            fmt.Println("level is odd")
            fraction :=  float64(child.Size) / float64(node.Size)
            area := scale.y * float64(fraction) 

            child.SizeX = scale.x 
            child.SizeY = area 

            child.PositionY = corner.y
            corner.y += child.SizeY
        }

        node.Children[i] = child

        if child.IsDir {
            directories = append(directories, i)
        }
    }

    for j := 0; j < len(directories); j++ {
        index := directories[j]
        /*fmt.Println("=====================")
        fmt.Println(node.Children[index])
        fmt.Println(node.Children[index].Name)
        fmt.Printf("posx %f \n", node.Children[index].PositionX)
        fmt.Printf("posy %f \n", node.Children[index].PositionY)
        fmt.Printf("scalex %f \n", node.Children[index].SizeX)
        fmt.Printf("scaley %f \n", node.Children[index].SizeY)
        */
        node.Children[index] = updateDisplay(node.Children[index], level + 1)
        //fmt.Println(node.Children[index])
    }

    return node
}
