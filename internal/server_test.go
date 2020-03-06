package internal

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

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
}

func quit(w http.ResponseWriter, r *http.Request) {
	go func() {
		time.Sleep(time.Second)
		errors <- server.Close()
		close(errors)
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
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}
