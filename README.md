# wasm_websocket 
golang wrapper for using browser WebSocket from WebAssembly build target


## Running Tests

Use [wasmbrowsertest](https://github.com/agnivade/wasmbrowsertest).
```bash
go get github.com/agnivade/wasmbrowsertest
GOOS=js GOARCH=wasm go test -exec wasmbrowsertest
```