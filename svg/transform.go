package svg

import "math"

type Transforms struct {
	collection []Transformer
}

func (t Transforms) Apply(text CalendarText) {
	for _, transform := range t.collection {
		transform.Apply(text)
	}
}

type Transformer interface {
	Apply(object HasPosition)
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

func (t Translate) Apply(object HasPosition) {
	object.Position().X += t.X
	object.Position().Y += t.Y
}

func (r Rotate) Apply(object HasPosition) {
	object.Position().Rotation += r.AngleDeg
}

func (s Scale) Apply(object HasPosition) {
	object.Position().Width *= s.ScaleX
	object.Position().Height *= s.ScaleY
}

func (m Matrix) Apply(object HasPosition) {
	dx := m.Values[4]
	dy := m.Values[5]
	angle := math.Atan(m.Values[2] / m.Values[0])
	scaleX := m.Values[0] / math.Cos(angle)
	scaleY := m.Values[3] / math.Cos(angle)
	angleDeg := -(angle * 180.0 / math.Pi)

	object.Position().X += dx
	object.Position().Y += dy
	object.Position().Width *= scaleX
	object.Position().Height *= scaleY
	object.Position().Rotation += angleDeg
}
