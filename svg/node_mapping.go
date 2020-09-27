package svg

import (
	"bytes"
	"golang.org/x/net/html"
	"regexp"
	"strconv"
	"sync"
)

type NodeMapping struct {
	root    *Node
	Mapping *sync.Map
	months  []*Node
}

type Node struct {
	Id           string
	Children     []*Node
	start        string
	end          string
	renderMonths []int
}

func NewMapping(data []byte) NodeMapping {
	root := &Node{}
	t := html.NewTokenizer(bytes.NewReader(data))
	m := &sync.Map{}
	nm := NodeMapping{root, m, make([]*Node, 0)}
	nm.walk(t, root)
	return nm
}

func (nm *NodeMapping) walk(t *html.Tokenizer, parent *Node) (end string) {
loop:
	for {
		ttype := t.Next()
		tok := t.Token()
		switch ttype {
		case html.ErrorToken:
			break loop
		case html.SelfClosingTagToken:
			nm.addToMapping(tok, parent)
		case html.StartTagToken:
			n := nm.addToMapping(tok, parent)
			if n != nil {
				n.end = nm.walk(t, n)
			}
			continue
		case html.EndTagToken:
			end = tok.String()
			return end
		}
	}

	return ""
}

var monthRegex = regexp.MustCompile("-m(\\d+)")
func (nm *NodeMapping) addToMapping(t html.Token, parent *Node) *Node {
	id := ""
	for _, attr := range t.Attr {
		if attr.Key == "id" {
			id = attr.Val
			break
		}
	}
	if id == "" {
		return nil
	}
	months := monthRegex.FindStringSubmatch(id)
	monthNumbers := make([]int, 0)
	for _, month := range months {
		n, _:= strconv.Atoi(month)
		monthNumbers = append(monthNumbers, n)
	}
	n := &Node{
		id,
		make([]*Node, 0),
		t.String(),
		"",
		monthNumbers,
	}
	nm.Mapping.Store(id, n)
	if len(monthNumbers) > 0 {

	}
	parent.Children = append(parent.Children, n)
	return n
}

//func filterByMonth(data []byte, month int) (Node, map[string]Node) {
//	var start Node
//	if err := xml.Unmarshal(data, &start); err != nil {
//		panic("cannot parse svg: " + err.Error())
//	}
//
//	mapping := make(map[string]Node)
//	walkAndFilter(start.Nodes, mapping, func(node Node) bool {
//		reg := regexp.MustCompile("^m(\\d+)") // todo check
//		ids := strings.Split(node.Id, "-")
//		for _, id := range ids {
//			subj := reg.FindString(id)
//			n, err := strconv.Atoi(subj)
//			if err != nil {
//				continue
//			}
//			if n == month {
//				return false
//			}
//		}
//		return true
//	})
//
//	return start, mapping
//}
