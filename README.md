# Ebiten Touchutils

Utilities to detect and handle mobile gestures.

It currently supports:

- Taps with 1, 2 or 3 fingers
- Pinch inwards and outwards
- Two finger pan (up, down, left, right)


## Demo

You can live test the demo app at:

https://blog.manuelpepe.com/demo-gestures.html


### Build demo

```
mkdir dist/
GOOS=js GOARCH=wasm go build -o dist/demo.wasm demo/main.go
```