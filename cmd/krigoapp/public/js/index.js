window.addEventListener("load", function () {
    // Connect websocket.
    var ws       = new ReconnectingWebSocket("ws://" + window.location.host + "/ws/");
    var duration = 0; 

    ws.onmessage = function (e) {
        var data = JSON.parse(e.data);
        console.log(data);

        switch (data.name) {
            case "windowTitle":
                updateWindowTitles(data.content);
                break;
            case "videoURL":
                updateVideoURLs(data.content);
                break;
            case "thumbnailURL":
                updateThumbnails(data.content);
                break;
            case "currentTime":
                updateProgress(parseFloat(data.content), duration);
                updateCurrentTime(parseInt(data.content));
                break;
            case "duration":
                duration = parseFloat(data.content);
                updateDuration(parseInt(data.content));
                break;
            default:
                console.log("Event: " + event.name + " not supported");
        }
    };
});


function updateThumbnails(thumbnailURL) {
    console.log("updating video thumbnails: " + thumbnailURL);

    var elems = document.getElementsByClassName("video-thumbnail");
    for (var i=0; i < elems.length; i++) {
        elems[i].src = thumbnailURL;
        reAnimate(elems[i]);
    }
}

function updateWindowTitles(str) {
    str = stripWindowSuffixes(str);

    console.log("updating window titles: " + str);
    var elems = document.getElementsByClassName("window-title");
    for (var i = 0; i < elems.length; i++) {
        elems[i].innerHTML = str;
        reAnimate(elems[i]);
    }
}

function updateVideoURLs(str) {
    console.log("updating video URLs: " + str);
    var elems = document.getElementsByClassName("video-url");
    for (var i = 0; i < elems.length; i++) {
        elems[i].innerHTML = str;
        reAnimate(elems[i]);
    }
}

function updateProgress(current, duration) {
    var percent = Math.floor((current / duration)*100);
    console.log("percent: " + percent);

    var elems = document.getElementsByClassName("video-progress");
    for (var i=0; i < elems.length; i++) {
        elems[i].style.backgroundImage = "linear-gradient(to right, rgba(255, 255, 255, 0.1) "+percent+"%, transparent "+percent+"%)";
    }
}

function updateCurrentTime(t) {
    console.log("updating current time: " + t);
    var elems = document.getElementsByClassName("video-current-time");
    for (var i=0; i < elems.length; i++) {
        elems[i].innerHTML = t;
    }
}

function updateDuration(t) {
    console.log("updating current time: " + t);
    var elems = document.getElementsByClassName("video-duration");
    for (var i=0; i < elems.length; i++) {
        elems[i].innerHTML = t;
    }
}

function stripWindowSuffixes(str) {
    suffixes = [
        " - YouTube - YouTube - Google Chrome",
        " - YouTube - Google Chrome",
        " - Google Chrome",
        " - Mozilla Firefox",
    ];

    var cutAmount = 0;
    for (var i = 0; i < suffixes.length; i++) {
        if (str.endsWith(suffixes[i])) {
            cutAmount = Math.max(cutAmount, suffixes[i].length);
        }
    }
    return str.substring(0, str.length - cutAmount);
}

function reAnimate(elem) {
    elem.parentNode.classList.remove("animate");
    setTimeout(function () {
        elem.parentNode.classList.add("animate");
    }, 50);
}
