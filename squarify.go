package main

import(
    "math"
)

func SquarifyDisplay(node Node) Node {

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
            node.Children[indx] = SquarifyDisplay(node.Children[indx])
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
