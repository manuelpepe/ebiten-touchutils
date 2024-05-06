# Ebiten Touchutils

Utilities to detect and handle mobile gestures.

It currently supports:

- Taps with 1, 2 or 3 fingers
- Pinch inwards and outwards
- Two finger pan (up, down, left, right)


## Demo

You can live test the demo app at:

https://blog.manuelpepe.com/demo-gestures.html

You'll of course need a touch screen.


### Build demo

```
mkdir dist/
GOOS=js GOARCH=wasm go build -o dist/demo.wasm demo/main.go
```

## Usage

I recommend the following import rename for simplicity:

```go
import (
	touchutils "github.com/manuelpepe/ebiten-touchutils"
)
```

Then you'll need to create a `TouchTracker` in your game and update it every tick:

```go
type Game struct {
    touch *touchutils.TouchTracker
}

func (g *Game) Update() {
    g.touch.Update()
}

func main() {
    g := &Game{touch: touchutils.NewTouchTracker()}
    if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
```

For examples on usage check out [the demo code](./demo/main.go).