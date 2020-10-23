package main

import (
	"encoding/xml"
	"fmt"
	"github.com/fogleman/gg"
	"io/ioutil"
	"svg/svg"
)

func main() {
	var foo svg.Svg
	svgPath := "examples/test.svg"
	data, _ := ioutil.ReadFile(svgPath)
	xml.Unmarshal(data, &foo)
	for _, g := range foo.Gs {
		for _, text := range g.Texts {
			w, h := calculateDimensions(text.Content, text.FontFamily, text.FontSize)
			println(w, h)
		}
	}

}

func calculateDimensions(text, fontString string, fontSize float64) (float64, float64) {
	ctx := gg.NewContext(0, 0)
	face, err := svg.GetFont(fontString, fontSize)
	if err != nil {
		panic(fmt.Sprintf("font-string could not be parsed: %s", fontString))
	}
	ctx.SetFontFace(*face)
	return ctx.MeasureString(text)
}
