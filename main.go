package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"svg/svg"
)

func main() {
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
	c.Parse(foo, string(data), 1.5)
	c.Render()

	return
}
