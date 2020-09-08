package svg

import (
	"fmt"
	"github.com/fogleman/gg"
	"image"
	"io/ioutil"
	"os"
	"os/exec"
)

func (c *Calendar) Render() {
	im := renderSvg(c.svgContent)
	ctx := gg.NewContextForImage(*im)
	err := ctx.SavePNG("saved.png")
	if err != nil {
		fmt.Println(err)
	}
}

func (c *Calendar) drawTexts(canvas *image.Image) {

}

func (c *Calendar) fillTable() {

}

func renderSvg(svg string) *image.Image {
	svgFile, _ := ioutil.TempFile("", "*.svg")
	defer os.Remove(svgFile.Name())
	pngFile, _ := ioutil.TempFile("", "*.png")
	defer os.Remove(pngFile.Name())

	_, err := svgFile.WriteString(svg)
	if err != nil {
		fmt.Println(err)
	}

	cmd := exec.Command("rsvg-convert", "-w", "2000", svgFile.Name(), "-o", pngFile.Name())
	err = cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
	im, _, err := image.Decode(pngFile)
	if err != nil {
		fmt.Println(err)
	}

	return &im
}
