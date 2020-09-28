package svg

import (
	"bytes"
	"golang.org/x/net/html"
	"regexp"
	"strconv"
	"sync"
)

var (
	keepMonthsRegex = regexp.MustCompile("-m(\\d+)")
	skipDayRegex    = regexp.MustCompile("-n(\\d+)")
)

type NodeMapping struct {
	root       *Node
	Mapping    *sync.Map
	keepMonths []*Node
	skipDays   []*Node
}

type Node struct {
	Id           string
	Children     []*Node
	start        string
	end          string
	monthsToKeep []int
	daysToSkip   []int
}

func NewMapping(data []byte) NodeMapping {
	root := &Node{}
	t := html.NewTokenizer(bytes.NewReader(data))
	m := &sync.Map{}
	nm := NodeMapping{
		root,
		m,
		make([]*Node, 0),
		make([]*Node, 0),
	}
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

	var keepMonths []int
	var skipDays []int

	months := keepMonthsRegex.FindStringSubmatch(id)
	if len(months) > 0 {
		keepMonths = make([]int, 0)
		for _, month := range months[1:] {
			x, _ := strconv.Atoi(month)
			keepMonths = append(keepMonths, x)
		}
	}

	days := skipDayRegex.FindStringSubmatch(id)
	if len(days) > 0 {
		skipDays = make([]int, 0)
		for _, day := range days[1:] {
			x, _ := strconv.Atoi(day)
			skipDays = append(skipDays, x)
		}
	}

	n := &Node{
		id,
		make([]*Node, 0),
		t.String(),
		"",
		keepMonths,
		skipDays,
	}
	nm.Mapping.Store(id, n)
	if len(keepMonths) > 0 {
		nm.keepMonths = append(nm.keepMonths, n)
	}
	if len(skipDays) > 0 {
		nm.skipDays = append(nm.skipDays, n)
	}
	parent.Children = append(parent.Children, n)
	return n
}
