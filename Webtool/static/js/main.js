function init() {
    openWS();
}

function openWS() {
    let table = document.getElementById("logContainer");

    const ws = new WebSocket("ws://localhost:8080/ws");
    ws.onmessage = (event) => {
        let data = JSON.parse(event.data);

        let row = table.insertRow(-1);
        row.insertCell(0).innerHTML = data.timestamp;
        row.insertCell(1).innerHTML = data.type;
        row.insertCell(2).innerHTML = data.message;
        row.scrollIntoView();
    };

    ws.onclose = () => {
        let row = table.insertRow(-1);
        row.insertCell(0).innerHTML = "-";
        row.insertCell(1).innerHTML = "ERROR";
        row.insertCell(2).innerHTML = "<span style=\"color: red;\">Connection to server lost!</span>";
        row.scrollIntoView();
    }

}

function launchPing() {
    createAPICall("GET", "/api/ping", (status, response) => {
        if (status !== STATUS_OK) alert("Ping failed!");
    });
}

window.onload = init;