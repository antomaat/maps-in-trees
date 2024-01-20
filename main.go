package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
)

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
    for i := 0; i < len(items); i++ {
        fmt.Println(items[i].Name())
    }

}

