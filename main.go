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
	//data, err := ioutil.ReadFile("examples/wandkalender_a4-quer_month-1.svg")
	data, err := ioutil.ReadFile("examples/wandkalender_a4-hoch_month-2.svg")
	//data, err := ioutil.ReadFile("examples/wandkalender_a3-hoch_month.svg")
	if err != nil {
		fmt.Errorf("shit happened: %s", err)
	}
	err = xml.Unmarshal(data, &foo)
	if err != nil {
		fmt.Errorf("shit happened: %s", err)
	}
	// todo: unite svg and node-mapping

	year, m, _ := time.Now().Date()
	month := int(m) - 8
	year += 1
	c := svg.NewCalendar(data)
	dims := strings.Split(foo.ViewBox, " ")
	size := 2000.0
	widthViewbox, _ := strconv.ParseFloat(dims[2], 64)
	heightViewbox, _ := strconv.ParseFloat(dims[3], 64)
	scalingRatio := size / widthViewbox

	c.Parse(foo, string(data), scalingRatio)
	c.Render(year, month, size, size*(heightViewbox/widthViewbox))

	fmt.Printf("done in: %vs\n", time.Since(start).Seconds())
	return
}
