package wasm_websocket

import (
	"syscall/js"
)

type WebSocket struct {
	value js.Value

	onOpen, onClose, onMessage, onError chan interface{}
}

// OnOpen
func (ws *WebSocket) OnOpen() <-chan interface{} {
	return ws.onOpen
}

// OnError
func (ws *WebSocket) OnError() <-chan interface{} {
	return ws.onError
}

// OnMessage
func (ws *WebSocket) OnMessage() <-chan interface{} {
	return ws.onMessage
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
