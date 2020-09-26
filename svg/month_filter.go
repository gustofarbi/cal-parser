package svg

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Node struct {
	XMLName  xml.Name
	DataName string `xml:"data-name,attr"`
	Id       string `xml:"id,attr"`
	Content  string `xml:",innerxml"`
	Nodes    []Node `xml:",any"`
	Cont     []byte
}

func walkAndFilter(nodes []Node, mapping map[string]Node, f func(Node) bool) {
	for i, n := range nodes {
		if n.Id != "" {
			mapping[n.Id] = n
		}
		if f(n) {
			walkAndFilter(n.Nodes, mapping, f)
		} else {
			nodes = append(nodes[:i], nodes[i+1:]...)
		}
	}
}

func filterByMonth(data []byte, month int) (Node, map[string]Node) {
	var start Node
	var buf bytes.Buffer
	xml.EscapeText(&buf, data)
	e := xml.NewDecoder(data)
	e.
	s := buf.String()
	println(s)
	if err := xml.Unmarshal(data, &start); err != nil {
		panic("cannot parse svg: " + err.Error())
	}

	mapping := make(map[string]Node)
	walkAndFilter(start.Nodes, mapping, func(node Node) bool {
		reg := regexp.MustCompile("^m(\\d+)") // todo check
		ids := strings.Split(node.Id, "-")
		for _, id := range ids {
			subj := reg.FindString(id)
			n, err := strconv.Atoi(subj)
			if err != nil {
				continue
			}
			if n == month {
				return false
			}
		}
		return true
	})

	return start, mapping
}

func (n *Node) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var str string
	for _, attr := range start.Attr {
		var a string
		if attr.Name.Space == "" {
			a = fmt.Sprintf("%s=\"%s\"", attr.Name.Local, attr.Value)

		} else {
			a = fmt.Sprintf("%s:%s=\"%s\"", attr.Name.Space, attr.Name.Local, attr.Value)
		}
		if str != "" {
			str += " "
		}
		str += a
	}
	n.Cont = make([]byte, 0)
	for {
		tok, err := d.Token()
		if err != nil {
			return nil
		}
		switch c := tok.(type) {
		default:
			println(c)
		case xml.CharData:
			for _, b := range []byte(c) {
				n.Cont = append(n.Cont, b)
			}
		case xml.EndElement:
			//n.Cont += tok.(string)
			return nil
			//default:
			//	n.Cont += tok.(string)
		}
	}
}
