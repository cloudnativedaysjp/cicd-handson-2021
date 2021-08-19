# CICD Landscape

![screenshot](./screenshot.png)

- Port: 9090
- Endpoints
  - /
  - /health
  - /landscape

## macOS, Linux

```bash
# build wasm
GOOS=js GOARCH=wasm go build -o web/static/main.wasm ./wasm

# run test
go test ./server

# run server
go run ./server
```

## Windows

```pwsh
# run test
go test .\server\

# run server
go run .\server\
```
