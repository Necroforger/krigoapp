package krigoapp

import (
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"

	"github.com/gorilla/websocket"
)

// Update is the update information sent over websocket
type Update struct {
	WindowTitle string `json:"windowTitle"`
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

// wsHandler handles websocket connections
func (s *Server) wsHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	var lastTitle string

	for {
		// Avoid sending duplicate data
		s.mu.Lock()
		if lastTitle == s.WindowTitle {
			s.mu.Unlock()
			continue
		}
		lastTitle = s.WindowTitle
		data := Update{
			WindowTitle: s.WindowTitle,
		}
		s.mu.Unlock()

		err = ws.WriteJSON(data)
		if err != nil {
			return
		}

		select {
		case <-time.After(time.Millisecond * 500):
		case <-s.close:
		}
	}
}
