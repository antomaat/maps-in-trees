
let treemap;
let canvas;
let shouldRedraw = true;
let items = [];

let canvasRect;

let mouse_pos;


const path = document.getElementById("path");

const context = createCanvas();

let selectedNode = undefined

loadJson().then(tree => {tree;});

function onSearch() {
    const searchTerm = document.getElementById("search_file").value.trim()
    for (node of items) {
        if (node.name === undefined) continue;
        if (node.name.includes(searchTerm)) {
            selectedNode = node;
            drawSelectedNode(node);
        }
    }
}

function drawSelectedNode(node) {
    context.fillStyle = "#F5F5DC";
    context.fillRect(node.positionX, node.positionY, node.sizeX, node.sizeY);
}

function createCanvas() {
    canvas = document.getElementById("myCanvas");
    canvas.addEventListener("mousemove", onMouseMove) 
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

function run() {
    while(shouldRedraw) {
        console.log("redraw");
        ctx.clearRect(0, 0, canvas.width, canvas.height);
        ctx.fillStyle = "#FF0000";
        ctx.fillRect(0, 0, 2000, 500);
        getItemAndDisplay();
        drawTreemap(ctx);
    }
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

function initItems(tree) {
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
                path: node.Path
            }
        );
        if (node.IsDir) {
            initTree(node)
        }
    }
}

function drawTreemap(context) {
    console.log(items);
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
    //initItems(treemap);
    initTree(treemap)
    drawTreemap(context);
}

