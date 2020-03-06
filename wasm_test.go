package wasm_websocket

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/WIZARDISHUNGRY/wasm_websocket/internal"
)

// avoid t.Parallel
// TestMain doesn't work with wasmbrowsertest
func testURL() string {
	fmt.Println("wss://" + internal.GetServerAddr() + "/ws")
	return "ws://" + internal.GetServerAddr() + "/ws"
}

func quitURL() string {
	return "http://" + internal.GetServerAddr() + "/quit"
}
func upURL() string {
	return "http://" + internal.GetServerAddr() + "/"
}

func TestAaaFirst(t *testing.T) {

	_, err := http.Get(upURL())
	if err != nil {
		t.Fatalf("error connecting to local http service! %v", err)
	}

}

func TestMustGlobal(t *testing.T) {
	ws := Must(Global(WebSocketArgs{url: testURL()}))
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
	ws, err := Global(WebSocketArgs{url: testURL()})
	if err == nil {
		t.Fatalf("nil error returned by Global")
	}
	if ws != nil {
		t.Fatalf("non-nil ws returned by bad call to Global")
	}
	fmt.Println(err.Error())
}

func TestZzzQuitService(t *testing.T) {

	_, err := http.Get(quitURL())
	if err != nil {
		t.Errorf("error closing local http service %v", err.Error())
	}

}
