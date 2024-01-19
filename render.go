package main

import (
	"fmt"

	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/pixelgl"
)

func run() {
	result := pixelgl.WindowConfig{
		Title:  "Pixel Rocks!",
		Bounds: pixel.R(0, 0, 1024, 768),
	}
	win, err := pixelgl.NewWindow(result)

	if err != nil {
		panic(err)
	}

	for !win.Closed() {
		win.Update()
	}
}

func InitRender() {
	fmt.Println("init render")
	pixelgl.Run(run)
}
