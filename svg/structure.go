package svg

import (
	"encoding/xml"
	"regexp"
	"strconv"
	"strings"
)

type Svg struct {
	Gs      []Group `xml:"g"`
	ViewBox string  `xml:"viewBox,attr"`
	Defs    Def     `xml:"defs"`
	Title   string  `xml:"title"`
}

type Group struct {
	Common
	Content           string    `xml:",innerxml"`
	Id                string    `xml:"id,attr,omitempty"`
	ClipPathReference string    `xml:"clip-path,attr"`
	Rects             []Rect    `xml:"rect"`
	Circles           []Circle  `xml:"circle"`
	Gs                []Group   `xml:"g"`
	Images            []Image   `xml:"image,omitempty"`
	Texts             []Text    `xml:"text"`
	Paths             []Path    `xml:"path"`
	Polygons          []Polygon `xml:"polygon"`
}

type Image struct {
	Common
	Position
	Transform Transforms `xml:"transform,attr"`
	Href      string     `xml:"http://www.w3.org/1999/xlink href,attr"`
}

type Text struct {
	Common
	Position
	Id         string     `xml:"id,attr"`
	FontSize   float64    `xml:"font-size,attr"`
	Fill       string     `xml:"fill,attr"`
	FontFamily string     `xml:"font-family,attr"`
	FontWeight int        `xml:"font-weight,attr,omitempty"`
	Tranform   Transforms `xml:"tranform,attr"`
	Content    string     `xml:",chardata"`
}

type Rect struct {
	Common
	Position
	Content string `xml:",innerxml"`
	Id      string `xml:"id,attr,omitempty"`
	Fill    string `xml:"fill,attr"`
}

type Def struct {
	ClipPaths []ClipPath `xml:"clipPath"`
}

type Polygon struct {
	Common
	PointsSlice Points `xml:"points,attr"`
	Fill        string `xml:"fill,attr"`
}

type Points struct {
	Ps []float64 // todo: is this necessary?
}

type Common struct {
	Raw string `xml:",innerxml"`
	Values []string `xml:"data-name,attr"`
}

type Circle struct {
	Common
	Cx   float64 `xml:"cx,attr"`
	Cy   float64 `xml:"cy,attr"`
	R    float64 `xml:"r,attr"`
	Fill string  `xml:"fill,attr"`
}

type Path struct {
	Common
	D       string  `xml:"d,attr"`
	Fill    string  `xml:"fill,attr"`
	Opacity float64 `xml:"opacity,attr,omitempty"`
}

type ClipPath struct {
	Id      string   `xml:"id"`
	Rects   []Rect   `xml:"rect"`
	Circles []Circle `xml:"circle"`
	Paths   []Path   `xml:"path"`
}

func (t *Transforms) UnmarshalXMLAttr(attr xml.Attr) error {
	reg := regexp.MustCompile("(\\w+\\(.*?\\))")
	transforms := reg.FindAllString(attr.Value, -1)
	for _, s := range transforms {
		switch {
		case strings.Contains(s, "rotate"):
			reg := regexp.MustCompile("\\((-?[0-9.]+)\\)")
			val := reg.FindAllStringSubmatch(s, -1)
			angleDeg, _ := strconv.ParseFloat(strings.TrimSpace(val[0][1]), 64)
			t.collection = append(t.collection, Rotate{AngleDeg: angleDeg})
			break
		case strings.Contains(s, "translate"):
			reg := regexp.MustCompile("([0-9.-]*)[\\s)]+")
			val := reg.FindAllStringSubmatch(s, -1)
			x, _ := strconv.ParseFloat(strings.TrimSpace(val[0][1]), 64)
			y, _ := strconv.ParseFloat(strings.TrimSpace(val[1][1]), 64)
			t.collection = append(t.collection, Translate{X: x, Y: y})
			break
		case strings.Contains(s, "matrix"):
			reg := regexp.MustCompile("\\(([-+.\\d\\s,]+)\\)")
			val := reg.FindAllStringSubmatch(s, -1)
			matrix := make(map[int]float64, len(val))
			for i, value := range val[0] {
				parsedValue, _ := strconv.ParseFloat(strings.TrimSpace(value), 64)
				matrix[i] = parsedValue
			}
			t.collection = append(t.collection, Matrix{Values: matrix})
			break
		case strings.Contains(s, "scale"):
			reg := regexp.MustCompile("\\(([-+.\\d]+)([\\s,][-+.\\d]+)?\\)")
			val := reg.FindAllStringSubmatch(s, -1)
			x, _ := strconv.ParseFloat(strings.TrimSpace(val[0][1]), 64)
			var y float64
			if len(val) == 2 {
				y, _ = strconv.ParseFloat(strings.TrimSpace(val[1][1]), 64)
			} else {
				y = x
			}
			t.collection = append(t.collection, Scale{x, y})
			break
		}
	}

	return nil
}

func (p *Points) UnmarshalXMLAttr(attr xml.Attr) error {
	points := strings.Split(attr.Value, " ")
	p.Ps = make([]float64, len(points))
	for i, point := range points {
		p.Ps[i], _ = strconv.ParseFloat(point, 64)
	}

	return nil
}

func (a *Common) UnmarshalXMLAttr(attr xml.Attr) error {
	a.Values = strings.Split(attr.Value, "-")

	return nil
}

type HasAnnotations interface {
	All() []string
}

func (a Common) All() []string {
	return a.Values
}

type Position struct {
	X        float64 `xml:"x,attr"`
	Y        float64 `xml:"y,attr"`
	Width    float64 `xml:"width,attr"`
	Height   float64 `xml:"height,attr"`
	Rotation float64 `xml:"rotation,attr"`
}

func (p Position) Position() *Position {
	return &p
}

type HasPosition interface {
	Position() *Position
}
