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
        if (confirm("A newer version is available, would you like to update?")) {
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
    sendRPC({action: "EXIT"}); // Prevent malware from hijacking this process, hopefully.
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
    bypassMode.value = localStorage.getItem("bypassMode");
    encryptionOnly.hidden = !tlsMode.checked;
}

if (detectIE() < 11) {
    alert("Your Internet Explorer version is too old, we require IE 11 or above.");
    sendRPC({action: "EXIT"});
}
// taken off https://codebottle.io/s/f17079e6cb
function detectIE() {
    var ua = window.navigator.userAgent;

    var msie = ua.indexOf('MSIE ');
    if (msie > 0) {
        // IE 10 or older => return version number
        return parseInt(ua.substring(msie + 5, ua.indexOf('.', msie)), 10);
    }

    var trident = ua.indexOf('Trident/');
    if (trident > 0) {
        // IE 11 => return version number
        var rv = ua.indexOf('rv:');
        return parseInt(ua.substring(rv + 3, ua.indexOf('.', rv)), 10);
    }

    var edge = ua.indexOf('Edge/');
    if (edge > 0) {
        // Edge (IE 12+) => return version number
        return parseInt(ua.substring(edge + 5, ua.indexOf('.', edge)), 10);
    }
    return false;
}
sendRPC({action: "READY"});