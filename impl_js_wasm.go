package wasm_websocket

import (
	"fmt"
	"syscall/js"
)

type WebSocket struct {
	value js.Value

	onOpen, onClose, onMessage, onError chan interface{}

	onOpenC    chan struct{}
	onMessageC chan string
}

// OnOpen
func (ws *WebSocket) OnOpen() <-chan struct{} {
	return ws.onOpenC
}

// OnError
func (ws *WebSocket) OnError() <-chan interface{} {
	return ws.onError
}

// OnMessage
func (ws *WebSocket) OnMessage() <-chan string {
	return ws.onMessageC
}

// OnClose
func (ws *WebSocket) OnClose() <-chan interface{} {
	return ws.onClose
}

// ReadyState https://developer.mozilla.org/en-US/docs/Web/API/WebSocket/readyState
type ReadyState uint16

const (
	ReadyStateConnecting = iota
	ReadyStateOpen
	ReadyStateClosing
	ReadyStateClosed
)

// ReadyState
func (ws *WebSocket) ReadyState() ReadyState {
	return ReadyState(ws.value.Get("readyState").Int())
}

// BufferedAmount
func (ws *WebSocket) BufferedAmount() int {
	return ws.value.Get("bufferedAmount").Int()
}

// Close
func (ws *WebSocket) Close() {
	ws.value.Call("close")
}

// Send
func (ws *WebSocket) Send(v interface{}) (err error) {
	switch t := v.(type) {
	case string:
	default:
		return fmt.Errorf("type unsupported %T!", t)
	}

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("send error %v", r)
		}
	}()

	ws.value.Call("send", js.ValueOf(v))
	return nil
}
