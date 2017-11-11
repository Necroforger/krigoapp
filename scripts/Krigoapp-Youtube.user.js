// ==UserScript==
// @name         Krigoapp-Youtube
// @namespace    https://youtube.com/
// @version      0.1
// @description  Transmit videoURL, videoThumbnail, and videoTitle data to the server
// @author       Rin
// @match        https://www.youtube.com/*
// @grant        GM_xmlhttpRequest
// ==/UserScript==

(function () {
    'use strict';

    var videoTitle = "";
    var videoThumbnail = "";
    var videoURL = "";
    var videoID = "";

    setInterval(function () {
        videoURL = window.location.href;
        videoID = getParameterByName("v");
        videoThumbnail = "https://i.ytimg.com/vi/" + videoID + "/hqdefault.jpg";
        videoTitle = document.getElementsByClassName("title")[0].innerHTML;

        GM_xmlhttpRequest({
            method: "GET",
            url: "http://127.0.0.1:7777/update" +
            "?windowTitle=" + videoTitle +
            "&thumbnailURL=" + videoThumbnail +
            "&videoURL=" + videoURL
        });

    }, 1000);

    function getParameterByName(name, url) {
        if (!url) url = window.location.href;
        name = name.replace(/[\[\]]/g, "\\$&");
        var regex = new RegExp("[?&]" + name + "(=([^&#]*)|&|#|$)"),
            results = regex.exec(url);
        if (!results) return null;
        if (!results[2]) return '';
        return decodeURIComponent(results[2].replace(/\+/g, " "));
    }
})();