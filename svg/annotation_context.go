package svg

import "sort"

type Context struct {
	Annotations map[int]map[string]Annotation
}

func (c Context) Update(h HasAnnotations) {
	newAnnotations := Parse(h)
	for _, anno := range newAnnotations {
		c.Annotations[anno.Priority()][anno.Id()] = anno
	}
}

func (c Context) ApplyEarly(text CalendarText) {
	for prio, annotations := range c.Annotations {
		for id, annotation := range annotations {
			if annotation.Priority() < 0 {
				annotation.Apply(text)
				c.Annotations[prio][id] = nil
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
		for _, anno := range c.Annotations[prio] {
			anno.Apply(text)
		}
	}
}
