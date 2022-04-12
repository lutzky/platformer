package main

import (
	"fmt"
	"image/color"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	playerlib "github.com/lutzky/platformer/player"
	"github.com/lutzky/platformer/rectangle"

	_ "embed"
	_ "image/png"
)

type Game struct {
	gravity float64
}

var game = Game{
	gravity: 0.8,
}

var player = playerlib.New()

const (
	tileHeight   = 30
	tileWidth    = 30
	screenHeight = 240
	screenWidth  = 320
)

var (
	debug = false
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

	if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		debug = !debug
	}

	var dVX float64
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		dVX = player.Acceleration
	} else if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		dVX = -player.Acceleration
	}

	if dVX*player.VX <= 0 {
		// Apply friction unless accelerating in the current movement direction
		player.VX *= (1 - player.Friction)
	}
	player.VX += dVX

	if player.VX > player.MaxVX {
		player.VX = player.MaxVX
	} else if player.VX < -player.MaxVX {
		player.VX = -player.MaxVX
	}

	player.HandleJump(ebiten.IsKeyPressed(ebiten.KeyUp))
}

func (g *Game) checkIsOnFloor() {
	hb := player.Hitbox()
	for _, t := range tiles {
		if hb.Max.X >= t.x && hb.Min.X <= t.x+tileWidth &&
			hb.Max.Y == t.y {
			player.IsOnFloor = true
			return
		}
	}
	player.IsOnFloor = false
}

func (g *Game) handleXCollisions() {
	hb := player.Hitbox()
	for _, t := range tiles {
		if hb.Overlaps(t.rect()) {
			if player.VX > 0 {
				player.SetRight(t.rect().Min.X)
				player.VX = 0
			} else if player.VX < 0 {
				player.SetLeft(t.rect().Max.X)
				player.VX = -0.01
			}
		}
	}
}

func (g *Game) handleYCollisions() {
	hb := player.Hitbox()

	for _, t := range tiles {
		if hb.Overlaps(t.rect()) {
			if player.VY > 0 {
				player.SetBottom(t.rect().Min.Y)
			} else {
				player.SetTop(t.rect().Max.Y)
				player.Scale()
			}
			player.VY = 0
		}
	}
}

func (g *Game) applyGravity() {
	if player.IsOnFloor {
		return
	}
	player.VY += g.gravity
	if player.VY > player.TerminalVelocityY {
		player.VY = player.TerminalVelocityY
	}
}

func (g *Game) Update() error {
	g.checkIsOnFloor()
	g.handleInput()
	player.MoveX(player.VX)
	g.handleXCollisions()
	player.MoveY(player.VY)
	g.applyGravity()
	g.handleYCollisions()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	player.Draw(screen)
	for _, t := range tiles {
		ebitenutil.DrawRect(screen, t.x, t.y, tileWidth, tileHeight, color.White)
	}
	if debug {
		ebitenutil.DebugPrint(screen, player.DebugString())
		player.DrawHitbox(screen)
	}
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
