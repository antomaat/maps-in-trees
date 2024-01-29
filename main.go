package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
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

type Rectangle struct {
    x float64
    y float64
    width float64
    height float64
    index int
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
    //node = updateDisplay(node, 0)
    node = NewSquarifyDisplay(node)
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


func NewSquarifyDisplay(node Node) Node {
    fillArea := Rectangle {
        x: node.PositionX,
        y: node.PositionY,
        width: node.SizeX,
        height: node.SizeY,
    }
    //fullSize := node.Size
    //vertical := fillArea.height <= fillArea.width 
    //children := node.Children
    directories := []int{}
    cache := [][]Node{}
    row := []Node{}
    for _, n := range node.Children {

        vertical := fillArea.height <= fillArea.width 
        widthN := fillArea.width
        if vertical {
            widthN = fillArea.height
        }
        rowWithChild := append(row, n)

        if len(row) == 0 || worst(row, widthN) >= worst(rowWithChild, widthN) {
            row = append(row, n)
        } else {
            row = layoutRow(row, widthN, vertical, &fillArea)
            cache = append(cache, row)
            // TODO: remove the added area from the fullArea
            row = []Node{}
            row = append(row, n)
        }
    }

    // maybe a single node is left unprocessed
    if len(row) > 0 {
        row = layoutRow(row, fillArea.height, true, &fillArea)
        cache = append(cache, row)
        row = []Node{}
    }

    fmt.Println(cache)

    // update references
    index := 0
    for _, nList := range cache {
        for _, n := range nList {
            if n.Name != node.Children[index].Name {
                panic("holy hell, a element is calculated wrongly")
            } 
            node.Children[index].PositionX = n.PositionX
            node.Children[index].PositionY = n.PositionY
            node.Children[index].SizeX = n.SizeX
            node.Children[index].SizeY = n.SizeY

            if n.IsDir {
                directories = append(directories, index)
            }
            index++
        }
    }

    for j := 0; j < len(directories); j++ {
        indx := directories[j]
        node.Children[indx] = NewSquarifyDisplay(node.Children[indx])
    }
    return node
}

func layoutRow(row []Node, width float64, vertical bool, parent *Rectangle) []Node {
    rowHeight := sumSizes(row) / width
    result := []Node{}
    cacheParent := Rectangle {
        x: parent.x,
        y: parent.y,
        width: parent.width,
        height: parent.height,
    }
    for _, node := range row {
        rowWidth := float64(node.Size) / rowHeight
        x := cacheParent.x
        y := cacheParent.y

        fmt.Println("=================")
        fmt.Println("node at beginning")
        printNode(node)

        if (vertical) {
            node.PositionX = x
            node.PositionY = y
            node.SizeY = float64(rowWidth)
            node.SizeX = rowHeight
            cacheParent.y += float64(rowWidth)
        } else {
            node.PositionX = x
            node.PositionY = y
            node.SizeX = float64(rowWidth)
            node.SizeY = rowHeight
            cacheParent.x += float64(rowWidth)
        }

        fmt.Println("---------------")
        fmt.Println("node at the end")
        printNode(node)

        if vertical {
            parent.x += rowHeight
            parent.width -= rowHeight
            parent.height -= rowWidth
        } else {
            parent.y += rowHeight
            parent.width -= rowWidth
            parent.height -= rowHeight
        }

        result = append(result, node)
    }
    return result 
}

func printNode(node Node) {
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

func worst(sizes []Node, w float64) float64 {
	max := math.Inf(-1)
	min := math.Inf(1)
	sum := 0.0
	for _, size := range sizes {
		sum += float64(size.Size)
		max = math.Max(max, float64(size.Size))
		min = math.Min(min, float64(size.Size))
	}
	return math.Max((w*w*max)/(sum*sum), (sum*sum)/(w*w*min))
}

/*func layoutCol(nodes []Node, rect Rectangle) {
    covered := sumSizes(nodes)
    height := covered / rect.width
    calcNodes:= []
    for i, node := range nodes {
        calcNodes = append(calcNodes, node)
    }
}*/

func sumSizes(nodes []Node) float64 {
    sum := 0.0
    for _, n := range nodes {
        sum += float64(n.Size)
    }
    return sum
}

func SquarifyDisplay(node Node) Rectangle {
    fillArea := Rectangle {
        x: node.PositionX,
        y: node.PositionY,
        width: node.SizeX,
        height: node.SizeY,
    }

    fullSize := node.Size

 //   results := [][]Rectangle{}
    //active := []Rectangle{}

    vertical := fillArea.height <= fillArea.width 
    cachedSize := 0
    for i := 0; i < len(node.Children); i++ {

        if vertical {
            child := node.Children[0]
            //rect := Rectangle{x: fillArea.x, y: fillArea.y} 
            cachedSize += int(child.Size)
            fmt.Print("cached size: ")
            fmt.Println(cachedSize)
            fraction := float64(cachedSize) / float64(fullSize)
            tmpWidth := fillArea.width * float64(fraction)
            tmpHeight := float64(cachedSize) / float64(tmpWidth)
            rect := Rectangle{width: tmpWidth, height: tmpHeight}
            
            // get the worst ratio there is from the list of rects
            worst := max(rect.width / rect.height, rect.height / rect.width)
            /*for _, r := range active {
                cached := max(r.width / r.height, r.height / r.width)
            }*/
            fmt.Print("w: ")
            fmt.Println(tmpWidth)
            fmt.Print("h: ")
            fmt.Println(tmpHeight)
            fmt.Print("worst ")
            fmt.Println(worst)
        } 
    }

    if vertical {
        child := node.Children[0]
        rect := Rectangle{x: fillArea.x, y: fillArea.y} 

        fraction := child.Size / fullSize
        rect.width = fillArea.width * float64(fraction)
        rect.height = fillArea.height
        return rect
    }
    return Rectangle {} 
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
