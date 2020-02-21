package wasm_websocket

import (
	"fmt"
	"testing"
)

// avoid t.Parallel https://travis-ci.community/t/goos-js-goarch-wasm-go-run-fails-panic-newosproc-not-implemented/1651

func TestMustGlobal(t *testing.T) {
	ws := Must(Global(WebSocketArgs{url: "wss://test.example.com/ws"}))
	if ws == nil {
		t.Fatalf("nil returned by Must")
	}
	defer ws.Close()

	for {
		select {
		case e := <-ws.OnClose():
			fmt.Println("OnClose! ", e)
			return
		case e := <-ws.OnError():
			fmt.Println("OnError ", e)
		case e := <-ws.OnOpen():
			fmt.Println("onOpen ", e)
		case e := <-ws.OnMessage():
			fmt.Println("OnMessage ", e)
		}
	}
}

func TestDoesntPanicOnConstructorError(t *testing.T) {
	ws, err := Global(WebSocketArgs{url: "http://test.example.com/ws"})
	if err == nil {
		t.Fatalf("nil error returned by Global")
	}
	if ws != nil {
		t.Fatalf("non-nil ws returned by bad call to Global")
	}
	fmt.Println(err.Error())
}
