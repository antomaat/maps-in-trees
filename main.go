package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"slices"
	"sort"
	"strings"
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
    //node := createTree(dir, dir, isJvm)
    node := createJvmTree(dir, dir, isJvm, dir)
    node.SizeX = areaX
    node.SizeY = areaY
    debug, _ := json.Marshal(node)
    _ = ioutil.WriteFile("debug.json", debug, 0644)
    //fmt.Println(node)
    //node = updateDisplay(node, 0)
    node = NewSquarifyDisplay(node)
    result, error := json.Marshal(node)
    if error != nil {
        fmt.Println(error.Error())
    }

    _ = ioutil.WriteFile(output, result, 0644)
}

func createJvmTree(dirName string, pathName string, isJvm bool, rootDir string) Node {
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
        if items[i].IsDir() && !isIgnoredDir(items[i].Name()) {
            child = createJvmTree(items[i].Name(), pathName + "/" + items[i].Name(), isJvm, rootDir)
        } else {
            classFile := findFileInBuildDir(items[i].Name(), rootDir, pathName)
            ParseFileInfo(classFile)
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

        // add only if the file is bigger than 0
        if (child.Size > 0) {
            node.Children = append(node.Children, child)
            //node.SizeX += child.SizeX
            node.Size += child.Size;
        }
    }

    sort.SliceStable(node.Children, func(i, j int) bool {
        return node.Children[i].Size > node.Children[j].Size
    })
   
    return node
}

func findFileInBuildDir(fileName string, rootDir string, pathName string) []byte  {
    //items, _ := os.ReadDir(rootDir + "/build/classes")
    if (strings.Contains(fileName, ".kt")) {
        fmt.Println("here is the info: \n")
        fmt.Printf("filename %s: \n", fileName)
        fmt.Printf("rootDir %s: \n", rootDir)
        fmt.Printf("pathName %s: \n", pathName)
        trimPathName := strings.TrimPrefix(pathName, rootDir + "/src/")
        trimPathName = strings.Replace(trimPathName, "kotlin/", "", 1)
        trimFileName := strings.Replace(fileName, ".kt", "", -1)
        trimFileName = strings.Replace(trimFileName, ".java", "", -1)
        trimFileName += ".class"
        buildDir := rootDir + "/build/classes/kotlin/" + trimPathName
        //buildDir = buildDir + "/" + trimFileName + ".class"
        fmt.Printf("trim path end: %s \n", trimPathName)
        fmt.Printf("last path: %s \n", buildDir)
        fmt.Println("------------------")
        items, _ := os.ReadDir(buildDir)
        fmt.Println(items)
        for _, n := range items {
            if n.Name() == trimFileName {
                fmt.Println("found the file")
                fmt.Println(n.Name())
                fl, _:= os.ReadFile(buildDir + "/" + trimFileName)
                return fl
            }
        }
    }
    return nil
}

func createTree(dirName string, pathName string, isJvm bool) Node {
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
        if items[i].IsDir() && !isIgnoredDir(items[i].Name()) {
            child = createTree(items[i].Name(), pathName + "/" + items[i].Name(), isJvm)
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

        // add only if the file is bigger than 0
        if (child.Size > 0) {
            node.Children = append(node.Children, child)
            //node.SizeX += child.SizeX
            node.Size += child.Size;
        }
    }

    sort.SliceStable(node.Children, func(i, j int) bool {
        return node.Children[i].Size > node.Children[j].Size
    })
   
    return node
}


func NewSquarifyDisplay(node Node) Node {

    fillArea := Rectangle {
        x: node.PositionX,
        y: node.PositionY,
        width: node.SizeX,
        height: node.SizeY,
        size: float64(node.Size),
    }
    directories := []int{}
    cache := [][]Node{}
    row := []Node{}
    for _, n := range node.Children {
        vertical := fillArea.height <= fillArea.width 
        widthN := fillArea.width
        heightN := fillArea.height
        if vertical {
            widthN = fillArea.height
            heightN = fillArea.width
        }

        rowWithChild := append(row, n)

        if len(row) == 0 || worst(row, widthN, heightN, float64(fillArea.size)) >= worst(rowWithChild, widthN, heightN, float64(fillArea.size)) {
            row = append(row, n)
        } else {
            row = layoutRow(row, widthN, vertical, &fillArea)
            cache = append(cache, row)
            row = []Node{}
            row = append(row, n)
        }
    }

    // maybe a single node is left unprocessed
    if len(row) > 0 {
        vertical := fillArea.height <= fillArea.width 
        widthN := fillArea.width
        if vertical {
            widthN = fillArea.height
        }
        row = layoutRow(row, widthN, vertical, &fillArea)
        cache = append(cache, row)
        row = []Node{}
    }
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

func layoutRow(row []Node, smallestSide float64, vertical bool, parent *Rectangle) []Node {
    result := []Node{}
    cacheParent := Rectangle {
        x: parent.x,
        y: parent.y,
        width: parent.width,
        height: parent.height,
        size: parent.size,
    }


    sizes := sumSizes(row)
    fraction :=  float64(sizes) / float64(cacheParent.size)

    longestSide := parent.width
    if !vertical {
        longestSide = parent.height
    }

    area := longestSide * float64(fraction) 

    for _, node := range row {
        // step 1 - calculate the size of the node
        nodeOtherSide := float64(node.Size) / sizes

        if (vertical) {
            node.SizeY = smallestSide * nodeOtherSide
            node.SizeX = area 
            node.PositionY = cacheParent.y
            node.PositionX = cacheParent.x

            cacheParent.y += node.SizeY

        } else {
            node.SizeX = smallestSide * nodeOtherSide
            node.SizeY = area 
            node.PositionX = cacheParent.x
            node.PositionY = cacheParent.y

            cacheParent.x += node.SizeX

        }

        result = append(result, node)
    }

    if vertical {
        (*parent).x += area
        (*parent).width -= area
    } else {
        (*parent).y += area
        (*parent).height -= area
    }
    (*parent).size -= sizes
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

func worst(sizes []Node, w float64, h float64, parentSize float64) float64 {
    max := math.Inf(-1)
    min := math.Inf(1)
    sum := 0.0
    for _, size := range sizes {
        sum += float64(size.Size)
        max = math.Max(max, float64(size.Size))
        min = math.Min(min, float64(size.Size))
    }

    ratio := math.Max(calculateRatio(max, sum, w, h, parentSize), calculateRatio(min, sum, w, h, parentSize))
    return ratio
}

func calculateRatio(value float64, sum float64, w float64, h float64, parentSize float64) float64 {
    fraction := value / sum
    width := w * fraction
    otherFraction := sum / parentSize
    height := h * otherFraction 
    ratio := math.Max(height / width, width / height)
    return ratio
}

func sumSizes(nodes []Node) float64 {
    sum := 0.0
    for _, n := range nodes {
        sum += float64(n.Size)
    }
    return sum
}

