package svg

import (
	"sort"
	"sync"
)

type AnnotationCollection map[int]map[string][]Annotation

type Context struct {
	Receiver     chan interface{}
	ReceiverWg   *sync.WaitGroup
	Annotations  AnnotationCollection
	RenderMonths []RenderMonthOnly
}

func NewContext(ch chan interface{}, wg *sync.WaitGroup) Context {
	return Context{
		ch,
		wg,
		make(map[int]map[string][]Annotation),
		make([]RenderMonthOnly, 0),
	}
}

var ctxLock sync.Mutex

func (c Context) Merge(h HasAnnotations) (result Context) {
	ctxLock.Lock()
	result.Receiver = c.Receiver
	result.ReceiverWg = c.ReceiverWg
	result.Annotations = make(map[int]map[string][]Annotation)
	for prio, group := range c.Annotations {
		result.Annotations[prio] = make(map[string][]Annotation)
		for id, collection := range group {
			result.Annotations[prio][id] = make([]Annotation, len(collection))
			copy(result.Annotations[prio][id], collection)
		}
	}
	result.Add(Parse(h))
	ctxLock.Unlock()

	return result
}

func (c Context) Add(as []Annotation) {
	for _, a := range as {
		if a != nil {
			if _, ok := c.Annotations[a.Priority()]; !ok {
				c.Annotations[a.Priority()] = make(map[string][]Annotation)
			}
			if c.Annotations[a.Priority()][a.Id()] == nil {
				c.Annotations[a.Priority()][a.Id()] = make([]Annotation, 0)
			}
			c.Annotations[a.Priority()][a.Id()] = append(c.Annotations[a.Priority()][a.Id()], a)
		}
	}
}

func (c Context) Get(prio int, id string) []Annotation {
	if group, ok := c.Annotations[prio][id]; ok {
		return group
	}

	return make([]Annotation, 0)
}

func (a AnnotationCollection) ApplyEarly(text *CalendarText) {
	for prio, annotations := range a {
		if prio < 0 {
			for id, group := range annotations {
				for pos, single := range group {
					single.Apply(text)
					a[prio][id][pos] = nil
				}
			}
		}
	}
}

func (a AnnotationCollection) ApplyLate(text *CalendarText) {
	order := make([]int, len(a))
	for prio := range a {
		order = append(order, prio)
	}
	sort.Ints(order)
	for prio := range order {
		for _, group := range a[prio] {
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

func (c Context) HandleSpecialAnnotation(annotations []Annotation, rawSvg string) {
	for _, annotation := range annotations {
		for _, single := range c.Get(annotation.Priority(), annotation.Id()) {
			switch x := single.(type) {
			case LineWeekendElement:
			case LineWeekdayElement: // todo check
				weekdayPosition := WeekdayPosition{}
				pos := c.Get(weekdayPosition.Priority(), weekdayPosition.Id())
				if len(pos) == 0 {
					continue
				}
				x.Attribute.Val = pos[0].Attr().(int)
			}
			annotationObject := AnnotationObject{single, rawSvg}
			c.Receiver <- annotationObject
		}
	}
}

type AnnotationObject struct {
	Annotation
	RawSvg string
}
