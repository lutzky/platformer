package sprites

import (
	"bytes"
	_ "embed"
	"image"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
)

type Sprite struct {
	frame         int
	frames        int
	frameDuration int
	frameWidth    int
	image         *ebiten.Image
}

func (s *Sprite) GetFrame() *ebiten.Image {
	s.frame = (s.frame + 1) % (s.frameDuration * s.frames)
	offset := s.frameWidth * (s.frame / s.frameDuration)
	r := image.Rect(offset, 0, offset+s.frameWidth, s.image.Bounds().Max.Y)
	return s.image.SubImage(r).(*ebiten.Image)
}

func loadSprite(b []byte, frameWidth, frameDuration int) *Sprite {
	img, _, err := image.Decode(bytes.NewReader(b))
	if err != nil {
		panic(err)
	}

	return &Sprite{
		frames:        img.Bounds().Max.X / frameWidth,
		frameDuration: frameDuration,
		frameWidth:    frameWidth,
		image:         ebiten.NewImageFromImage(img),
	}
}

var (
	// Sprites are from https://pixelfrog-assets.itch.io/pixel-adventure-1

	//go:embed "Run (32x32).png"
	runBytes []byte
	Run      = loadSprite(runBytes, 32, 5)

	//go:embed "Idle (32x32).png"
	idleBytes []byte
	Idle      = loadSprite(idleBytes, 32, 5)

	//go:embed "Jump (32x32).png"
	jumpBytes []byte
	Jump      = loadSprite(jumpBytes, 32, 5)

	//go:embed "Fall (32x32).png"
	fallBytes []byte
	Fall      = loadSprite(fallBytes, 32, 5)
)
