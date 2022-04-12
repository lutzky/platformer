package player

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/lutzky/platformer/rectangle"
	"github.com/lutzky/platformer/sprites"
)

type Player struct {
	VX, VY float64
	rect   rectangle.Rectangle[float64]

	IsOnFloor, isJumping         bool
	color, colorJump, colorFloor color.Color

	Friction, Acceleration float64
	MaxVX                  float64
	TerminalVelocityY      float64

	jumpSpeed      float64
	jumpHoverSpeed float64
	jumpStarted    bool

	scaling float64

	scalingFactorOnHit float64

	marginTop, marginLeft, marginBottom, marginRight float64
}

func New() Player {
	return Player{
		color:      color.RGBA{0, 0, 255, 255},
		colorJump:  color.RGBA{255, 0, 255, 255},
		colorFloor: color.RGBA{0, 255, 255, 255},

		rect: rectangle.Rect[float64](0, 0, 32, 32),

		Friction:          0.08,
		Acceleration:      0.3,
		MaxVX:             5,
		TerminalVelocityY: 9,

		jumpSpeed:      12,
		jumpHoverSpeed: 3,

		scaling:            1,
		scalingFactorOnHit: 0.9,

		marginTop:    6,
		marginLeft:   6,
		marginRight:  6,
		marginBottom: 0,
	}
}

func (p *Player) SetLeft(x float64) {
	p.rect.SetLeft(x - math.Ceil(p.marginLeft*p.scaling))
}

func (p *Player) SetRight(x float64) {
	p.rect.SetRight(x + math.Ceil(p.marginRight*p.scaling))
}

func (p *Player) SetTop(y float64) {
	p.rect.SetTop(y - math.Ceil(p.marginTop*p.scaling))
}

func (p *Player) SetBottom(y float64) {
	p.rect.SetBottom(y - math.Ceil(p.marginBottom*p.scaling))
}

func (p *Player) MoveX(dX float64) {
	p.rect.MoveX(dX)
}

func (p *Player) MoveY(dY float64) {
	p.rect.MoveY(dY)
}

func (p *Player) Hitbox() rectangle.Rectangle[float64] {
	return rectangle.Rect(
		p.rect.Min.X+math.Ceil(p.marginLeft*p.scaling),
		p.rect.Min.Y+math.Ceil(p.marginTop*p.scaling),
		p.rect.Max.X-math.Ceil(p.marginRight*p.scaling),
		p.rect.Max.Y-math.Ceil(p.marginBottom*p.scaling),
	)
}

func (p *Player) Draw(dst *ebiten.Image) {
	var sprite *sprites.Sprite
	switch {
	case p.VY < -0.2:
		sprite = sprites.Jump
	case p.VY > 0.2:
		sprite = sprites.Fall
	case math.Abs(p.VX) > 0.2:
		sprite = sprites.Run
	default:
		sprite = sprites.Idle
	}

	op := &ebiten.DrawImageOptions{}
	if p.VX < 0 {
		op.GeoM.Scale(-1, 1)
		op.GeoM.Translate(32, 0)
	}
	op.GeoM.Scale(p.scaling, p.scaling)
	op.GeoM.Translate(p.rect.Min.X, p.rect.Min.Y)
	dst.DrawImage(sprite.GetFrame(), op)
}

func (p *Player) DrawHitbox(dst *ebiten.Image) {
	hb := p.Hitbox()
	c := color.White
	ebitenutil.DrawLine(dst, hb.Min.X, hb.Min.Y, hb.Max.X, hb.Min.Y, c)
	ebitenutil.DrawLine(dst, hb.Min.X, hb.Min.Y, hb.Min.X, hb.Max.Y, c)
	ebitenutil.DrawLine(dst, hb.Min.X, hb.Max.Y, hb.Max.X, hb.Max.Y, c)
	ebitenutil.DrawLine(dst, hb.Max.X, hb.Min.Y, hb.Max.X, hb.Max.Y, c)
}

func (p *Player) HandleJump(jumpKeyDown bool) {
	if p.VY >= 0 {
		p.isJumping = false
	}
	if jumpKeyDown {
		if !p.jumpStarted {
			p.jumpStarted = true
			if p.IsOnFloor {
				p.VY = -p.jumpSpeed
				p.isJumping = true
			}
		}
	} else {
		p.jumpStarted = false
		if p.isJumping && p.VY < -p.jumpHoverSpeed {
			p.VY = -p.jumpHoverSpeed
		}
		p.isJumping = false
	}
}

func (p *Player) Scale() {
	p.rect.Scale(p.scalingFactorOnHit)
	p.scaling *= p.scalingFactorOnHit
}

func (p *Player) DebugString() string {
	return fmt.Sprintf("%.2f V: (%3.1f,%3.1f)\nHITBOX: %.2f",
		p.rect, p.VX, p.VY, p.Hitbox())
}
