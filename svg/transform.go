package svg

import "math"

type Transforms struct {
	collection []Transformer
}

func (t Transforms) Apply(text *CalendarText) {
	for _, transform := range t.collection {
		transform.Apply(&text.Position)
	}
}

type Transformer interface {
	Apply(object *Position)
}

type Translate struct {
	X, Y float64
}

type Rotate struct {
	AngleDeg float64
}

type Scale struct {
	ScaleX, ScaleY float64
}

type Matrix struct {
	Values map[int]float64
}

func (t Translate) Apply(object *Position) {
	object.X += t.X
	object.Y += t.Y
}

func (r Rotate) Apply(object *Position) {
	object.Rotation += r.AngleDeg
}

func (s Scale) Apply(object *Position) {
	object.Width *= s.ScaleX
	object.Height *= s.ScaleY
}

func (m Matrix) Apply(object *Position) {
	dx := m.Values[4]
	dy := m.Values[5]
	angle := math.Atan(m.Values[2] / m.Values[0])
	scaleX := m.Values[0] / math.Cos(angle)
	scaleY := m.Values[3] / math.Cos(angle)
	angleDeg := -(angle * 180.0 / math.Pi)

	object.X += dx
	object.Y += dy
	object.Width *= scaleX
	object.Height *= scaleY
	object.Rotation += angleDeg
}
