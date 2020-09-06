package svg

import (
	"sort"
)

type AnnotationCollection map[int]map[string][]Annotation

type Context struct {
	Receiver     chan interface{}
	Annotations  AnnotationCollection
	RenderMonths []RenderMonthOnly
}

func (c Context) Update(h HasAnnotations) {
	c.Add(Parse(h))
}

func (c Context) Add(as []Annotation) {
	for _, a := range as {
		c.Annotations[a.Priority()][a.Id()] = append(c.Annotations[a.Priority()][a.Id()], a)
	}
}

func (c Context) Get(prio int, id string) []Annotation {
	if group, ok := c.Annotations[prio][id]; ok {
		return group
	}

	return make([]Annotation, 0)
}

func (c Context) ApplyEarly(text CalendarText) {
	for prio, annotations := range c.Annotations {
		for id, group := range annotations {
			for pos, single := range group {
				if single.Priority() < 0 {
					single.Apply(text)
					c.Annotations[prio][id] = append(c.Annotations[prio][id][:pos], c.Annotations[prio][id][pos+1:]...)
				}
			}
		}
	}
}

func (c Context) ApplyLate(text CalendarText) {
	order := make([]int, len(c.Annotations))
	for prio := range c.Annotations {
		order = append(order, prio)
	}
	sort.Ints(order)
	for prio := range order {
		for _, group := range c.Annotations[prio] {
			for _, single := range group {
				single.Apply(text)
			}
		}
	}
}

func (c Context) RenderPrevNext() bool {
	tmp := RenderPrevNextMonth{}
	prio := tmp.Priority()
	id := tmp.Id()
	if group, ok := c.Annotations[prio][id]; ok {
		if len(group) != 0 {
			return group[0].Attr().(bool)
		}
	}
	return false
}

func (c Context) HandleSpecialAnnotation(anno Annotation) {
	for _, single := range c.Get(anno.Priority(), anno.Id()) {
		c.Receiver <- single
	}
}
