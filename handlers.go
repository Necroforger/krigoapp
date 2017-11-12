package krigoapp

import (
	"log"
	"net/http"
	"strconv"
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
		lastDuration  float64
		lastProgress  float64
	)

	setError := func(er error) {
		if er != nil {
			err = er
		}
	}

	// Update loop
	for {
		s.mu.Lock()
		if s.WindowTitle != lastTitle {
			setError(ws.WriteJSON(Event{
				Name: "windowTitle",
				Data: s.WindowTitle,
			}))
			lastTitle = s.WindowTitle
		}
		if s.ThumbnailURL != lastThumbnail {
			setError(ws.WriteJSON(Event{
				Name: "thumbnailURL",
				Data: s.ThumbnailURL,
			}))
			lastThumbnail = s.ThumbnailURL
		}
		if s.VideoURL != lastVideo {
			setError(ws.WriteJSON(Event{
				Name: "videoURL",
				Data: s.VideoURL,
			}))
			lastVideo = s.VideoURL
		}
		if s.Duration != lastDuration {
			setError(ws.WriteJSON(Event{
				Name: "duration",
				Data: strconv.FormatFloat(s.Duration, 'f', -1, 64),
			}))
			lastDuration = s.Duration
		}
		if s.CurrentTime != lastProgress {
			setError(ws.WriteJSON(Event{
				Name: "currentTime",
				Data: strconv.FormatFloat(s.CurrentTime, 'f', -1, 64),
			}))
			lastProgress = s.CurrentTime
		}
		s.mu.Unlock()

		// If there is an error, the websocket was probably dropped.
		if err != nil {
			// log.Println(err)
			return
		}

		select {
		case <-time.After(time.Millisecond * 500):
			continue
		case <-s.close:
			return
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
	if f, ok := r.Form["currentTime"]; ok {
		n, err := strconv.ParseFloat(f[0], 64)
		if err != nil {
			log.Println("UpdateHandler: Error parsing currentTime: ", err)
			return
		}
		s.SetCurrentTime(n)
	}
	if f, ok := r.Form["duration"]; ok {
		n, err := strconv.ParseFloat(f[0], 64)
		if err != nil {
			log.Println("UpdateHandler: Error parsing duration: ", err)
		}
		s.SetDuration(n)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(200)))
}
