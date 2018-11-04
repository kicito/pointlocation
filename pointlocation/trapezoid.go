package pointlocation

import (
	"fmt"
	"log"
)

type trapezoid struct {
	leftp       *point
	rightp      *point
	top         *segment
	bottom      *segment
	upperLeftN  *trapezoid
	lowerLeftN  *trapezoid
	upperRightN *trapezoid
	lowerRightN *trapezoid
	dagRef      *node
}

func (t trapezoid) String() string {
	leftpStr := "nil"
	rightpStr := "nil"
	topStr := "nil"
	bottomStr := "nil"
	neighborLStr := "nil"
	neighborRStr := "nil"
	if t.leftp != nil {
		leftpStr = (*t.leftp).String()
	}
	if t.rightp != nil {
		rightpStr = (*t.rightp).String()
	}
	if t.top != nil {
		topStr = (*t.top).String()
	}
	if t.bottom != nil {
		bottomStr = (*t.bottom).String()
	}
	if t.upperLeftN != nil && t.lowerLeftN != nil {
		neighborLStr = "2"
	} else if t.upperLeftN != nil || t.lowerLeftN != nil {
		neighborLStr = "1"
	}
	if t.upperRightN != nil && t.lowerRightN != nil {
		neighborRStr = "2"
	} else if t.upperRightN != nil || t.lowerRightN != nil {
		neighborRStr = "1"
	}

	return fmt.Sprintf("trapezoid lp: %v\nt: %v\nrp: %v \nb: %v\nleft neighbor: %v\nright neighbor: %v",
		leftpStr,
		topStr,
		rightpStr,
		bottomStr,
		neighborLStr,
		neighborRStr,
	)
}

func (t trapezoid) topSegment() segment {
	startY, err := (*t.top).y((*t.leftp).y)
	if err != nil {
		log.Fatal(err)
	}
	startPoint := point{x: (*t.leftp).x, y: startY}

	endY, err := (*t.top).y((*t.rightp).y)
	if err != nil {
		log.Fatal(err)
	}
	endPoint := point{x: (*t.rightp).x, y: endY}
	return newSegment(startPoint, endPoint)
}

func (t trapezoid) bottomSegment() segment {
	startY, err := (*t.bottom).y((*t.leftp).y)
	if err != nil {
		log.Fatal(err)
	}
	startPoint := point{x: (*t.leftp).x, y: startY}

	endY, err := (*t.bottom).y((*t.rightp).y)
	if err != nil {
		log.Fatal(err)
	}
	endPoint := point{x: (*t.rightp).x, y: endY}
	return newSegment(startPoint, endPoint)

}

func (t trapezoid) leftSegment() segment {
	startY, err := (*t.top).y((*t.leftp).y)
	if err != nil {
		log.Fatal(err)
	}
	startPoint := point{x: (*t.leftp).x, y: startY}

	endY, err := (*t.bottom).y((*t.leftp).y)
	if err != nil {
		log.Fatal(err)
	}
	endPoint := point{x: (*t.leftp).x, y: endY}
	return newSegment(startPoint, endPoint)
}

func (t trapezoid) rightSegment() segment {
	startY, err := (*t.top).y((*t.rightp).y)
	if err != nil {
		log.Fatal(err)
	}
	startPoint := point{x: (*t.rightp).x, y: startY}

	endY, err := (*t.bottom).y((*t.rightp).y)
	if err != nil {
		log.Fatal(err)
	}
	endPoint := point{x: (*t.rightp).x, y: endY}
	return newSegment(startPoint, endPoint)
}

// comparing if the point is in side trapezoid by check orientation of point to each segment is inside trapeziod
func (t trapezoid) pointInTrapezoid(p point) (b bool) {
	var tmpSegment segment
	var tmpOrt int
	tmpSegment = t.leftSegment()
	tmpOrt = p.orientationFromSegment(tmpSegment)
	if tmpOrt == counterclockwise {
		b = false
		return
	}
	tmpSegment = t.topSegment()
	tmpOrt = p.orientationFromSegment(tmpSegment)
	if tmpOrt == counterclockwise {
		b = false
		return
	}
	tmpSegment = t.rightSegment()
	tmpOrt = p.orientationFromSegment(tmpSegment)
	if tmpOrt == clockwise {
		b = false
		return
	}

	tmpSegment = t.bottomSegment()
	tmpOrt = p.orientationFromSegment(tmpSegment)
	if tmpOrt == clockwise {
		b = false
		return
	}

	b = true
	return
}

