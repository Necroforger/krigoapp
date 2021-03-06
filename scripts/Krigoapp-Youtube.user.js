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


    var ytplayer,
        videoData,
        videoURL,
        videoThumbnail;

    setInterval(function () {
        ytplayer       = document.querySelector("div[id^='player_']") || document.getElementById("movie_player");
        videoData      = ytplayer.getVideoData();
        videoURL       = "youtu.be/"+ videoData.video_id;
        // // Include the playlist ID in the link
        // if (ytplayer.getPlaylistId() !== "") videoURL += "?list="+ytplayer.getPlaylistId();
        videoThumbnail = "https://i.ytimg.com/vi/" + videoData.video_id + "/hqdefault.jpg";

        GM_xmlhttpRequest({
            method: "GET",
            url: "http://127.0.0.1:7777/update"  +
            "?windowTitle="  + encodeURIComponent(videoData.title) +
            "&thumbnailURL=" + encodeURIComponent(videoThumbnail)  +
            "&videoURL="     + encodeURIComponent(videoURL) +
            "&currentTime="  + encodeURIComponent(ytplayer.getCurrentTime()) +
            "&duration="     + encodeURIComponent(ytplayer.getDuration())
        });

    }, 1000);
})();
