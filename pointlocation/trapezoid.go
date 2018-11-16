package pointlocation

import (
	"fmt"
	"log"
)

type trapezoid struct {
	leftp       Point
	rightp      Point
	top         Segment
	bottom      Segment
	upperLeftN  *trapezoid
	lowerLeftN  *trapezoid
	upperRightN *trapezoid
	lowerRightN *trapezoid
	dagRef      node
}

func (t trapezoid) String() string {
	neighborLStr := "nil"
	neighborRStr := "nil"
	dagRefStr := "no reference to dag"
	if t.upperLeftN != nil && t.lowerLeftN != nil {
		neighborLStr = fmt.Sprintf(`left neighbors: 2
	lowerLeft.leftp: %v
	lowerLeft.top: %v
	lowerLeft.rightp: %v
	lowerLeft.bottom: %v
	upperLeft.leftp: %v
	upperLeft.top: %v
	upperLeft.rightp: %v
	upperLeft.bottom: %v`,
			t.lowerLeftN.leftp,
			t.lowerLeftN.top,
			t.lowerLeftN.rightp,
			t.lowerLeftN.bottom,
			t.upperLeftN.leftp,
			t.upperLeftN.top,
			t.upperLeftN.rightp,
			t.upperLeftN.bottom)
	} else if t.upperLeftN != nil {
		neighborLStr = fmt.Sprintf(`left neighbors: 1
	upperLeft.leftp: %v
	upperLeft.top: %v
	upperLeft.rightp: %v
	upperLeft.bottom: %v`,
			t.upperLeftN.leftp,
			t.upperLeftN.top,
			t.upperLeftN.rightp,
			t.upperLeftN.bottom)
	} else if t.lowerLeftN != nil {
		neighborLStr = fmt.Sprintf(`left neighbors: 1
	lowerLeft.leftp: %v
	lowerLeft.top: %v
	lowerLeft.rightp: %v
	lowerLeft.bottom: %v`,
			t.lowerLeftN.leftp,
			t.lowerLeftN.top,
			t.lowerLeftN.rightp,
			t.lowerLeftN.bottom)

	}
	if t.upperRightN != nil && t.lowerRightN != nil {
		neighborRStr = fmt.Sprintf(`right neighbors: 2
	lowerRight.leftp: %v
	lowerRight.top: %v
	lowerRight.rightp: %v
	lowerRight.bottom: %v
	upperRight.leftp: %v
	upperRight.top: %v
	upperRight.rightp: %v
	upperRight.bottom: %v`,
			t.lowerRightN.leftp,
			t.lowerRightN.top,
			t.lowerRightN.rightp,
			t.lowerRightN.bottom,
			t.upperRightN.leftp,
			t.upperRightN.top,
			t.upperRightN.rightp,
			t.upperRightN.bottom)
	} else if t.upperRightN != nil {
		neighborRStr = fmt.Sprintf(`right neighbors: 1
	upperRight.leftp: %v
	upperRight.top: %v
	upperRight.rightp: %v
	upperRight.bottom: %v`,
			t.upperRightN.leftp,
			t.upperRightN.top,
			t.upperRightN.rightp,
			t.upperRightN.bottom)
	} else if t.lowerRightN != nil {
		neighborRStr = fmt.Sprintf(`right neighbors: 1
	lowerRight.leftp: %v
	lowerRight.top: %v
	lowerRight.rightp: %v
	lowerRight.bottom: %v`,
			t.lowerRightN.leftp,
			t.lowerRightN.top,
			t.lowerRightN.rightp,
			t.lowerRightN.bottom)
	}

	if t.dagRef != nil {
		dagRefStr = "has reference to dag"
	}

	return fmt.Sprintf(`trapezoid lp: %v
	t: %v
	rp: %v 
	b: %v
	%v
	%v
	%v`,
		t.leftp,
		t.top,
		t.rightp,
		t.bottom,
		neighborLStr,
		neighborRStr,
		dagRefStr,
	)
}

