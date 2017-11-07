package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Necroforger/krigoapp"
	"github.com/Necroforger/krigoapp/window"
)

// vars
var (
	mu            sync.Mutex
	TrackedWindow *window.HWND

	PublicFolder = flag.String("f", "./public", "public folder")
	Address      = flag.String("a", "127.0.0.1:7777", "address the server runs on")
)

func handleError(err error) {
	if err != nil {
		log.Println(err)
		panic(err)
	}
}

func printMenu() {
	fmt.Println("================================")
	fmt.Println("| 1. Select window              |")
	fmt.Println("| 2. Select window from list    |")
	fmt.Println("| 9. Exit                       |")
	fmt.Println("================================")
}

func shell(s *krigoapp.Server) {
	for {
		printMenu()
		fmt.Print("> ")

		line, err := readline()
		handleError(err)

		execCommand(s, line)
	}
}

func readline() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(line), nil
}

func execCommand(s *krigoapp.Server, line string) {
	switch {
	case line == "1":
		setWindow(s)
	case line == "2":
		setWindowFromList(s)
	case line == "9":
		fmt.Println("Exiting...")
		s.Close()
		os.Exit(0)
	default:
		fmt.Println("Invalid command")
	}
}

func setWindow(s *krigoapp.Server) {
	fmt.Println("Please open the window you wish to capture the title of")
	for i := 3; i > 0; i-- {
		fmt.Println("Selecting window in ", i, "seconds")
		time.Sleep(time.Second * 1)
	}

	mu.Lock()
	TrackedWindow = window.GetForegroundWindow()
	mu.Unlock()

	fmt.Println("Capturing window [", TrackedWindow.Title(), "]")
}

func setWindowFromList(s *krigoapp.Server) {
	windows := window.GetAllWindows()
	for i := len(windows) - 1; i >= 0; i-- {
		w := windows[i]
		if title := w.Title(); strings.TrimSpace(title) != "" {
			fmt.Printf("[%d] : %s\n", i, title)
		}
	}

	fmt.Println("Please select the desired window from the list.")
	fmt.Println("Enter a noninteger or number less than 0 to cancel")
	fmt.Printf("> ")
	line, err := readline()
	handleError(err)

	n, err := strconv.Atoi(strings.TrimSpace(line))
	if err != nil || n < 0 || n >= len(windows) {
		fmt.Println("Cancelled..")
		return
	}

	mu.Lock()
	TrackedWindow = windows[n]
	mu.Unlock()

	fmt.Println("Capturing window [", TrackedWindow.Title(), "]")
}

// updateTitle continuously updates the window title
func updateTitle(s *krigoapp.Server) {
	ticker := time.NewTicker(time.Second * 1)
	for _ = range ticker.C {
		mu.Lock()
		if TrackedWindow != nil {
			s.SetWindowTitle(TrackedWindow.Title())
		}
		mu.Unlock()
	}
}

func main() {
	server := krigoapp.NewServer(*PublicFolder, *Address)
	fmt.Printf("Running server on [%s]\n", *Address)
	go shell(server)
	go updateTitle(server)
	err := server.Start()
	if err != nil {
		fmt.Println(err)
	}
}
