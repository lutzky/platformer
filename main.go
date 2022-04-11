package main

import (
	"fmt"
	"image/color"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"github.com/lutzky/platformer/rectangle"
)

type Game struct {
	gravity float64
}

var game = Game{
	gravity: 0.8,
}

type Player struct {
	vX, vY                 float64
	rect                   rectangle.Rectangle[float64]
	isOnFloor              bool
	color                  color.Color
	friction, acceleration float64
	maxVX                  float64
}

var player = Player{
	color:        color.RGBA{0, 0, 255, 255},
	rect:         rectangle.Rect[float64](0, 0, 30, 30),
	friction:     0.08,
	acceleration: 0.3,
	maxVX:        5,
}

func (player *Player) draw(dst *ebiten.Image) {
	ebitenutil.DrawRect(dst, player.rect.Min.X, player.rect.Min.Y,
		player.rect.Width(), player.rect.Height(), player.color)
}

const (
	tileHeight   = 30
	tileWidth    = 30
	screenHeight = 240
	screenWidth  = 320
)

var (
	jumpPressed = false
)

var tileMap = []string{
	"..........",
	"..........",
	"..xxxx....",
	"..........",
	".x.......x",
	"..x.....xx",
	"...x...xxx",
	"......xxxx",
	"xxxxxxxxxx",
}

type tile struct {
	x, y float64
}

func (t tile) rect() rectangle.Rectangle[float64] {
	return rectangle.Rect(t.x, t.y, t.x+tileWidth, t.y+tileHeight)
}

var tiles []tile

func loadTiles() {
	for row, s := range tileMap {
		for col, c := range s {
			if c == 'x' {
				tiles = append(tiles, tile{
					x: float64(tileWidth * col),
					y: float64(tileHeight * row),
				})
			}
		}
	}
}

func (g *Game) handleInput() {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		fmt.Println("Exiting")
		os.Exit(0)
	}

	var dVX float64
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		dVX = player.acceleration
	} else if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		dVX = -player.acceleration
	}

	if dVX*player.vX <= 0 {
		// Apply friction unless accelerating in the current movement direction
		player.vX *= (1 - player.friction)
	}
	player.vX += dVX

	if player.vX > player.maxVX {
		player.vX = player.maxVX
	} else if player.vX < -player.maxVX {
		player.vX = -player.maxVX
	}

	if ebiten.IsKeyPressed(ebiten.KeySpace) || ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		if !jumpPressed {
			jumpPressed = true
			if player.isOnFloor {
				player.vY = -16
			}
		}
	} else {
		jumpPressed = false
	}
}

func (g *Game) checkIsOnFloor() {
	for _, t := range tiles {
		if player.rect.Max.X >= t.x && player.rect.Min.X <= t.x+tileWidth &&
			player.rect.Max.Y == t.y {
			player.isOnFloor = true
			return
		}
	}
	player.isOnFloor = false
}

func (g *Game) handleXCollisions() {
	for _, t := range tiles {
		if player.rect.Overlaps(t.rect()) {
			if player.vX > 0 {
				player.rect.SetRight(t.rect().Min.X)
			} else if player.vX < 0 {
				player.rect.SetLeft(t.rect().Max.X)
			}
			player.vX = 0
		}

	}

	if player.rect.Max.X > screenWidth {
		player.rect.SetRight(screenWidth)
		player.vX *= -1
	} else if player.rect.Min.X < 0 {
		player.rect.SetLeft(0)
		player.vX *= -1
	}
}

func (g *Game) handleYCollisions() {
	for _, t := range tiles {
		if player.rect.Overlaps(t.rect()) {
			if player.vY > 0 {
				player.rect.SetBottom(t.rect().Min.Y)
			} else {
				player.rect.SetTop(t.rect().Max.Y)
				player.rect.Scale(0.9)
			}
			player.vY = 0
		}
	}
}

func (g *Game) applyGravity() {
	if player.isOnFloor {
		return
	}
	player.vY += g.gravity
}

func (g *Game) Update() error {
	g.checkIsOnFloor()
	g.handleInput()
	player.rect.MoveX(player.vX)
	g.handleXCollisions()
	player.rect.MoveY(player.vY)
	g.applyGravity()
	g.handleYCollisions()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	player.draw(screen)
	for _, t := range tiles {
		ebitenutil.DrawRect(screen, t.x, t.y, tileWidth, tileHeight, color.White)
	}
	ebitenutil.DebugPrint(screen, fmt.Sprintf("Pos: (%3.0f,%3.0f) V: (%3.1f,%3.1f), IoF: %t",
		player.rect.Min.X, player.rect.Min.Y, player.vX, player.vY, player.isOnFloor))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	log.Print("Initializing, setting window size")
	ebiten.SetWindowSize(2*screenWidth, 2*screenHeight)
	log.Print("Setting window title")
	ebiten.SetWindowTitle("Hello, World!")
	log.Print("Loading tiles")
	loadTiles()
	log.Print("Running game")
	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}
