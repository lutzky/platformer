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

const (
	maxVX    = 5
	friction = 0.05
)

var (
	playerColor = color.RGBA{0, 0, 255, 255}
	isOnFloor   = true
	jumpPressed = false
)

var tileMap = []string{
	"..........",
	"..........",
	"..........",
	"..........",
	"..........",
	"..........",
	"..........",
	"......x...",
	//"xxxxxxxxxx",
}

type tile struct {
	x, y float64
}

var tiles []tile

func loadTiles() {
	for row, s := range tileMap {
		for col, c := range s {
			if c == 'x' {
				tiles = append(tiles, tile{
					x: float64(30 * col),
					y: float64(240 - 30 - 30*row),
				})
			}
		}
	}
}

func (g *Game) handleInput() {
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.vX += 0.1
	} else if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.vX -= 0.1
	} else {
		g.vX *= (1 - friction)
	}
	if g.vX > maxVX {
		g.vX = maxVX
	} else if g.vX < -maxVX {
		g.vX = -maxVX
	}

	if ebiten.IsKeyPressed(ebiten.KeySpace) || ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		if !jumpPressed {
			jumpPressed = true
			if isOnFloor {
				g.vY = 10
			}
		}
	} else {
		jumpPressed = false
	}
}

func (g *Game) overlaps(t tile) bool {
	// TODO(lutzky): This is ugly, we need to pass the player, not the game

	// TODO(lutzky): Things actually aren't pixel-perfect; the tiles are 30x30 and they are
	// positions 0..30,30..60 - i.e. they have an overlap

	return g.RectX+30 >= t.x && g.RectX <= t.x+30 &&
		g.RectY+30 >= t.y && g.RectY <= t.y+30
}

func (g *Game) handleXCollisions() {
	for _, t := range tiles {
		if g.overlaps(t) {
			if g.vX > 0 {
				g.RectX = t.x - 30
			} else {
				g.RectX = t.x + 30
			}
			g.vX = 0
		}

	}

	if g.RectX > 320-30 {
		g.RectX = 320 - 30
		g.vX = -g.vX
	} else if g.RectX < 0 {
		g.RectX = 0
		g.vX = -g.vX
	}
}

func (g *Game) handleYCollisions() {
	for _, t := range tiles {
		if g.overlaps(t) {
			if g.vY > 0 {
				g.RectY = t.y - 30
			} else {
				g.RectY = t.y + 30
				isOnFloor = true
			}
			g.vY = 0
		}
	}

	if g.RectY <= 0 {
		g.vY = 0
		g.RectY = 0
		isOnFloor = true
	}
}

func (g *Game) applyGravity() {
	g.vY -= 1.0
	g.vY *= 0.99
}

func (g *Game) Update() error {
	g.handleInput()
	g.RectX += g.vX
	g.handleXCollisions()
	g.RectY += g.vY
	g.applyGravity()
	g.handleYCollisions()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DrawRect(screen, g.RectX, 240-30-g.RectY, 30, 30, playerColor)
	for _, t := range tiles {
		ebitenutil.DrawRect(screen, t.x, 240-30-t.y, 30, 30, color.White)
	}
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
	log.Print("Loading tiles")
	loadTiles()
	log.Print("Running game")
	if err := ebiten.RunGame(&Game{
		RectX: 30,
		RectY: 100,
		vX:    0,
	}); err != nil {
		log.Fatal(err)
	}
}
