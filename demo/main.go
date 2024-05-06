package main

import (
	"bytes"
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	touchutils "github.com/manuelpepe/ebiten-touchutils"
)

var (
	mplusFaceSource *text.GoTextFaceSource
)

func init() {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))
	if err != nil {
		log.Fatal(err)
	}
	mplusFaceSource = s
}

type Gesture struct {
	w, h int

	touch *touchutils.TouchTracker

	tapCounter int
	tapMessage string
}

func NewGestureDemo(width, height int) *Gesture {

	return &Gesture{
		w: width,
		h: height,

		touch: touchutils.NewTouchTracker(),
	}
}

// Max TPS as specified by ebiten
const MAX_TPS = 60

// Delay between updates
const DELAY_SEC = 1

func (g *Gesture) Update() error {
	g.touch.Update()
	return nil
}

func (g *Gesture) Draw(screen *ebiten.Image) {
	msgs := make([]string, 4)

	if g.tapMessage != "" {
		msgs = append(msgs, g.tapMessage)
		g.tapCounter++
		if g.tapCounter > 60 {
			g.tapMessage = ""
			g.tapCounter = 0
		}
	}

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		vector.DrawFilledCircle(screen, float32(x), float32(y), 5, color.RGBA{0, 0, 255, 1}, true)
		msgs = append(msgs, "left mouse button")
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
		x, y := ebiten.CursorPosition()
		vector.DrawFilledCircle(screen, float32(x), float32(y), 5, color.RGBA{0, 255, 0, 1}, true)
		msgs = append(msgs, "right mouse button")
	}

	if _, _, _, ok := g.touch.TappedThree(); ok {
		g.tapMessage = "tapped three"
		g.tapCounter = 0
	} else if _, _, ok := g.touch.TappedTwo(); ok {
		g.tapMessage = "tapped two"
		g.tapCounter = 0
	} else if _, ok := g.touch.TappedOne(); ok {
		g.tapMessage = "tapped one"
		g.tapCounter = 0
	}

	if g.touch.IsTouchingThree() {
		msgs = append(msgs, "touching three")
	} else if g.touch.IsTouchingTwo() {
		msgs = append(msgs, "touching two")
		if pan, ok := g.touch.TwoFingerPan(); ok {
			if pan.IsHorizontal() {
				msgs = append(msgs, "horizontal pan")
				deltaX := pan.OriginX - pan.LastX
				if deltaX < -10 {
					msgs = append(msgs, fmt.Sprintf("swipe right - delta: %d", deltaX))
				} else if deltaX > 10 {
					msgs = append(msgs, fmt.Sprintf("swipe left - delta: %d", deltaX))
				}

				vector.DrawFilledCircle(screen, float32(pan.OriginX), float32(g.h)/2, 5, color.RGBA{255, 0, 0, 1}, true)
				vector.DrawFilledCircle(screen, float32(pan.LastX), float32(g.h)/2, 5, color.RGBA{0, 255, 0, 1}, true)
				vector.StrokeLine(screen, float32(pan.OriginX), float32(g.h)/2, float32(pan.LastX), float32(g.h)/2, 1, color.White, true)
			}

			if pan.IsVertical() {
				msgs = append(msgs, "vertical pan")
				deltaY := pan.OriginY - pan.LastY
				if deltaY < -10 {
					msgs = append(msgs, fmt.Sprintf("swipe up - delta: %d", deltaY))
				} else if deltaY > 10 {
					msgs = append(msgs, fmt.Sprintf("swipe down - delta: %d", deltaY))
				}

				vector.DrawFilledCircle(screen, float32(g.w)/2, float32(pan.OriginY), 5, color.RGBA{255, 0, 0, 1}, true)
				vector.DrawFilledCircle(screen, float32(g.w)/2, float32(pan.LastY), 5, color.RGBA{0, 255, 0, 1}, true)
				vector.StrokeLine(screen, float32(g.w)/2, float32(pan.OriginY), float32(g.w)/2, float32(pan.LastY), 1, color.White, true)
			}
		}

		if pinch, ok := g.touch.Pinch(); ok {
			if pinch.IsInward() {
				msgs = append(msgs, "inward pinch")
			}
			if pinch.IsOutward() {
				msgs = append(msgs, "outward pinch")
			}

			vector.DrawFilledCircle(screen, float32(pinch.CenterX)-float32(pinch.Distance/2), float32(pinch.CenterY), 5, color.RGBA{255, 0, 0, 1}, true)
			vector.DrawFilledCircle(screen, float32(pinch.CenterX)+float32(pinch.Distance/2), float32(pinch.CenterY), 5, color.RGBA{0, 255, 0, 1}, true)
			vector.StrokeLine(screen, float32(pinch.CenterX)-float32(pinch.Distance/2), float32(pinch.CenterY), float32(pinch.CenterX)+float32(pinch.Distance/2), float32(pinch.CenterY), 1, color.White, true)
		}
	} else if g.touch.IsTouching() {
		x, y, _ := g.touch.GetFirstTouchPosition()
		msgs = append(msgs, "touching one")
		vector.DrawFilledCircle(screen, float32(x), float32(y), 5, color.RGBA{0, 0, 255, 1}, true)
	}

	for ix, m := range msgs {
		textFace := &text.GoTextFace{
			Source: mplusFaceSource,
			Size:   24,
		}
		w, h := text.Measure(m, textFace, 0)
		textOp := &text.DrawOptions{}
		textOp.GeoM.Translate(w/2+10, h*float64(ix+1))
		textOp.PrimaryAlign = text.AlignCenter
		textOp.SecondaryAlign = text.AlignCenter
		text.Draw(screen, m, textFace, textOp)
	}

	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %.2f\nFPS: %.2f", ebiten.ActualTPS(), ebiten.ActualFPS()))
}

func (g *Gesture) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func main() {
	W, H := 300, 500
	ebiten.SetWindowSize(W, H)
	ebiten.SetWindowTitle("Hello, World!")
	game := NewGestureDemo(W, H)
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
