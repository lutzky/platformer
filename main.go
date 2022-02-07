package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Game struct {
	RectX float64
	RectY float64

	vX float64
	vY float64
}

func (g *Game) Update() error {
	g.RectX += g.vX
	g.RectY += g.vY
	g.vY -= 1.0
	g.vY *= 0.99
	if g.RectY < 0 {
		g.vY = -g.vY
	}
	if g.RectX > 320-30 || g.RectX < 0 {
		g.vX = -g.vX
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DrawRect(screen, g.RectX, 240-30-g.RectY, 30, 30, color.White)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("Pos: (%3.0f,%3.0f) V: (%3.1f,%3.1f)", g.RectX, g.RectY, g.vX, g.vY))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func main() {
	log.Print("Initializing, setting window size")
	ebiten.SetWindowSize(640, 480)
	log.Print("Setting window title")
	ebiten.SetWindowTitle("Hello, World!")
	log.Print("Running game")
	if err := ebiten.RunGame(&Game{
		RectX: 30,
		RectY: 100,
		vX:    3.5,
	}); err != nil {
		log.Fatal(err)
	}
}
