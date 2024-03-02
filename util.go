package main 

import(
    "fmt"
)

func PrintNode(node Node) {
    fmt.Print("node: ")
    fmt.Println(node.Name)
    fmt.Print("x: ")
    fmt.Println(node.PositionX)
    fmt.Print("y: ")
    fmt.Println(node.PositionY)
    fmt.Print("width: ")
    fmt.Println(node.SizeX)
    fmt.Print("height: ")
    fmt.Println(node.SizeY)
}
