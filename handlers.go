package krigoapp

import (
	"log"
	"net/http"
	"time"
)

// wsHandler handles websocket connections
func (s *Server) wsHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	var (
		lastTitle     string
		lastThumbnail string
		lastVideo     string
	)

	for {
		s.mu.Lock()
		if s.WindowTitle != lastTitle {
			ws.WriteJSON(Event{
				Name: "windowTitle",
				Data: s.WindowTitle,
			})
			lastTitle = s.WindowTitle
		}
		if s.ThumbnailURL != lastThumbnail {
			ws.WriteJSON(Event{
				Name: "thumbnailURL",
				Data: s.ThumbnailURL,
			})
			lastThumbnail = s.ThumbnailURL
		}
		if s.VideoURL != lastVideo {
			ws.WriteJSON(Event{
				Name: "videoURL",
				Data: s.VideoURL,
			})
			lastVideo = s.VideoURL
		}
		s.mu.Unlock()

		select {
		case <-time.After(time.Millisecond * 500):
		case <-s.close:
		}
	}
}

// UpdateHandler ...
func (s *Server) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Println("UpdateHandler: Error parsing form: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
		return
	}

	if f, ok := r.Form["windowTitle"]; ok {
		s.SetWindowTitle(f[0])
	}
	if f, ok := r.Form["thumbnailURL"]; ok {
		s.SetThumbnailURL(f[0])
	}
	if f, ok := r.Form["videoURL"]; ok {
		s.SetVideoURL(f[0])
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(200)))
}
