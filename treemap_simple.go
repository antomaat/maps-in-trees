package main

import(
    "os"
    "sort"
)

func CreateTree(dirName string, pathName string, isJvm bool) Node {
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
            child = CreateTree(items[i].Name(), pathName + "/" + items[i].Name(), isJvm)
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

