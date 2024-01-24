
let treemap;
let canvas;
let shouldRedraw = true;
let items = [];

let canvasRect;

let mouse_pos;


const path = document.getElementById("path");

const context = createCanvas();

loadJson().then(tree => {tree;});

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
                path: node.Path
            }
        );
    }
}

function initTree(tree) {
    console.log("tree: ", tree);
    if (tree.Children === null) {
        return
    }
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
    for (node of items) {
        drawBorder(context, node.positionX, node.positionY, node.sizeX, node.sizeY);
        context.fillStyle = "#000000";
        context.fillRect(node.positionX, node.positionY, node.sizeX, node.sizeY);
    }
}

function drawBorder(ctx, xPos, yPos, width, height, thickness = 2)
{
  ctx.fillStyle='#fff';
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

