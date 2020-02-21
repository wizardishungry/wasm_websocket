package wasm_websocket

import (
	"syscall/js"
)

type WebSocket struct {
	value                               js.Value
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
