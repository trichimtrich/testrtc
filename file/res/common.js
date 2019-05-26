//// Declare in global scope
let wsConn, myID;
let rtcPeer, localSDP, remoteSDP;


//// Utils
log = (msg, isMail) => {
    var pp = document.createElement("p");
    pp.textContent= msg;
    if (isMail != undefined) {
        pp.classList.add("mail");
    }
    document.getElementById("log").append(pp);

    console.log(msg);
}

//// Alter default 'send' function a bit
WebSocket.prototype.send2 = WebSocket.prototype.send;
WebSocket.prototype.send = function (id, data, partnerID) {
    if (id == undefined) {
        log("[!] Must input action ID on websocket send");
        return;
    }
    if (data == undefined) {
        data = "";
    }
    if (partnerID == undefined) {
        partnerID = "";
    }
    this.send2(JSON.stringify({ "id": id, "data": data, "cid": partnerID }));
}

WebSocket.prototype.sendMail = function (partnerID, id, data) {
    if (partnerID == undefined) {
        log("[!] Must input partnerID on mail send");
        return;
    }
    if (id == undefined) {
        log("[!] Must input action ID on mail send");
        return;
    }
    if (data == undefined) {
        data = "";
    }
    this.send("mail", JSON.stringify({ "id": id, "data": data }), partnerID);
}


//// WebSocket stuffs
wsConn = new WebSocket(`ws://${location.host}/ws`);

wsConn.onopen = function () {
    log("WebSocket is opened");
    this.send("hello"); // to get myID

    // create ping interval (every second) to keep connection
    setInterval(() => { this.send("ping"); }, 1000);
}

wsConn.onerror = error => { log(`Websocket error: ${error}`); }
wsConn.onclose = () => { log("Websocket closed"); }

// ws main loop
wsConn.onmessage = function (e) {
    var d = JSON.parse(e.data);
    switch (d["id"]) {
        case "error":
            log(`[!] Error from server: ${d["data"]}`);
            break;

        case "pong":
            break;

        case "hello":
            log(`Got my identity: ${d["data"]}`);
            myID = d["data"];
            break;

        case "mail":
            log(`Got mail from <${d["cid"]}>`);
            try {
                mailObj = JSON.parse(d["data"]);
                handleMail(mailObj, d["cid"]);
            } catch {
                // test data maybe, not in JSON format
            }
            break;

    }
}


//// Common RTC stuffs
function basicRTC() {
    // clean up previous webrtc session
    try {
        rtcPeer.close();
    } catch {
        // ...
    }

    rtcPeer = new RTCPeerConnection({
        iceServers: [
            {
                urls: [
                    'stun:stun.l.google.com:19302',
                    'stun:stun1.l.google.com:19302'
                ],
            },
        ],
    });


    // keep track on state, 'connected' => our GOAL!
    rtcPeer.oniceconnectionstatechange = function(event) {
        log(`iceConnectionState: ${this.iceConnectionState}`);

        if (this.iceConnectionState === "connected") { }
        else if (this.iceConnectionState === "completed") { }
        else if (this.iceConnectionState === "failed") { }
        else if (this.iceConnectionState === "disconnected") { }
    }


    // candidate packet from STUN
    rtcPeer.onicecandidate = function(event) {
        if (event.candidate === null) {
            // ...
        } else {
            log(`Got icecandidate: ${JSON.stringify(event.candidate)}`);
        }
    }
}