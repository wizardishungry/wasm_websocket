package wasm_websocket

import (
	"fmt"
	"testing"
)

// avoid t.Parallel https://travis-ci.community/t/goos-js-goarch-wasm-go-run-fails-panic-newosproc-not-implemented/1651

func TestMustGlobal(t *testing.T) {
	return
	ws := Must(Global(WebSocketArgs{url: "wss://test.example.com/ws"}))
	if ws == nil {
		t.Fatalf("nil returned by Must")
	}
	select {
	case e := <-ws.OnError():
		fmt.Println("onError ", e)
	case e := <-ws.OnOpen():
		fmt.Println("onOpen ", e)
	}
}

func TestDoesntPanicOnConstructorError(t *testing.T) {
	return

	ws, err := Global(WebSocketArgs{url: "http://test.example.com/ws"})
	if err == nil {
		t.Fatalf("nil error returned by Global")
	}
	if ws != nil {
		t.Fatalf("non-nil ws returned by bad call to Global")
	}
}
