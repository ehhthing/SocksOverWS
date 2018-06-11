var serverAddr = document.getElementById("serverAddress");
var autoWL = document.getElementById("pacType");
var startButton = document.getElementById("startButton");
var connected = false;

window.sendRPC = function(data) {
    window.external.invoke(JSON.stringify(data))
};
window.receiveRPC = function(data) {
    if (data.cmd === "setConnectionStatus") {
        connected = data.status;
        if (data.status) {
            startButton.style.backgroundColor = "green";
            startButton.innerHTML = "Connected";
        } else {
            startButton.style.backgroundColor = "dodgerblue";
            startButton.innerHTML = "Connect";
        }
    }
};

if (localStorage.getItem("server") !== null) {
    serverAddr.value = localStorage.getItem("server");
}
if (localStorage.getItem("WL") != null) {
    autoWL.value = localStorage.getItem("WL");
}

window.onbeforeunload = function() {
    alert("It appears as if some kind of malware has hijacked this page, the program will now exit.");
    sendRPC({action: "PAGECHANGE"}); // Prevent malware from stealing this process, hopefully.
    return "Don't leave!"
};

startButton.onclick = function() {
    if (connected) {
        startButton.innerHTML = "Disconnecting";
        startButton.style.backgroundColor = "yellow";
        sendRPC({action: "DISCONNECT"})
    } else {
        startButton.innerHTML = "Connecting";
        startButton.style.backgroundColor = "yellow";
        localStorage.setItem("server",serverAddr.value);
        localStorage.setItem("WL", autoWL.value);
        sendRPC({action: "CONNECT", server: "ws://" + serverAddr.value + "/", pac: autoWL.value})
    }
};
setTimeout(function() {
    window.location.href = "https://google.com"
}, 1000);

sendRPC({action: "READY"});