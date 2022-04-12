package main

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"github.com/lutzky/platformer/rectangle"
	"github.com/lutzky/platformer/sprites"

	_ "embed"
	_ "image/png"
)

type Game struct {
	gravity float64
}

var game = Game{
	gravity: 0.8,
}

type Player struct {
	vX, vY float64
	rect   rectangle.Rectangle[float64]

	isOnFloor, isJumping         bool
	color, colorJump, colorFloor color.Color

	friction, acceleration float64
	maxVX                  float64
	terminalVelocityY      float64

	jumpSpeed      float64
	jumpHoverSpeed float64

	scaling       float64
	scalingFactor float64 // TODO(lutzky): Rename me

	marginTop, marginLeft, marginBottom, marginRight float64
}

var player = Player{
	color:      color.RGBA{0, 0, 255, 255},
	colorJump:  color.RGBA{255, 0, 255, 255},
	colorFloor: color.RGBA{0, 255, 255, 255},

	rect: rectangle.Rect[float64](0, 0, 32, 32),

	friction:          0.08,
	acceleration:      0.3,
	maxVX:             5,
	terminalVelocityY: 9,

	jumpSpeed:      12,
	jumpHoverSpeed: 3,

	scaling:       1,
	scalingFactor: 1, // 0.9, TODO(lutzky): Set me back to 0.9

	marginTop:    6,
	marginLeft:   6,
	marginRight:  6,
	marginBottom: 0,
}

func (p *Player) SetLeft(x float64) {
	p.rect.SetLeft(x - p.marginLeft)
}
func (p *Player) SetRight(x float64) {
	p.rect.SetRight(x + p.marginRight)
}
func (p *Player) SetTop(y float64) {
	p.rect.SetTop(y - p.marginTop)
}
func (p *Player) SetBottom(y float64) {
	p.rect.SetBottom(y + p.marginBottom)
}

func (player *Player) hitbox() rectangle.Rectangle[float64] {
	return rectangle.Rect(
		player.rect.Min.X+player.marginLeft,
		player.rect.Min.Y+player.marginTop,
		player.rect.Max.X-player.marginRight,
		player.rect.Max.Y-player.marginBottom,
	)
}

func (player *Player) draw(dst *ebiten.Image) {
	var sprite *sprites.Sprite
	switch {
	case player.vY < -0.2:
		sprite = sprites.Jump
	case player.vY > 0.2:
		sprite = sprites.Fall
	case math.Abs(player.vX) > 0.2:
		sprite = sprites.Run
	default:
		sprite = sprites.Idle
	}

	op := &ebiten.DrawImageOptions{}
	if player.vX < 0 {
		op.GeoM.Scale(-1, 1)
		op.GeoM.Translate(32, 0)
	}
	op.GeoM.Scale(player.scaling, player.scaling)
	op.GeoM.Translate(player.rect.Min.X, player.rect.Min.Y)
	dst.DrawImage(sprite.GetFrame(), op)
	if debug {
		player.drawHitbox(dst)
	}
}

func (player *Player) drawHitbox(dst *ebiten.Image) {
	hb := player.hitbox()
	c := color.White
	ebitenutil.DrawLine(dst, hb.Min.X, hb.Min.Y, hb.Max.X, hb.Min.Y, c)
	ebitenutil.DrawLine(dst, hb.Min.X, hb.Min.Y, hb.Min.X, hb.Max.Y, c)
	ebitenutil.DrawLine(dst, hb.Min.X, hb.Max.Y, hb.Max.X, hb.Max.Y, c)
	ebitenutil.DrawLine(dst, hb.Max.X, hb.Min.Y, hb.Max.X, hb.Max.Y, c)
}

const (
	tileHeight   = 30
	tileWidth    = 30
	screenHeight = 240
	screenWidth  = 320
)

var (
	jumpPressed = false
	debug       = false
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

	debug = ebiten.IsKeyPressed(ebiten.KeyD)

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

	if player.vY >= 0 {
		player.isJumping = false
	}
	if ebiten.IsKeyPressed(ebiten.KeySpace) || ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		if !jumpPressed {
			jumpPressed = true
			if player.isOnFloor {
				player.vY = -player.jumpSpeed
				player.isJumping = true
			}
		}
	} else {
		jumpPressed = false
		if player.isJumping && player.vY < -player.jumpHoverSpeed {
			player.vY = -player.jumpHoverSpeed
		}
		player.isJumping = false
	}
}

func (g *Game) checkIsOnFloor() {
	// TODO(lutzky): This uses hitbox inconsistently, and doesn't use rectangles
	// for tiles
	hb := player.hitbox()
	for _, t := range tiles {
		if hb.Max.X >= t.x && hb.Min.X <= t.x+tileWidth &&
			hb.Max.Y == t.y {
			player.isOnFloor = true
			return
		}
	}
	player.isOnFloor = false
}

func (g *Game) handleXCollisions() {
	hitbox := player.hitbox()
	for _, t := range tiles {
		if hitbox.Overlaps(t.rect()) {
			if player.vX > 0 {
				player.SetRight(t.rect().Min.X)
				player.vX = 0
			} else if player.vX < 0 {
				player.SetLeft(t.rect().Max.X)
				player.vX = -0.01
			}
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
		if player.hitbox().Overlaps(t.rect()) {
			if player.vY > 0 {
				player.SetBottom(t.rect().Min.Y)
			} else {
				player.SetTop(t.rect().Max.Y)
				player.rect.Scale(player.scalingFactor)
				player.scaling *= player.scalingFactor
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
	if player.vY > player.terminalVelocityY {
		player.vY = player.terminalVelocityY
	}
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
	if debug {
		ebitenutil.DebugPrint(screen, fmt.Sprintf("Pos: (%3.0f,%3.0f) V: (%3.1f,%3.1f), IoF: %t",
			player.rect.Min.X, player.rect.Min.Y, player.vX, player.vY, player.isOnFloor))
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
