
let treemap;
let canvas;
let shouldRedraw = true;
let items = [];

let canvasRect;

let mouse_pos;


const path = document.getElementById("path");

const context = createCanvas();

let selected = []
let selectedChilds = []
let selectedParents = []

loadJson().then(tree => {tree;});

function onAdvancedSearch() {
    const searchTerm = document.getElementById('advanced_search').value.trim()
    selected = [];
    redrawTreemap();
    for (node of items) {
        if (node.optional !== undefined) {
            if (node.optional.Fields !== null) {
                for (const name of node.optional.Fields) {
                    if (name.includes(searchTerm)) {
                        selected.push(node);
                        drawSelectedNode(node);
                    }
                }
            }
        } 
    }
}

function onSearch() {
    const searchTerm = document.getElementById("search_file").value.trim()
    selected = [];
    redrawTreemap();
    for (node of items) {
        if (node.name === undefined) continue;
        if (node.name.includes(searchTerm)) {
            selected.push(node);
            drawSelectedNode(node);
        }
    }
}

function drawSelectedNode(node) {
    context.fillStyle = "#F5F5DC";
    context.fillRect(node.positionX, node.positionY, node.sizeX, node.sizeY);
}

function drawSelectedChildNode(node) {
    context.fillStyle = "#FFC133";
    context.fillRect(node.positionX, node.positionY, node.sizeX, node.sizeY);
}

function drawSelectedParentNode(node) {
    context.fillStyle = "#FF4F33";
    context.fillRect(node.positionX, node.positionY, node.sizeX, node.sizeY);
}

function createCanvas() {
    canvas = document.getElementById("myCanvas");
    canvas.addEventListener("mousemove", onMouseMove) 
    canvas.addEventListener("click", onClickCanvas) 
    canvasRect = canvas.getBoundingClientRect();

    const ctx = canvas.getContext("2d");

    ctx.fillStyle = "#FF0000";
    ctx.fillRect(0, 0, 2000, 500);

    return ctx;
}

function onMouseMove(event) {
    mouse_pos = {
        x: event.offsetX,
        y: event.offsetY
    };
    getItemAndDisplay();
}

function onClickCanvas(event) {
    selected = [];
    selectedChilds = [];
    selectedParents = [];
    const node = getItemFromMouse();
    console.log(node);
    if (node != null) {
        redrawTreemap();
        selected.push(node);

        if (node.optional.Fields !== null) {
            for (const name of node.optional.FieldClasses) {
                let nameSplit = name.split('/')
                nameSplit = nameSplit[nameSplit.length - 1].replace(';', '');
                for (const childNode of items) {
                    if (childNode.name === undefined) continue;
                    if (childNode.name.toLowerCase().includes(nameSplit.toLowerCase())) {
                        selectedChilds.push(childNode);
                        break;
                    }
                }
            }
        }
        for (const n of items) {
            if (n.name === undefined) continue;
            if (n.optional !== null && n.optional.FieldClasses !== null) {
                for (const nField of n.optional.FieldClasses) {
                    let nameSplit = nField.split('/')
                    nameSplit = nameSplit[nameSplit.length - 1].replace(';', '');
                    if (node.name.toLowerCase().includes(nameSplit.toLowerCase())) {
                        console.log("name split");
                        console.log(nameSplit);
                        selectedParents.push(n)
                    }
                    
                }
            }
        }


        for (const n of selected) {
            drawSelectedNode(n);
        }
        for (const n of selectedChilds) {
            drawSelectedChildNode(n);
        }

        for (const n of selectedParents) {
            drawSelectedParentNode(n);
        }
    }
}

function redrawTreemap() {
    context.clearRect(0, 0, canvas.width, canvas.height);
    context.fillStyle = "#FF0000";
    context.fillRect(0, 0, 2000, 500);
    drawTreemap(context);
}

function getItemFromMouse() {
    for (node of items) {
        if (!node.isDir) {
            if (mouse_pos.x > node.positionX && mouse_pos.x < node.positionX + node.sizeX) {
                if (mouse_pos.y > node.positionY && mouse_pos.y < node.positionY + node.sizeY) {
                    return node;
                }
            }

        }
    }
    return null;
}

function getItemAndDisplay() {
    for (node of items) {
        if (mouse_pos.x > node.positionX && mouse_pos.x < node.positionX + node.sizeX) {
            if (mouse_pos.y > node.positionY && mouse_pos.y < node.positionY + node.sizeY) {
                path.innerHTML = node.path + "/" + node.name;
            }
        }
    }
}

function initTree(tree) {
    console.log("tree: ", tree);
    if (tree.Children === null) {
        return
    }
    items.push(tree);
    for (node of tree.Children) {
        items.push(
            {
                name: node.Name,
                positionX: node.PositionX,
                positionY: node.PositionY,
                sizeX: node.SizeX,
                sizeY: node.SizeY,
                children: node.Children,
                isDir: node.IsDir,
                path: node.Path,
                optional: node.OptionalInfo,
            }
        );
        if (node.IsDir) {
            initTree(node)
        }
    }
}

function drawTreemap(context) {
    dirs = [];
    for (node of items) {
        if (node.isDir) {
            dirs.push(node);
            context.fillStyle = "#808080";
        } else {
            drawBorder(context, node.positionX, node.positionY, node.sizeX, node.sizeY);
            context.fillStyle = "#000000";
        }
        context.fillRect(node.positionX, node.positionY, node.sizeX, node.sizeY);
    }

    for (dir of dirs) {
        drawDirBorder(context, dir.positionX, dir.positionY, dir.sizeX, dir.sizeY, '#8A2BE2', 4);
    }

}

function drawDirBorder(ctx, xPos, yPos, width, height, color = '#fff', thickness = 2)
{
    ctx.lineWidth = thickness;
    ctx.strokeStyle = color;
    ctx.beginPath();
    ctx.roundRect(xPos, yPos, width, height, 0);
    ctx.stroke();
}


function drawBorder(ctx, xPos, yPos, width, height, color = '#fff', thickness = 2)
{
  ctx.fillStyle=color;
  ctx.fillRect(xPos - (thickness), yPos - (thickness), width + (thickness * 2), height + (thickness * 2));
}


async function loadJson() {
    const response = await fetch('./input.json');
    const treeJson = await response.json();
    treemap = treeJson;
    initTree(treemap)
    drawTreemap(context);
}

