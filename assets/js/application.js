require("expose-loader?$!expose-loader?jQuery!jquery");
require("bootstrap/dist/js/bootstrap.bundle.js");

$(() => {
    var tabId = -1;
    console.log("Extension started")
    const HOST = "127.0.0.1";
    const PORT = "3000";
    const PATH = "/api/v1/websocket/";
    console.log(`Connecting to websocket at: ${HOST}:${PORT}${PATH}`);
    let ws = new WebSocket(`ws://${HOST}:${PORT}${PATH}`);
    
    ws.onopen = (e) => {
        console.log("Connection open!");
        ws.send("Test!")
    }
    
    ws.onclose = (e) => console.log("connection closed!");
});
