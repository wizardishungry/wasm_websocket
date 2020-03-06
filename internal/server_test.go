package internal

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/gorilla/handlers"
	"github.com/gorilla/websocket"
)

var (
	server http.Server
	errors chan error
)

func TestMustGlobal(t *testing.T) {
	addr := GetServerAddr()
	errors = make(chan error)

	mux := http.NewServeMux()
	mux.HandleFunc("/ws", ws)
	mux.HandleFunc("/quit", quit)
	mux.HandleFunc("/", home)
	loggedMux := handlers.CombinedLoggingHandler(os.Stdout, mux)
	server := http.Server{
		Addr:     addr,
		Handler:  loggedMux,
		ErrorLog: log.New(os.Stderr, "httpd", log.LstdFlags),
	}
	fmt.Println("Going to listen", addr)
	go func() {
		errors <- server.ListenAndServe()
	}()
	for err := range errors {
		if err == nil {
			continue
		}
		if err != http.ErrServerClosed {
			t.Errorf("unexpected server exit: %v", err)
		}
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	headers(w)
	fmt.Println(r)
}

func headers(w http.ResponseWriter) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
}

func quit(w http.ResponseWriter, r *http.Request) {
	headers(w)
	defer func() {
		//	go func() {
		errors <- server.Close()
		close(errors)
		//	}()
	}()
}

func ws(w http.ResponseWriter, r *http.Request) {
	var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	{
		err := c.WriteMessage(websocket.TextMessage, []byte("text message"))
		if err != nil {
			errors <- fmt.Errorf("write txt: %w", err)
			return
		}
	}

	{
		err := c.WriteMessage(websocket.BinaryMessage, []byte("binary message")) // unsupported so far
		if err != nil {
			errors <- fmt.Errorf("write bin: %w", err)
			return
		}
	}

	{
		mt, message, err := c.ReadMessage()
		if err != nil {
			errors <- fmt.Errorf("read: %w", err)
			return
		}
		log.Printf("recv: %s", message)

		err = c.WriteMessage(mt, message)
		if err != nil {
			errors <- fmt.Errorf("write: %w", err)
			return
		}
	}
}
