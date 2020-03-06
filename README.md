# wasm_websocket 
golang wrapper for using browser WebSocket from WebAssembly build target


## Running Tests

Use [wasmbrowsertest](https://github.com/agnivade/wasmbrowsertest).
```bash
go get github.com/agnivade/wasmbrowsertest
go test -v ./internal/... -count 1 &
GOOS=js GOARCH=wasm go test -v -exec wasmbrowsertest
```
Right now the test case doesn't exit because the server never quits.