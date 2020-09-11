package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"svg/svg"
	"testing"
	"time"
)

func TestParseRender(t *testing.T) {
	var foo svg.Svg
	data, err := ioutil.ReadFile("wandkalender_a4-hoch_month.svg")
	if err != nil {
		fmt.Errorf("shit happened: %s", err)
	}
	err = xml.Unmarshal(data, &foo)
	if err != nil {
		fmt.Errorf("shit happened: %s", err)
	}

	c := svg.NewCalendar()
	dims := strings.Split(foo.ViewBox, " ")
	size := 2000.0
	widthViewbox, _ := strconv.ParseFloat(dims[2], 64)
	scalingRatio := size / widthViewbox
	c.Parse(foo, string(data), scalingRatio)
	year, month, _ := time.Now().Date()
	c.Render(year, int(month), size)
}