package svg

import (
	"fmt"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"io/ioutil"
	"math"
	"os"
	"strconv"
	"strings"
)

var (
	availableFonts []string
	fontCache      = make(map[string]*font.Face)
)

func init() {
	wd, _ := os.Getwd()
	fontsPrefix := wd + "/resources/fonts/"
	fileInfos, err := ioutil.ReadDir(fontsPrefix)

	if err != nil {
		fmt.Println("loading fonts: " + err.Error())
	}

	for _, fileInfo := range fileInfos {
		availableFonts = append(availableFonts, fontsPrefix+fileInfo.Name())
	}
}

func GetFont(fonts string, size float64) (face *font.Face, e error) {
	fontFamilies := strings.Split(fonts, ",")
	for _, family := range fontFamilies {
		var ok bool
		family = strings.Trim(family, " ")
		face, ok = getFamily(family, size)
		if ok {
			return face, nil
		} else {
			for _, fontFile := range availableFonts {
				if strings.Contains(fontFile, family) {
					bytes, err := ioutil.ReadFile(fontFile)
					if err != nil {
						panic("cannot load fontFile: " + fontFile)
					}
					f, err := truetype.Parse(bytes)
					if err != nil {
						panic("error loading faceTmp: " + err.Error())
					}
					faceTmp := truetype.NewFace(f, &truetype.Options{Size: size})
					cacheFamily(fontFamilies, size, &faceTmp)
					return &faceTmp, nil
				}
			}
		}
	}

	return nil, fmt.Errorf("no fontfile found for family %s", fonts)
}

func getFamily(identifier string, size float64) (face *font.Face, ok bool) {
	sizeInt := int(math.Round(size))
	face, ok = fontCache[identifier+strconv.Itoa(sizeInt)]
	return face, ok
}

func cacheFamily(identifiers []string, size float64, face *font.Face) {
	for _, identifier := range identifiers {
		sizeInt := int(math.Round(size))
		identifier = strings.TrimSpace(identifier)
		fontCache[identifier+strconv.Itoa(sizeInt)] = face
	}
}
