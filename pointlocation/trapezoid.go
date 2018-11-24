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
	label       string
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

	return fmt.Sprintf(`trapezoid[%v] lp: %v
	t: %v
	rp: %v 
	b: %v
	%v
	%v
	%v`,
		t.label,
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
	return NewSegment(startPoint, endPoint, "topsegment tapezoid "+t.label), nil
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
	return NewSegment(startPoint, endPoint, "bottomsegment tapezoid "+t.label), nil

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
	return NewSegment(startPoint, endPoint, "leftsegment tapezoid "+t.label), nil
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
	return NewSegment(startPoint, endPoint, "rightsegment tapezoid "+t.label), nil
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

	if !(t.leftp.x > tr.leftp.x) {
		debugPl := PointLocation{
			Trs: []*trapezoid{},
		}
		debugPl.Trs = append(debugPl.Trs, t, tr)
		debugPl.PlotTrs("error")
		return fmt.Errorf("this shouln't be called assigning %v \nleft neighbor to\n %v", t, tr)
	}

	if t.top == tr.top {
		t.upperLeftN = tr
		tr.upperRightN = t
	}

	if t.bottom == tr.bottom {
		t.lowerLeftN = tr
		tr.lowerRightN = t
	}
	return nil
}

func (t *trapezoid) assignRightNeighborToTrapezoid(tr *trapezoid) (err error) {

	if !(t.leftp.x < tr.leftp.x) {
		debugPl := PointLocation{
			Trs: []*trapezoid{},
		}
		debugPl.Trs = append(debugPl.Trs, t, tr)
		debugPl.PlotTrs("error")
		return fmt.Errorf("this shouln't be called assigning %v \nright neighbor to\n %v", t, tr)
	}

	if t.top == tr.top {
		t.upperRightN = tr
		tr.upperLeftN = t
	}
	if t.bottom == tr.bottom {
		t.lowerRightN = tr
		tr.lowerLeftN = t
	}
	return nil
}

func (t *trapezoid) replaceUpperLeftNeighborWith(tr *trapezoid) (err error) {
	if t.upperLeftN != nil {
		if (t.upperLeftN.upperRightN != nil && *t.upperLeftN.upperRightN == *t) ||
			(t.upperLeftN.lowerRightN != nil && *t.upperLeftN.lowerRightN == *t) {
			if tr != nil {
				return t.upperLeftN.assignRightNeighborToTrapezoid(tr)
			} else if t.upperLeftN.upperRightN != nil && *t.upperLeftN.upperRightN == *t {
				t.upperLeftN.upperRightN = nil
			} else if t.upperLeftN.lowerRightN != nil && *t.upperLeftN.lowerRightN == *t {
				t.upperLeftN.lowerRightN = nil
			}
		}
	}
	return
}

func (t *trapezoid) replaceLowerLeftNeighborWith(tr *trapezoid) (err error) {
	if t.lowerLeftN != nil {
		if (t.lowerLeftN.upperRightN != nil && *t.lowerLeftN.upperRightN == *t) ||
			(t.lowerLeftN.lowerRightN != nil && *t.lowerLeftN.lowerRightN == *t) {
			if tr != nil {
				return t.lowerLeftN.assignRightNeighborToTrapezoid(tr)
			} else if t.lowerLeftN.upperRightN != nil && *t.lowerLeftN.upperRightN == *t {
				t.lowerLeftN.upperRightN = nil
			} else if t.lowerLeftN.lowerRightN != nil && *t.lowerLeftN.lowerRightN == *t {
				t.lowerLeftN.lowerRightN = nil
			}
		}
	}
	return
}

func (t *trapezoid) replaceUpperRightNeighborWith(tr *trapezoid) (err error) {
	if t.upperRightN != nil {
		if (t.upperRightN.upperLeftN != nil && *t.upperRightN.upperLeftN == *t) ||
			(t.upperRightN.lowerLeftN != nil && *t.upperRightN.lowerLeftN == *t) {
			if tr != nil {
				return t.upperRightN.assignLeftNeighborToTrapezoid(tr)
			} else if t.upperRightN.upperLeftN != nil && *t.upperRightN.upperLeftN == *t {
				t.upperRightN.upperLeftN = nil
			} else if t.upperRightN.lowerLeftN != nil && *t.upperRightN.lowerLeftN == *t {
				t.upperRightN.lowerLeftN = nil
			}
		}
	}
	return
}

