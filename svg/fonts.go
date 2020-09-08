package svg

import (
	"fmt"
	"io/ioutil"
	"strings"
)

var availableFonts []string

func init() {
	fileInfos, err := ioutil.ReadDir("fonts/")

	if err != nil {
		fmt.Println("loading fonts: " + err.Error())
	}

	for _, fileInfo := range fileInfos {
		availableFonts = append(availableFonts, fileInfo.Name())
	}
}

func GetFont(family string) (string, error) {
	for _, font := range availableFonts {
		if strings.Contains(font, family) {
			return "fonts/" + font, nil
		}
	}

	return "", fmt.Errorf("could not load font: %s", family)
}