func (t trapezoid) topSegment() Segment {
	startY, err := t.top.y(t.leftp.x)
	if err != nil {
		log.Fatal(err)
	}
	startPoint := Point{x: t.leftp.x, y: startY}

	endY, err := t.top.y(t.rightp.x)
	if err != nil {
		log.Fatal(err)
	}
	endPoint := Point{x: t.rightp.x, y: endY}
	return NewSegment(startPoint, endPoint)
}

func (t trapezoid) bottomSegment() Segment {
	startY, err := t.bottom.y(t.leftp.x)
	if err != nil {
		log.Fatal(err)
	}
	startPoint := Point{x: t.leftp.x, y: startY}

	endY, err := t.bottom.y(t.rightp.x)
	if err != nil {
		log.Fatal(err)
	}
	endPoint := Point{x: t.rightp.x, y: endY}
	return NewSegment(startPoint, endPoint)

}

func (t trapezoid) leftSegment() Segment {
	startY, err := t.top.y(t.leftp.x)
	if err != nil {
		log.Fatal(err)
	}
	startPoint := Point{x: t.leftp.x, y: startY}

	endY, err := t.bottom.y(t.leftp.x)
	if err != nil {
		log.Fatal(err)
	}
	endPoint := Point{x: t.leftp.x, y: endY}
	return NewSegment(startPoint, endPoint)
}

func (t trapezoid) rightSegment() Segment {
	startY, err := t.top.y(t.rightp.x)
	if err != nil {
		log.Fatal(err)
	}
	startPoint := Point{x: t.rightp.x, y: startY}

	endY, err := t.bottom.y(t.rightp.x)
	if err != nil {
		log.Fatal(err)
	}
	endPoint := Point{x: t.rightp.x, y: endY}
	return NewSegment(startPoint, endPoint)
}