func (t *trapezoid) replaceLowerRightNeighborWith(tr *trapezoid) (err error) {
	if t.lowerRightN != nil {
		if (t.lowerRightN.upperLeftN != nil && *t.lowerRightN.upperLeftN == *t) ||
			(t.lowerRightN.lowerLeftN != nil && *t.lowerRightN.lowerLeftN == *t) {
			if tr != nil {
				return t.lowerRightN.assignLeftNeighborToTrapezoid(tr)
			} else if t.lowerRightN.upperLeftN != nil && *t.lowerRightN.upperLeftN == *t {
				t.lowerRightN.upperLeftN = nil
			} else if t.lowerRightN.lowerLeftN != nil && *t.lowerRightN.lowerLeftN == *t {
				t.lowerRightN.lowerLeftN = nil
			}
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

func (t trapezoid) addSegmentLeftTrapeziod(s Segment) (lt *trapezoid, err error) {
	lt = &trapezoid{
		leftp:  t.leftp,
		rightp: s.startPoint,
		top:    t.top,
		bottom: t.bottom,
		label:  fmt.Sprintf("lt for %v", s.label),
	}
	return
}

func (t trapezoid) addSegmentRightTrapezoid(s Segment) (rt *trapezoid, err error) {
	rt = &trapezoid{
		leftp:  s.endPoint,
		rightp: t.rightp,
		top:    t.top,
		bottom: t.bottom,
		label:  fmt.Sprintf("rt for %v", s.label),
	}
	return
}

func (t trapezoid) hasNeighborPointedTo() bool {
	if (t.lowerRightN.upperLeftN != nil && t.lowerRightN.upperLeftN == &t) ||
		(t.lowerRightN.lowerLeftN != nil && t.lowerRightN.lowerLeftN == &t) ||
		(t.upperRightN.upperLeftN != nil && t.upperRightN.upperLeftN == &t) ||
		(t.upperRightN.lowerLeftN != nil && t.upperRightN.lowerLeftN == &t) {
		return true
	}
	return false
}

func (t trapezoid) addSegment(s Segment) (newLT, newRT, newUT, newBT *trapezoid, err error) {

	isFirstTr := t.leftp.x <= s.startPoint.x && t.rightp.x >= s.startPoint.x
	isLastTr := t.leftp.x <= s.endPoint.x && t.rightp.x >= s.endPoint.x
	newUT, newBT, err = t.splitY(s)

	// if left p is not on same coordinate as segment's endpoint, we do create lt
	if isFirstTr {
		if !t.leftp.sameCoordinate(s.startPoint) && !t.rightp.sameCoordinate(s.startPoint) {
			newLT, err = t.addSegmentLeftTrapeziod(s)
			if err != nil {
				return
			}
		}
	}

	if isLastTr {
		if !t.leftp.sameCoordinate(s.endPoint) && !t.rightp.sameCoordinate(s.endPoint) {
			// if right p is not on same coordinate as segment's endpoint, we do create rt
			// else, we do nothing because splitY already done repalceing neighbor
			newRT, err = t.addSegmentRightTrapezoid(s)
			if err != nil {
				return
			}
		}
	}
	return
}

func (t trapezoid) splitY(s Segment) (ut, bt *trapezoid, err error) {
	if s.startPoint.x >= t.leftp.x &&
		s.endPoint.x <= t.rightp.x {
		ut = &trapezoid{
			leftp:  s.startPoint,
			rightp: s.endPoint,
			top:    t.top,
			bottom: s,
			label:  fmt.Sprintf("ut for %v", s.label),
		}
		bt = &trapezoid{
			leftp:  s.startPoint,
			rightp: s.endPoint,
			top:    s,
			bottom: t.bottom,
			label:  fmt.Sprintf("bt for %v", s.label),
		}
	} else {
		ut = &trapezoid{
			leftp:  t.leftp,
			rightp: t.rightp,
			top:    t.top,
			bottom: s,
			label:  fmt.Sprintf("ut for %v", s.label),
		}
		bt = &trapezoid{
			leftp:  t.leftp,
			rightp: t.rightp,
			top:    s,
			bottom: t.bottom,
			label:  fmt.Sprintf("bt for %v", s.label),
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
	tr.leftp = leftp
	tr.rightp = rightp
	t.leftp = leftp
	t.rightp = rightp
	tr.lowerLeftN = t.lowerLeftN
	tr.upperLeftN = t.upperLeftN
	t.lowerRightN = tr.lowerRightN
	t.upperRightN = tr.upperRightN
	t = tr
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
		"boundery tapezoid",
	)

	bounderyBotSegment := NewSegment(
		boundingBoxBot,
		Point{x: boundingBoxTop.x, y: boundingBoxBot.y},
		"boundery tapezoid",
	)

	tr = trapezoid{
		leftp:  boundingBoxBot,
		rightp: Point{x: boundingBoxTop.x, y: boundingBoxBot.y},
		top:    bounderyTopSegment,
		bottom: bounderyBotSegment,
	}
	return
}

func (tr trapezoid) plotData() (leftp, rightp *plotter.Scatter, top, bottom, boundl, boundr *plotter.Line, label *plotter.Labels, err error) {
	if leftp, err = tr.leftp.scatter(); err != nil {
		return
	}
	if rightp, err = tr.rightp.scatter(); err != nil {
		return
	}
	top, _ = tr.top.line()
	bottom, _ = tr.bottom.line()

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

	labelXY := make(plotter.XYs, 1)
	labelXY[0].X = (tr.leftp.x + tr.rightp.x) / 2.0
	labelXY[0].Y = ((tr.top.maxY()+tr.top.minY())/2.0 + (tr.bottom.maxY()+tr.bottom.minY())/2.0) / 2.0
	if label, err = plotter.NewLabels(plotter.XYLabels{XYs: labelXY, Labels: []string{tr.label}}); err != nil {
		return
	}

	return
}
