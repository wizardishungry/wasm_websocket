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

func TestIntegrationScenario(t *testing.T) {
	testService(t)

	ws := Must(Global(WebSocketArgs{url: testURL()}))
	if ws == nil {
		t.Fatalf("nil returned by Must")
	}
	defer ws.Close()

	for {
		select {
		// TODO having a timeout case is weird in the wasm go runtime
		case e := <-ws.OnClose():
			fmt.Println("OnClose! ", e)
			return
		case e := <-ws.OnError():
			fmt.Println("OnError ", e)
		case e := <-ws.OnOpen():
			fmt.Println("onOpen ", e)
		case e := <-ws.OnMessage():
			fmt.Println("OnMessage ", e)
			ws.Send("on message back")
		}
	}
}

func TestDoesntPanicOnConstructorError(t *testing.T) {
	ws, err := Global(WebSocketArgs{url: "wdsfsdfs://fu.example.com/"})
	if err == nil {
		t.Errorf("nil error returned by Global")
	} else {
		fmt.Println(err.Error())
	}
	if ws != nil {
		t.Errorf("non-nil ws returned by bad call to Global")
	}
}

func testService(t *testing.T) {
	_, err := http.Get(upURL())
	if err != nil {
		t.Fatalf("error connecting to local http service! %v", err)
	}

}

func quitService(t *testing.T) {
	_, err := http.Get(quitURL())
	if err != nil {
		t.Errorf("error closing local http service %v", err.Error())
	}

}

func TestExitService(t *testing.T) {
	testService(t)
	defer quitService(t)
}