// comparing if the Point is in side trapezoid by check orientation of Point to each Segment is inside trapeziod
func (t trapezoid) pointInTrapezoid(p Point) (b bool) {
	var tmpSegment Segment
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

func (t trapezoid) segmentInTrapezoid(s Segment) (b bool) {
	b = t.pointInTrapezoid(s.startPoint) && t.pointInTrapezoid(s.endPoint)
	return
}

func (t trapezoid) lineIntersectOrInside(s Segment) (b bool) {

	// if Segment is in trapezoid
	if t.segmentInTrapezoid(s) {
		b = true
		return
	}

	var tmpSegment Segment
	// if Segment is intersect trapezoid
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

func (t *trapezoid) assignLeftNeighborToTrapezoid(tr *trapezoid) (err error) {
	// if t.leftp != tr.rightp {
	// 	return fmt.Errorf("trapezoidal %v cannot be left Neighbor to %v", tr, t)
	// }
	if t.top == tr.top {
		t.upperLeftN = tr
		tr.upperRightN = t
	}

	if t.bottom == tr.bottom {
		t.lowerLeftN = tr
		tr.lowerRightN = t
	}

	return
}

func (t *trapezoid) assignRightNeighborToTrapezoid(tr *trapezoid) (err error) {
	// if t.rightp != tr.leftp {
	// 	return fmt.Errorf("trapezoidal %v cannot be right Neighbor to %v", tr, t)
	// }
	if t.top == tr.top {
		tr.upperLeftN = t
		t.upperRightN = tr
	}
	if t.bottom == tr.bottom {
		tr.lowerLeftN = t
		t.lowerRightN = tr
	}
	return
}

func (t *trapezoid) replaceUpperLeftNeighborWith(tr *trapezoid) (err error) {
	if t.upperLeftN != nil {
		if (t.upperLeftN.upperRightN != nil && *t.upperLeftN.upperRightN == *t) ||
			(t.upperLeftN.lowerRightN != nil && *t.upperLeftN.lowerRightN == *t) {
			return t.upperLeftN.assignRightNeighborToTrapezoid(tr)
		}
	}
	return
}

func (t *trapezoid) replaceLowerLeftNeighborWith(tr *trapezoid) (err error) {
	if t.lowerLeftN != nil {
		if (t.lowerLeftN.upperRightN != nil && *t.lowerLeftN.upperRightN == *t) ||
			(t.lowerLeftN.lowerRightN != nil && *t.lowerLeftN.lowerRightN == *t) {
			return t.lowerLeftN.assignRightNeighborToTrapezoid(tr)
		}
	}
	return
}

func (t *trapezoid) replaceUpperRightNeighborWith(tr *trapezoid) (err error) {
	if t.upperRightN != nil {
		if (t.upperRightN.upperLeftN != nil && *t.upperRightN.upperLeftN == *t) ||
			(t.upperRightN.lowerLeftN != nil && *t.upperRightN.lowerLeftN == *t) {
			return t.upperRightN.assignLeftNeighborToTrapezoid(tr)
		}
	}
	return
}

func (t *trapezoid) replaceLowerRightNeighborWith(tr *trapezoid) (err error) {
	if t.lowerRightN != nil {
		if (t.lowerRightN.upperLeftN != nil && *t.lowerRightN.upperLeftN == *t) ||
			(t.lowerRightN.lowerLeftN != nil && *t.lowerRightN.lowerLeftN == *t) {
			return t.lowerRightN.assignLeftNeighborToTrapezoid(tr)
		}
	}
	return
}

func (t *trapezoid) replaceLeftNeighborsWith(tr *trapezoid) (err error) {
	if err = t.replaceUpperLeftNeighborWith(tr); err != nil {
		return
	}
	if err = t.replaceLowerLeftNeighborWith(tr); err != nil {
		return
	}
	return
}
func (t *trapezoid) replaceRightNeighborsWith(tr *trapezoid) (err error) {
	if err = t.replaceUpperRightNeighborWith(tr); err != nil {
		return
	}
	if err = t.replaceLowerRightNeighborWith(tr); err != nil {
		return
	}
	return
}

func (t *trapezoid) replaceUpperNeighborsWith(tr *trapezoid) (err error) {
	if err = t.replaceUpperRightNeighborWith(tr); err != nil {
		return
	}
	if err = t.replaceUpperLeftNeighborWith(tr); err != nil {
		return
	}
	return
}

func (t *trapezoid) replaceLowerNeighborsWith(tr *trapezoid) (err error) {
	if err = t.replaceLowerRightNeighborWith(tr); err != nil {
		return
	}
	if err = t.replaceLowerLeftNeighborWith(tr); err != nil {
		return
	}
	return
}

func (t *trapezoid) addSegmentLeftTrapeziod(s Segment) (lt *trapezoid, dagNode node, err error) {
	lt = &trapezoid{
		leftp:  t.leftp,
		rightp: s.startPoint,
		top:    t.top,
		bottom: t.bottom,
	}
	if err = t.replaceLeftNeighborsWith(lt); err != nil {
		return
	}
	x1 := xNode{xCoordinate: s.startPoint.x}
	ltNode := trapezoidNode{tr: lt, parents: make([]node, 0)}
	lt.dagRef = &ltNode
	x1.assignLeft(&ltNode)
	dagNode = &x1
	return
}

func (t *trapezoid) addSegmentRightTrapezoid(s Segment) (rt *trapezoid, dagNode node, err error) {
	rt = &trapezoid{
		leftp:  s.endPoint,
		rightp: t.rightp,
		top:    t.top,
		bottom: t.bottom,
	}
	if err = t.replaceRightNeighborsWith(rt); err != nil {
		return
	}
	rtNode := trapezoidNode{tr: rt, parents: make([]node, 0)}
	rt.dagRef = &rtNode
	x1 := xNode{xCoordinate: s.endPoint.x}
	x1.assignRight(&rtNode)
	dagNode = node(&x1)
	return
}

func (t *trapezoid) addSegmentInside(oldUT, oldBT *trapezoid, s Segment) (ut, bt *trapezoid, dagNode node, trs []trapezoid, err error) {
	trs = make([]trapezoid, 0)
	var lt, rt *trapezoid
	var xNodeStart, xNodeEnd, yNode node
	if ut, bt, yNode, err = t.splitY(oldUT, oldBT, s); err != nil {
		return
	}
	isFirstTr := t.leftp.x <= s.startPoint.x && t.rightp.x >= s.startPoint.x
	isLastTr := t.leftp.x <= s.endPoint.x && t.rightp.x >= s.endPoint.x

	// if left p is not on same coordinate as segment's endpoint, we do create lt
	if isFirstTr {
		if !t.leftp.sameCoordinate(s.startPoint) {
			lt, xNodeStart, err = t.addSegmentLeftTrapeziod(s)
			if err != nil {
				return
			}
			if err = t.replaceLeftNeighborsWith(lt); err != nil {
				return
			}
			if err = lt.assignRightNeighborToTrapezoid(ut); err != nil {
				return
			}
			if err = lt.assignRightNeighborToTrapezoid(bt); err != nil {
				return
			}
		} else {
			// leftp and segment's start endpoint is same we check how it share the point
			if t.leftp.sameCoordinate(t.top.startPoint) { // connect by segment is lower
				if err = t.replaceLeftNeighborsWith(bt); err != nil {
					return
				}
			} else if t.leftp.sameCoordinate(t.bottom.startPoint) { // connect by segment is higher
				if err = t.replaceLeftNeighborsWith(ut); err != nil {
					return
				}
			} else { // connect by continuation
				if err = t.replaceLeftNeighborsWith(ut); err != nil {
					return
				}
				if err = t.replaceLeftNeighborsWith(bt); err != nil {
					return
				}
			}
		}

	}

	if isLastTr {
		// if right p is not on same coordinate as segment's endpoint, we do create rt
		// else, we do nothing because splitY already done repalcing neighbor
		if !t.rightp.sameCoordinate(s.endPoint) {
			rt, xNodeEnd, err = t.addSegmentRightTrapezoid(s)
			if err != nil {
				return
			}
			if err = ut.replaceRightNeighborsWith(rt); err != nil {
				return
			}
			if err = bt.replaceRightNeighborsWith(rt); err != nil {
				return
			}
			if err = rt.assignLeftNeighborToTrapezoid(ut); err != nil {
				return
			}
			if err = rt.assignLeftNeighborToTrapezoid(bt); err != nil {
				return
			}
		}
	}

	// creating sub tree
	if xNodeStart != nil && xNodeEnd != nil {
		xNodeStart.assignRight(xNodeEnd)
		xNodeEnd.assignLeft(yNode)
		dagNode = node(xNodeStart)
	} else if xNodeStart != nil {
		xNodeStart.assignRight(yNode)
		dagNode = node(xNodeStart)
	} else if xNodeEnd != nil {
		xNodeEnd.assignLeft(yNode)
		dagNode = node(xNodeEnd)
	} else {
		dagNode = node(yNode)
	}

	// add result to trapeziod list
	if lt != nil {
		trs = append(trs, *lt)
	}
	if rt != nil {
		trs = append(trs, *rt)
	}
	trs = append(trs, *ut, *bt)
	return
}

func (t trapezoid) addSegment(oldUT, oldBT *trapezoid, s Segment) (newUT, newBT *trapezoid, dagNode node, trs []trapezoid, err error) {
	trs = make([]trapezoid, 0)

	isLastTr := t.leftp.x <= s.endPoint.x && t.rightp.x >= s.endPoint.x

	if oldUT == nil && oldBT == nil || isLastTr { // first or last iteration
		return t.addSegmentInside(oldUT, oldBT, s)
	}

	// middle segment
	newUT, newBT, dagNode, err = t.splitY(oldUT, oldBT, s)
	trs = append(trs, *newUT, *newBT)
	return
}

func (t trapezoid) splitY(oldUT, oldBT *trapezoid, s Segment) (ut, bt *trapezoid, y node, err error) {
	if t.segmentInTrapezoid(s) {
		ut = &trapezoid{
			leftp:  s.startPoint,
			rightp: s.endPoint,
			top:    t.top,
			bottom: s,
		}
		bt = &trapezoid{
			leftp:  s.startPoint,
			rightp: s.endPoint,
			top:    s,
			bottom: t.bottom,
		}
	} else {
		if oldUT != nil && oldUT.top == t.top {
			ut = oldUT
		} else {
			ut = &trapezoid{
				leftp:  t.leftp,
				rightp: t.rightp,
				top:    t.top,
				bottom: s,
			}
			if oldUT != nil {
				oldUT.assignRightNeighborToTrapezoid(ut)
			}
		}
		if oldBT != nil && oldBT.bottom == t.bottom {
			bt = oldBT
		} else {
			bt = &trapezoid{
				leftp:  t.leftp,
				rightp: t.rightp,
				top:    s,
				bottom: t.bottom,
			}
			if oldBT != nil {
				oldBT.assignRightNeighborToTrapezoid(ut)
			}
		}
		if s.startPoint.x > t.leftp.x {
			ut.leftp = s.startPoint
			bt.leftp = s.startPoint
		} else {
			ut.leftp = t.leftp
			bt.leftp = t.leftp
		}
		if s.endPoint.x < t.rightp.x {
			ut.rightp = s.endPoint
			bt.rightp = s.endPoint
		} else {
			ut.rightp = t.rightp
			bt.rightp = t.rightp
		}
	}

	if oldUT != nil {
		if err = oldUT.replaceLeftNeighborsWith(ut); err != nil {
			return
		}
	} else {
		if err = t.replaceLeftNeighborsWith(ut); err != nil {
			return
		}
	}
	if oldBT != nil {
		if err = oldBT.replaceLeftNeighborsWith(bt); err != nil {
			return
		}
	} else {
		if err = t.replaceLeftNeighborsWith(bt); err != nil {
			return
		}
	}
	if err = t.replaceRightNeighborsWith(ut); err != nil {
		return
	}
	if err = t.replaceRightNeighborsWith(bt); err != nil {
		return
	}

	utNode := trapezoidNode{tr: ut, parents: make([]node, 0)}
	btNode := trapezoidNode{tr: bt, parents: make([]node, 0)}
	ut.dagRef = &utNode
	bt.dagRef = &btNode
	y = &yNode{s: s}
	y.assignLeft(&utNode)
	y.assignRight(&btNode)
	return
}

func boundingBox(ss []Segment) (tr trapezoid) {
	boundingBoxTop := Point{x: ss[0].maxX(), y: ss[0].maxY()}
	boundingBoxBot := Point{x: ss[0].minX(), y: ss[0].minY()}
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
		index++
	}
	boundingBoxTop.x += 5
	boundingBoxTop.y += 5
	boundingBoxBot.x -= 5
	boundingBoxBot.y -= 5

	bounderyTopSegment := NewSegment(
		Point{x: boundingBoxBot.x, y: boundingBoxTop.y},
		boundingBoxTop,
	)

	bounderyBotSegment := NewSegment(
		boundingBoxBot,
		Point{x: boundingBoxTop.x, y: boundingBoxBot.y},
	)

	tr = trapezoid{
		leftp:  boundingBoxBot,
		rightp: Point{x: boundingBoxTop.x, y: boundingBoxBot.y},
		top:    bounderyTopSegment,
		bottom: bounderyBotSegment,
	}
	return
}
