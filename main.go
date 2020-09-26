package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"svg/svg"
	"time"
)

func main() {
	start := time.Now()
	var foo svg.Svg
	data, err := ioutil.ReadFile("examples/wandkalender_a4-quer_month-1.svg")
	//data, err := ioutil.ReadFile("wandkalender_a4-hoch_month.svg")
	if err != nil {
		fmt.Errorf("shit happened: %s", err)
	}
	err = xml.Unmarshal(data, &foo)
	if err != nil {
		fmt.Errorf("shit happened: %s", err)
	}

	year, m, _ := time.Now().Date()
	month := int(m) + 1
	c := svg.NewCalendar(data, month)
	dims := strings.Split(foo.ViewBox, " ")
	size := 2000.0
	widthViewbox, _ := strconv.ParseFloat(dims[2], 64)
	scalingRatio := size / widthViewbox

	c.Parse(foo, string(data), scalingRatio)
	c.Render(year, month, size)

	fmt.Printf("done in: %vs\n", time.Since(start).Seconds())
	return
}
