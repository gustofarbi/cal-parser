package svg

type Attribute struct {
	Val interface{}
}

func (a Attribute) Attr() interface{} {
	return a.Val
}

func (a *Attribute) Set(o interface{}) {
	a.Val = o
}
