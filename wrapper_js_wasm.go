package wasm_websocket

import (
	"fmt"
	"syscall/js"
)

// WebSocketArgs arguments https://developer.mozilla.org/en-US/docs/Web/API/WebSocket/WebSocket
type WebSocketArgs struct {
	url       string
	protocols []string
}

func (wsa *WebSocketArgs) args() []interface{} {
	v := []interface{}{
		wsa.url,
	}

	if wsa.protocols != nil && len(wsa.protocols) > 0 {
		v = append(v, wsa.protocols)
	}
	return v
}

// Global wraps a new instance of the global WebSocket implementation
func Global(wsa WebSocketArgs) (*WebSocket, error) {
	return Wrap(js.Global().Get("WebSocket"), wsa)
}

// Wrap a new instance of the provided websocket constructor in our WebSocket
func Wrap(constructor js.Value, wsa WebSocketArgs) (ws *WebSocket, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("constructor error %v", r)
		}
	}()

	if t := constructor.Type(); t != js.TypeFunction {
		return nil, fmt.Errorf("constructor is not js.TypeFunction (was %s)", t)
	}
	v := constructor.New(wsa.args())
	if t := v.Type(); t != js.TypeObject {
		return nil, fmt.Errorf("WebSocket type is not js.TypeObject (was %s)", t)
	}

	ws = &WebSocket{
		value:   v,
		onOpen:  make(chan interface{}),
		onError: make(chan interface{}),
	}

	v.Set("onopen", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		fmt.Println("onopen")
		if len(args) > 0 {
			ws.onOpen <- args[0]
		}
		return nil
	}))

	v.Set("onerror", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		fmt.Println("onerror")
		if len(args) > 0 {
			// https://developer.mozilla.org/en-US/docs/Web/API/Event
			fmt.Println("onerror type", args[0].Get("type").String())

			ws.onError <- args[0]
		}
		return nil
	}))

	// TODO finalizer to reap callbacks

	return
}

// Must is used for simplifying panic chains
func Must(ws *WebSocket, err error) *WebSocket {
	if err != nil {
		panic(fmt.Errorf("wasm_websocket.Must: %w", err))
	}
	return ws
}
