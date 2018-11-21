package pointlocation

import (
	"fmt"
	"image/color"

	"github.com/pkg/errors"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
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

func (t trapezoid) equalTrapezoid(tr trapezoid, checkNeighbor bool) bool {
	if t.leftp.sameCoordinate(tr.leftp) &&
		t.rightp.sameCoordinate(tr.rightp) &&
		t.top.startPoint.sameCoordinate(tr.top.startPoint) &&
		t.top.endPoint.sameCoordinate(tr.top.endPoint) &&
		t.bottom.startPoint.sameCoordinate(tr.bottom.startPoint) &&
		t.bottom.endPoint.sameCoordinate(tr.bottom.endPoint) {

		if checkNeighbor {
			return t.upperLeftN.equalTrapezoid(*tr.upperLeftN, false) &&
				t.lowerRightN.equalTrapezoid(*tr.lowerRightN, false) &&
				t.upperRightN.equalTrapezoid(*tr.upperRightN, false) &&
				t.lowerRightN.equalTrapezoid(*tr.lowerRightN, false)
		}
		return true
	}
	return false
}

func (t trapezoid) topSegment() (Segment, error) {
	startY, err := t.top.y(t.leftp.x)
	if err != nil {
		return Segment{}, errors.Wrapf(err, "topsegment trapezoid %v", t)
	}
	startPoint := Point{x: t.leftp.x, y: startY}

	endY, err := t.top.y(t.rightp.x)
	if err != nil {
		return Segment{}, errors.Wrapf(err, "topsegment trapezoid %v", t)
	}
	endPoint := Point{x: t.rightp.x, y: endY}
	return NewSegment(startPoint, endPoint), nil
}

func (t trapezoid) bottomSegment() (Segment, error) {
	startY, err := t.bottom.y(t.leftp.x)
	if err != nil {
		return Segment{}, errors.Wrapf(err, "bottomsegment trapezoid %v", t)
	}
	startPoint := Point{x: t.leftp.x, y: startY}

	endY, err := t.bottom.y(t.rightp.x)
	if err != nil {
		return Segment{}, errors.Wrapf(err, "bottomsegment trapezoid %v", t)
	}
	endPoint := Point{x: t.rightp.x, y: endY}
	return NewSegment(startPoint, endPoint), nil

}

func (t trapezoid) leftSegment() (Segment, error) {
	startY, err := t.top.y(t.leftp.x)
	if err != nil {
		return Segment{}, errors.Wrapf(err, "leftSegment trapezoid %v", t)
	}
	startPoint := Point{x: t.leftp.x, y: startY}

	endY, err := t.bottom.y(t.leftp.x)
	if err != nil {
		return Segment{}, errors.Wrapf(err, "leftSegment trapezoid %v", t)
	}
	endPoint := Point{x: t.leftp.x, y: endY}
	return NewSegment(startPoint, endPoint), nil
}

func (t trapezoid) rightSegment() (Segment, error) {
	startY, err := t.top.y(t.rightp.x)
	if err != nil {
		return Segment{}, errors.Wrapf(err, "rightSegment trapezoid %v", t)
	}
	startPoint := Point{x: t.rightp.x, y: startY}

	endY, err := t.bottom.y(t.rightp.x)
	if err != nil {
		fmt.Println("yellow", t)
		return Segment{}, errors.Wrapf(err, "rightSegment trapezoid %v", t)
	}
	endPoint := Point{x: t.rightp.x, y: endY}
	return NewSegment(startPoint, endPoint), nil
}

// comparing if the Point is in side trapezoid by check orientation of Point to each Segment is inside trapeziod
func (t trapezoid) pointInTrapezoid(p Point) (b bool, err error) {
	var tmpSegment Segment
	var tmpOrt int
	if tmpSegment, err = t.leftSegment(); err != nil {
		return
	}
	tmpOrt = p.orientationFromSegment(tmpSegment)
	if tmpOrt == counterclockwise {
		b = false
		return
	}
	if tmpSegment, err = t.topSegment(); err != nil {
		return
	}
	tmpOrt = p.orientationFromSegment(tmpSegment)
	if tmpOrt == counterclockwise {
		b = false
		return
	}
	if tmpSegment, err = t.rightSegment(); err != nil {
		return
	}
	tmpOrt = p.orientationFromSegment(tmpSegment)
	if tmpOrt == clockwise {
		b = false
		return
	}

	if tmpSegment, err = t.bottomSegment(); err != nil {
		return
	}
	tmpOrt = p.orientationFromSegment(tmpSegment)
	if tmpOrt == clockwise {
		b = false
		return
	}

	b = true
	return
}

func (t trapezoid) segmentInTrapezoid(s Segment) (b bool, err error) {
	if b, err = t.pointInTrapezoid(s.startPoint); err != nil {
		return
	}
	if b, err = t.pointInTrapezoid(s.endPoint); err != nil {
		return
	}
	return
}

func (t *trapezoid) assignLeftNeighborToTrapezoid(tr *trapezoid) (err error) {
	// if t.leftp != tr.rightp {
	// 	return fmt.Errorf("trapezoidal %v cannot be left Neighbor to %v", tr, t)
	// }
	if tr == nil || *tr == (trapezoid{}) {
		t.upperLeftN = nil
		t.lowerLeftN = nil
		return
	}
	// if t.leftp.x == tr.rightp.x {
	if t.top == tr.top {
		t.upperLeftN = tr
		tr.upperRightN = t
	}

	if t.bottom == tr.bottom {
		t.lowerLeftN = tr
		tr.lowerRightN = t
	}
	return
	// }

	return nil
	// return fmt.Errorf("trapezoid %v can't be neighbor with %v, trapezoid is not connect", t, tr)
}

func (t *trapezoid) assignRightNeighborToTrapezoid(tr *trapezoid) (err error) {
	if tr == nil || *tr == (trapezoid{}) {
		t.upperRightN = nil
		t.lowerRightN = nil
		return
	}
	// if t.rightp.x == tr.leftp.x {
	if t.top == tr.top {
		t.upperRightN = tr
		tr.upperLeftN = t
	}
	if t.bottom == tr.bottom {
		t.lowerRightN = tr
		tr.lowerLeftN = t
	}
	return
	// }
	return nil

	// return fmt.Errorf("trapezoid %v can't be neighbor with %v, trapezoid is not connect", t, tr)
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

func (t trapezoid) addSegmentLeftTrapeziod(s Segment) (lt trapezoid, err error) {
	lt = trapezoid{
		leftp:  t.leftp,
		rightp: s.startPoint,
		top:    t.top,
		bottom: t.bottom,
	}
	return
}

func (t trapezoid) addSegmentRightTrapezoid(s Segment) (rt trapezoid, err error) {
	rt = trapezoid{
		leftp:  s.endPoint,
		rightp: t.rightp,
		top:    t.top,
		bottom: t.bottom,
	}
	return
}

func (t trapezoid) addSegmentInside(oldUT, oldBT *trapezoid, s Segment) (lt, rt trapezoid, ut, bt *trapezoid, err error) {
	if ut, bt, err = t.splitY(oldUT, oldBT, s); err != nil {
		return
	}
	isFirstTr := t.leftp.x <= s.startPoint.x && t.rightp.x >= s.startPoint.x
	isLastTr := t.leftp.x <= s.endPoint.x && t.rightp.x >= s.endPoint.x

	// if left p is not on same coordinate as segment's endpoint, we do create lt
	if isFirstTr {
		if !t.leftp.sameCoordinate(s.startPoint) && !t.rightp.sameCoordinate(s.startPoint) {
			lt, err = t.addSegmentLeftTrapeziod(s)
			if err != nil {
				return
			}
		}

	}

	if isLastTr {
		if !t.leftp.sameCoordinate(s.endPoint) && !t.rightp.sameCoordinate(s.endPoint) {
			// if right p is not on same coordinate as segment's endpoint, we do create rt
			// else, we do nothing because splitY already done repalceing neighbor
			rt, err = t.addSegmentRightTrapezoid(s)
			if err != nil {
				return
			}
		}
	}

	return
}

func (t trapezoid) addSegment(oldUT, oldBT *trapezoid, s Segment) (newLT, newRT trapezoid, newUT, newBT *trapezoid, err error) {

	isLastTr := t.leftp.x <= s.endPoint.x && t.rightp.x >= s.endPoint.x

	if oldUT == nil && oldBT == nil || isLastTr { // first or last iteration
		if newLT, newRT, newUT, newBT, err = t.addSegmentInside(oldUT, oldBT, s); err != nil {
			err = errors.Wrap(err, "unable to add segment")
			return
		}
		return
	}

	// middle segment
	newUT, newBT, err = t.splitY(oldUT, oldBT, s)
	return
}

func (t trapezoid) splitY(oldUT, oldBT *trapezoid, s Segment) (ut, bt *trapezoid, err error) {
	if s.startPoint.x >= t.leftp.x &&
		s.endPoint.x <= t.rightp.x {
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
		if oldUT != nil && oldUT.top == t.top && oldUT.bottom == t.bottom {
			ut = oldUT
		} else {
			ut = &trapezoid{
				leftp:  t.leftp,
				rightp: t.rightp,
				top:    t.top,
				bottom: s,
			}
		}
		if oldBT != nil && oldBT.bottom == t.bottom && oldBT.top == t.top {
			bt = oldBT
		} else {
			bt = &trapezoid{
				leftp:  t.leftp,
				rightp: t.rightp,
				top:    s,
				bottom: t.bottom,
			}
		}
		if s.startPoint.x > t.leftp.x {
			ut.leftp = s.startPoint
			bt.leftp = s.startPoint
		}
		if s.endPoint.x < t.rightp.x {
			ut.rightp = s.endPoint
			bt.rightp = s.endPoint
		}
	}

	return
}

func (t *trapezoid) mergeWith(tr *trapezoid) error {
	if tr.bottom != t.bottom && tr.top != t.top {
		return fmt.Errorf("trapezoid cannot be merge")
	}
	var leftp Point
	var rightp Point
	if tr.leftp.x > t.leftp.x {
		leftp = t.leftp
	}

	if tr.rightp.x > t.rightp.x {
		rightp = tr.rightp
	}
	t.leftp = leftp
	t.rightp = rightp
	tr = t
	return nil
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
		if ss[index].minX() < boundingBoxBot.x {
			boundingBoxBot.x = ss[index].minX()
		}
		if ss[index].minY() < boundingBoxBot.y {
			boundingBoxBot.y = ss[index].minY()
		}
		index++
	}
	boundingBoxTop.x += 0.001
	boundingBoxTop.y += 0.001
	boundingBoxBot.x -= 0.001
	boundingBoxBot.y -= 0.001

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

func (tr trapezoid) plotData() (leftp, rightp *plotter.Scatter, top, bottom, boundl, boundr *plotter.Line, err error) {
	if leftp, err = tr.leftp.scatter(); err != nil {
		return
	}
	if rightp, err = tr.rightp.scatter(); err != nil {
		return
	}
	top, _ = tr.top.line()
	bottom, _ = tr.bottom.line()
	// if top, err = tr.top.lineWithXBound(tr.leftp.x, tr.rightp.x); err != nil {
	// 	return
	// }
	// if bottom, err = tr.bottom.lineWithXBound(tr.leftp.x, tr.rightp.x); err != nil {
	// 	return
	// }

	lseg, err := tr.leftSegment()
	if err != nil {
		return
	}
	if boundl, err = lseg.line(); err != nil {
		return
	}
	boundl.LineStyle.Width = vg.Points(1)
	boundl.LineStyle.Dashes = []vg.Length{vg.Points(5), vg.Points(5)}
	boundl.LineStyle.Color = color.RGBA{G: 255, A: 255}

	rseg, err := tr.rightSegment()
	if err != nil {
		return
	}
	if boundr, err = rseg.line(); err != nil {
		return
	}

	boundr.LineStyle.Width = vg.Points(1)
	boundr.LineStyle.Dashes = []vg.Length{vg.Points(5), vg.Points(5)}
	boundr.LineStyle.Color = color.RGBA{G: 255, A: 255}

	return
}
