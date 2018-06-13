var serverAddr = document.getElementById("serverAddress");
var autoWL = document.getElementById("pacType");
var tlsMode = document.getElementById("tls");
var validateCert = document.getElementById("validateCert");
var startButton = document.getElementById("startButton");
var connected = false;

window.sendRPC = function (data) {
    window.external.invoke(JSON.stringify(data))
};
window.receiveRPC = function (data) {
    if (data.cmd === "setConnectionStatus") {
        connected = data.status;
        if (data.status) {
            startButton.style.backgroundColor = "green";
            startButton.innerHTML = "Connected";
        } else {
            startButton.style.backgroundColor = "dodgerblue";
            startButton.innerHTML = "Connect";
        }
    } else if (data.cmd === "showUpdateScreen") {
        document.body.innerHTML = '<h1 id="title" class="center-all">Updating...</h1>';
    } else if (data.cmd === "showUpdatePrompt") {
        if (confirm("A newer version is available, would you like to update?") === true) {
            window.sendRPC({action: "UPDATE"})
        }
    }
};

if (localStorage.getItem("server") !== null) {
    serverAddr.value = localStorage.getItem("server");
    autoWL.value = localStorage.getItem("WL");
    tlsMode.checked = localStorage.getItem("TLS") === "true";
    validateCert.checked = localStorage.getItem("validateCert") === "true";
}

window.onbeforeunload = function () {
    alert("It appears that some kind of malware has hijacked this page, the program will now exit.");
    sendRPC({action: "PAGECHANGE"}); // Prevent malware from hijacking this process, hopefully.
    return "Don't leave!"
};

startButton.onclick = function () {
    if (connected) {
        startButton.innerHTML = "Disconnecting";
        startButton.style.backgroundColor = "yellow";
        sendRPC({action: "DISCONNECT"})
    } else {
        startButton.innerHTML = "Connecting";
        startButton.style.backgroundColor = "yellow";
        localStorage.setItem("server", serverAddr.value);
        localStorage.setItem("WL", autoWL.value);
        localStorage.setItem("TLS", tlsMode.checked.toString());
        localStorage.setItem("validateCert", validateCert.checked.toString());
        var addr;
        if (tlsMode.checked) {
            addr = "wss://" + serverAddr.value + "/";
        } else {
            addr = "ws://" + serverAddr.value + "/";
        }
        sendRPC({
            action: "CONNECT",
            server: addr,
            pac: autoWL.value,
            validateCertificate: validateCert.checked.toString()
        })
    }
};
sendRPC({action: "READY"});