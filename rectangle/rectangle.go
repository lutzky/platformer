// Package rectangle is a generic version of image.Rectangle
package rectangle

type number interface {
	~int | ~float32 | ~float64
}

type point[C number] struct {
	X, Y C
}

type Rectangle[C number] struct {
	Min, Max point[C]
}

func Rect[C number](x0, y0, x1, y1 C) Rectangle[C] {
	return Rectangle[C]{
		Min: point[C]{x0, y0},
		Max: point[C]{x1, y1},
	}
}

func (r Rectangle[C]) Width() C {
	return r.Max.X - r.Min.X
}

func (r Rectangle[C]) Height() C {
	return r.Max.Y - r.Min.Y
}

func (r Rectangle[C]) Overlaps(s Rectangle[C]) bool {
	return r.Min.X < s.Max.X && s.Min.X < r.Max.X &&
		r.Min.Y < s.Max.Y && s.Min.Y < r.Max.Y
}

func (r *Rectangle[C]) SetLeft(x C) {
	w := r.Width()
	r.Min.X = x
	r.Max.X = x + w
}

func (r *Rectangle[C]) SetRight(x C) {
	w := r.Width()
	r.Min.X = x - w
	r.Max.X = x
}

func (r *Rectangle[C]) SetBottom(y C) {
	h := r.Height()
	r.Min.Y = y - h
	r.Max.Y = y
}

func (r *Rectangle[C]) SetTop(y C) {
	h := r.Height()
	r.Min.Y = y
	r.Max.Y = y + h
}

func (r *Rectangle[C]) MoveX(dX C) {
	r.Min.X += dX
	r.Max.X += dX
}

func (r *Rectangle[C]) MoveY(dY C) {
	r.Min.Y += dY
	r.Max.Y += dY
}

func (r *Rectangle[C]) Scale(factor float64) {
	dX := C(0.5 * float64(r.Width()) * (1 - factor))
	dY := C(0.5 * float64(r.Height()) * (1 - factor))
	r.Min.X += dX
	r.Max.X -= dX
	r.Min.Y += 2 * dY
}
