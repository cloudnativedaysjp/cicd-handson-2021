#!/bin/bash

GOOS=js GOARCH=wasm go build -o web/static/main.wasm ./wasm
