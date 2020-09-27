package svg

import (
	"bytes"
	"golang.org/x/net/html"
	"sync"
)

type NodeMapping struct {
	root    *Node
	Mapping *sync.Map
}

type Node struct {
	Id       string
	Children []*Node
	start    string
	end      string
}

func NewMapping(data []byte) NodeMapping {
	root := &Node{}
	t := html.NewTokenizer(bytes.NewReader(data))
	m := &sync.Map{}
	walk(t, root, m)
	return NodeMapping{root, m}
}

func walk(t *html.Tokenizer, parent *Node, m *sync.Map) (end string) {
loop:
	for {
		ttype := t.Next()
		tok := t.Token()
		switch ttype {
		case html.ErrorToken:
			break loop
		case html.SelfClosingTagToken:
			addToMapping(tok, parent, m)
		case html.StartTagToken:
			n := addToMapping(tok, parent, m)
			if n != nil {
				n.end = walk(t, n, m)
			}
			continue
		case html.EndTagToken:
			end = tok.String()
			return end
		}
	}

	return ""
}

func addToMapping(t html.Token, parent *Node, m *sync.Map) *Node {
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
	n := &Node{
		id,
		make([]*Node, 0),
		t.String(),
		"",
	}
	m.Store(id, n)
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
