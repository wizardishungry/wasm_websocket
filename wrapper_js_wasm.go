package wasm_websocket

import (
	"encoding/json"
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
		value: v,
	}

	ws.onOpen = ws.regCb("onopen")
	ws.onError = ws.regCb("onerror")
	ws.onMessage = ws.regCb("onmessage")
	ws.onClose = ws.regCb("onclose")

	// TODO finalizer to reap callbacks

	return
}

func (ws *WebSocket) regCb(call string) chan interface{} {
	c := make(chan interface{})
	ws.value.Set(call, js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		fmt.Println(call)
		if len(args) > 0 {
			// https://developer.mozilla.org/en-US/docs/Web/API/Event
			fmt.Println(call, "type", args[0].Get("type").String())
			fmt.Println(call, "ReadyState()", ws.ReadyState())

			if m, err := asMap(args[0]); err == nil {
				c <- m
			} else {
				fmt.Println("asMap error", err)
			}
		}
		return nil
	}))
	return c
}

// Must is used for simplifying panic chains
func Must(ws *WebSocket, err error) *WebSocket {
	if err != nil {
		panic(fmt.Errorf("wasm_websocket.Must: %w", err))
	}
	return ws
}

// asMap holy hacks
func asMap(v js.Value) (map[string]interface{}, error) {
	// todo panic/recover?
	s := js.Global().Get("JSON").Call("stringify", v).String()
	m := make(map[string]interface{})
	err := json.Unmarshal([]byte(s), &m)
	return m, err
}
