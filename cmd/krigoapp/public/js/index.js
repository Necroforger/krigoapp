window.addEventListener("load", function () {
    // Connect websocket.
    var ws = new ReconnectingWebSocket("ws://" + window.location.host + "/ws/");

    ws.onmessage = function (e) {
        var data = JSON.parse(e.data);
        console.log(data);
        if (data.windowTitle) {
            updateWindowTitles(data.windowTitle);
        }

    };
});

function updateWindowTitles(str) {  
    str = stripWindowSuffixes(str);

    console.log("updating window titles: " + str);
    var elems = document.getElementsByClassName("window-title");
    for( var i=0; i < elems.length; i++ ) {
        elems[i].innerHTML = str;
        reAnimate(elems[i]);
    }
}

function stripWindowSuffixes(str) {
    suffixes = [
        " - YouTube - YouTube - Google Chrome",
        " - YouTube - Google Chrome",
        " - Google Chrome",
    ];

    var cutAmount = 0;
    for (var i=0; i < suffixes.length; i++) {
        if (str.endsWith(suffixes[i])) {
            cutAmount = Math.max(cutAmount, suffixes[i].length);
        }
    }
    return str.substring(0, str.length-cutAmount);
}

function reAnimate(elem) {
    elem.parentNode.classList.remove("animate");
    setTimeout(function() {
        elem.parentNode.classList.add("animate");
    }, 50);
}
