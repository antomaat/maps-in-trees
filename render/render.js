
const treemap = loadJson();
let canvas
const context = createCanvas();

function createCanvas() {
    canvas = document.getElementById("myCanvas");
    const ctx = canvas.getContext("2d");

    ctx.fillStyle = "#FF0000";
    ctx.fillRect(0, 0, 1000, 500);

    return ctx;
}


function drawTreemap() {

}


function loadJson() {
    return fetch('./input.json')
    .then((response) => response.json())
    .then((json) => console.log(json));
}
