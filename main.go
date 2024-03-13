package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"slices"
)

type Node struct {
    Name string
    PositionX float64 
    PositionY float64 
    SizeX float64 
    SizeY float64 
    Size int64
    Children []Node
    IsDir bool
    Path string
    IsModule bool
    OptionalInfo OptionalInfo
}

type OptionalInfo struct {
    Fields []string
    FieldClasses []string
    SuperClass string
}

type Vector2 struct {
    x float64 
    y float64 
}

type Rectangle struct {
    x float64
    y float64
    width float64
    height float64
    index int
    size float64
}


func main() {

    arguments := os.Args[1:]
    if (len(arguments) == 0) {
        fmt.Println("missing directory path")
        os.Exit(1) 
    }


    dir := arguments[0]
    output := "./input.json"
    isJvm := true 
    if len(arguments) > 1 {
        output = arguments[1]
    }
    createTreemap(dir, output, 1000, 500, isJvm);
}

func isIgnoredDir(name string) bool {
    ignoredList := []string{".git", "build", "gradle", ".gradle", ".idea"}
    return slices.Contains(ignoredList, name)
}

func createTreemap(dir string, output string, areaX float64, areaY float64, isJvm bool) {
    //node := CreateTree(dir, dir, isJvm)
    node := CreateJvmTree(dir, dir, isJvm, dir)
    node.SizeX = areaX
    node.SizeY = areaY
    debug, _ := json.Marshal(node)
    _ = ioutil.WriteFile("debug.json", debug, 0644)
    //fmt.Println(node)
    //node = updateDisplay(node, 0)
    node = SquarifyDisplay(node)
    result, error := json.Marshal(node)
    if error != nil {
        fmt.Println(error.Error())
    }

    _ = ioutil.WriteFile(output, result, 0644)
}

