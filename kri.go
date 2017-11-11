package krigoapp

import (
	"net/http"
	"sync"

	"github.com/gorilla/mux"

	"github.com/gorilla/websocket"
)

// Event is the update information sent over websocket
type Event struct {
	Name string `json:"name"`
	Data string `json:"content"`
}

// Server ...
type Server struct {
	mu       sync.Mutex
	upgrader websocket.Upgrader
	close    chan int

	// Server underlying http server
	Server *http.Server

	// WindowTitle current title of the selected window
	WindowTitle string

	// Thumbnail is a URL to the thumbnail of the currently playing song
	ThumbnailURL string

	// VideoURL is the URL of the current song
	VideoURL string

	// Directory the server will serve files from
	Servedir string
}

// NewServer creates a new server with the default settings
func NewServer(servedir string, addr string) *Server {
	s := &Server{
		Servedir: servedir,
	}

	r := mux.NewRouter()
	r.HandleFunc("/ws/", s.wsHandler)
	r.HandleFunc("/update", s.UpdateHandler)
	r.PathPrefix("/").Handler(http.FileServer(http.Dir(servedir)))

	s.Server = &http.Server{
		Addr:    addr,
		Handler: r,
	}
	return s
}

// SetWindowTitle sets the currently tracked window title
func (s *Server) SetWindowTitle(title string) {
	s.mu.Lock()
	s.WindowTitle = title
	s.mu.Unlock()
}

// SetThumbnailURL sets the thumbnail URL of the server.
func (s *Server) SetThumbnailURL(URL string) {
	s.mu.Lock()
	s.ThumbnailURL = URL
	s.mu.Unlock()
}

// SetVideoURL sets the video URL
func (s *Server) SetVideoURL(URL string) {
	s.mu.Lock()
	s.VideoURL = URL
	s.mu.Unlock()
}

// Start begins listening for connections and hotkeys
func (s *Server) Start() error {
	s.close = make(chan int)
	return s.Server.ListenAndServe()
}

// Close closes the running server
func (s *Server) Close() error {
	if s.close != nil {
		close(s.close)
	}
	return s.Server.Close()
}