func (t trapezoid) segmentInTrapezoid(s segment) (b bool) {
	b = t.pointInTrapezoid(s.startPoint) && t.pointInTrapezoid(s.endPoint)
	return
}

func (t trapezoid) lineIntersectOrInside(s segment) (b bool) {

	// if segment is in trapezoid
	if t.segmentInTrapezoid(s) {
		b = true
		return
	}

	var tmpSegment segment
	// if segment is intersect trapezoid
	tmpSegment = t.leftSegment()
	if tmpSegment.isSegmentIntersect(s) {
		b = true
		return
	}
	tmpSegment = t.topSegment()
	if tmpSegment.isSegmentIntersect(s) {
		b = true
		return
	}
	tmpSegment = t.rightSegment()
	if tmpSegment.isSegmentIntersect(s) {
		b = true
		return
	}
	tmpSegment = t.bottomSegment()
	if tmpSegment.isSegmentIntersect(s) {
		b = true
		return
	}
	return
}

func (t trapezoid) addSegmentInside(s segment) (ts []trapezoid) {
	ts = make([]trapezoid, 0)
	//|---------------|
	//|	lt|__ut__|rt  |
	//|	  |  dt  |	  |
	//|---------------|
	lt := trapezoid{
		leftp:  t.leftp,
		rightp: &s.startPoint,
		top:    t.top,
		bottom: t.bottom,
	}
	ut := trapezoid{
		leftp:  &s.startPoint,
		rightp: &s.endPoint,
		top:    t.top,
		bottom: &s,
	}
	rt := trapezoid{
		leftp:  &s.endPoint,
		rightp: t.rightp,
		top:    t.top,
		bottom: t.bottom,
	}
	bt := trapezoid{
		leftp:  &s.startPoint,
		rightp: &s.endPoint,
		top:    &s,
		bottom: t.bottom,
	}
	// assign neighbor
	lt.upperRightN = &ut
	lt.lowerRightN = &bt

	ut.upperLeftN = &lt
	ut.upperRightN = &rt

	rt.upperLeftN = &ut
	rt.lowerLeftN = &bt

	bt.lowerLeftN = &lt
	bt.lowerRightN = &rt

	ts = append(ts, lt)
	ts = append(ts, ut)
	ts = append(ts, rt)
	ts = append(ts, bt)
	return
}

func (t trapezoid) addSegment(s segment) (ts []trapezoid) {
	// case segment inside
	if t.segmentInTrapezoid(s) {
		ts = t.addSegmentInside(s)
		return
	}

	return
}

func boundingBox(ss []segment) (tr trapezoid) {
	boundingBoxTop := point{x: ss[0].maxX(), y: ss[0].maxY()}
	boundingBoxBot := point{x: ss[0].minX(), y: ss[0].minY()}
	index := 1
	for _ = range ss[1:] {
		if ss[index].maxX() > boundingBoxTop.x {
			boundingBoxTop.x = ss[index].maxX()
		}
		if ss[index].maxY() > boundingBoxTop.y {
			boundingBoxTop.y = ss[index].maxY()
		}
		if ss[index].minX() > boundingBoxBot.x {
			boundingBoxBot.x = ss[index].minX()
		}
		if ss[index].minY() > boundingBoxBot.y {
			boundingBoxBot.y = ss[index].minY()
		}
	}
	boundingBoxTop.x += 5
	boundingBoxTop.y += 5
	boundingBoxBot.x -= 5
	boundingBoxBot.y -= 5

	bounderyTopSegment := newSegment(
		point{x: boundingBoxBot.x, y: boundingBoxTop.y},
		boundingBoxTop,
	)

	bounderyBotSegment := newSegment(
		boundingBoxBot,
		point{x: boundingBoxTop.x, y: boundingBoxBot.y},
	)

	tr = trapezoid{
		leftp:  &boundingBoxBot,
		rightp: &point{x: boundingBoxTop.x, y: boundingBoxBot.y},
		top:    &bounderyTopSegment,
		bottom: &bounderyBotSegment,
	}
	return
}

// implementation of node interface
func (trapezoid) check(i point) node {
	return nil
}

func (*trapezoid) leftChild() node {
	return nil
}

func (*trapezoid) rightChild() node {
	return nil
}
