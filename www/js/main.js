var serverAddr = document.getElementById("serverAddress");
var autoWL = document.getElementById("pacType");
var tlsMode = document.getElementById("tls");
var validateCert = document.getElementById("validateCert");
var startButton = document.getElementById("startButton");
var encryptionType = document.getElementById("encryptionType");
var encryptionOnly = document.getElementById("encryptionOnly");
var bypassMode = document.getElementById("bypassMode");
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
window.sendLine = function(s) {
    window.open("https://www.w3schools.com");
    alert(s);
};


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
        localStorage.setItem("encryptionType", encryptionType.value);
        localStorage.setItem("bypassMode", bypassMode.value);
        var addr;
        if (tlsMode.checked) {
            addr = "wss://" + serverAddr.value + "/socksproxy";
        } else {
            addr = "ws://" + serverAddr.value + "/socksproxy";
        }
        sendRPC({
            action: "CONNECT",
            server: addr,
            pac: autoWL.value,
            validateCertificate: validateCert.checked.toString(),
            encryptionType: encryptionType.value,
            bypassType: bypassMode.value
        })
    }
};

tlsMode.onchange = function() {
    encryptionOnly.hidden = !tlsMode.checked;
};

if (localStorage.getItem("server") !== null) {
    serverAddr.value = localStorage.getItem("server");
    autoWL.value = localStorage.getItem("WL");
    tlsMode.checked = localStorage.getItem("TLS") === "true";
    encryptionType.value = localStorage.getItem("encryptionType");
    validateCert.checked = localStorage.getItem("validateCert") === "true";
    bypassMode.value = localStorage.getItem("bypassMode")
    encryptionOnly.hidden = !tlsMode.checked;
}


sendRPC({action: "READY"});