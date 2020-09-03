package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"svg/svg"
)

func main() {
	var foo svg.Svg
	data, err := ioutil.ReadFile("wandkalender_a3-hoch_month.svg")
	if err != nil {
		fmt.Errorf("shit happened: %s", err)
	}
	err = xml.Unmarshal(data, &foo)
	if err != nil {
		fmt.Errorf("shit happened: %s", err)
	}
	fmt.Printf("%#v\n", foo)
}
