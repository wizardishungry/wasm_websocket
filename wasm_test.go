package wasm_websocket

import (
	"testing"
)

func TestMustGlobal(t *testing.T) {
	t.Parallel()
	ws := Must(Global(WebSocketArgs{url: "wss://test.example.com/ws"}))
	if ws == nil {
		t.Fatalf("nil returned by Must")
	}
}

func TestDoesntPanicOnConstructorError(t *testing.T) {
	t.Parallel()
	ws, err := Global(WebSocketArgs{url: "http://test.example.com/ws"})
	if err == nil {
		t.Fatalf("nil error returned by Global")
	}
	if ws != nil {
		t.Fatalf("non-nil ws returned by bad call to Global")
	}
}
