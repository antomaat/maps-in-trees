package main

import(
    "os"
    "sort"
    "strings"
    "fmt"
)

func CreateJvmTree(dirName string, pathName string, isJvm bool, rootDir string) Node {
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
            child = CreateJvmTree(items[i].Name(), pathName + "/" + items[i].Name(), isJvm, rootDir)
        }
        if !items[i].IsDir() {
            classFile := findFileInBuildDir(items[i].Name(), rootDir, pathName)
            optionalInfo := OptionalInfo{}
            if classFile != nil {
                fileInfo := ParseFileInfo(classFile)
                optionalInfo.Fields = fileInfo.fields
            }
            child = Node { 
                Name: items[i].Name(),
                PositionX: 0,
                PositionY: 0,
                SizeX: 1,
                SizeY: 1,
                Size: info.Size(), 
                IsDir: false,
                Path: pathName,
                OptionalInfo: optionalInfo,
            }
            fmt.Println("optional info fields ")
            fmt.Println(child.OptionalInfo.Fields)
            fmt.Println("------------------")
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
        //fmt.Println("here is the info: ")
        fmt.Printf("filename %s: \n", fileName)
        //fmt.Printf("rootDir %s: \n", rootDir)
        //fmt.Printf("pathName %s: \n", pathName)
        trimPathName := strings.TrimPrefix(pathName, rootDir + "/src/")
        trimPathName = strings.Replace(trimPathName, "kotlin/", "", 1)
        trimFileName := strings.Replace(fileName, ".kt", "", -1)
        trimFileName = strings.Replace(trimFileName, ".java", "", -1)
        trimFileName += ".class"
        buildDir := rootDir + "/build/classes/kotlin/" + trimPathName
        //buildDir = buildDir + "/" + trimFileName + ".class"
        //fmt.Printf("trim path end: %s \n", trimPathName)
        //fmt.Printf("last path: %s \n", buildDir)
        //fmt.Println("------------------")
        items, _ := os.ReadDir(buildDir)
        for _, n := range items {
            if n.Name() == trimFileName {
                fmt.Println(n.Name())
                fl, _:= os.ReadFile(buildDir + "/" + trimFileName)
                return fl
            }
        }
    }
    return nil
}


