package pointlocation

import (
	"fmt"

	"gonum.org/v1/plot/vg"

	"github.com/pkg/errors"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
)

type PointLocation struct {
	Trs []*trapezoid
	DAG dag
}

func NewPointLocation(ss []Segment) (pl PointLocation, err error) {

	// Determine a bounding box R that contains all segments of S, and initialize
	// the trapezoidal map structure T and search structure D for it.
	tr := boundingBox(ss)
	pl.DAG = newDAG(&tr)
	pl.addTrapeziod(&tr)
	segmentAdded := make([]Segment, 0)

	for index := range ss {
		if index == len(ss)-2 {
			fmt.Println("yo")
		}
		pl.PlotTrsWithSegment(fmt.Sprintf("[%v]segment %v", index, ss[index]), ss[index])
		if err = pl.addSegment(ss[index]); err != nil {
			return pl, errors.Wrapf(err, "added segment %v\n adding segment %v\n", segmentAdded, ss[index])
		}
		segmentAdded = append(segmentAdded, ss[index])
	}
	trLen := len(pl.Trs)
	for trIndex := 0; trIndex < trLen; trIndex++ {
		if pl.Trs[trIndex].dagRef == nil {
			if err = pl.removeTrapezoid(*pl.Trs[trIndex]); err != nil {
				return
			}
			trLen--
		}
	}
	return
}

func (pl PointLocation) findTrapeziod(tr trapezoid) (result *trapezoid, index int, err error) {
	for trIndex := range pl.Trs {
		if *pl.Trs[trIndex] == tr {
			index = trIndex
			result = pl.Trs[trIndex]
			return
		}
	}
	err = errors.New("trapezoid not found")
	return
}

func (pl PointLocation) replaceTrapezoid(origin, changed *trapezoid) error {
	for trIndex := range pl.Trs {
		if pl.Trs[trIndex] == origin {
			oldNode := pl.Trs[trIndex].dagRef.(*trapezoidNode)
			oldNode.tr = changed
			changed.dagRef = oldNode
			pl.Trs[trIndex] = changed
			return nil
		}
	}
	err := errors.New("trapezoid not found")
	return err
}

func (pl *PointLocation) addTrapeziod(tr *trapezoid) (err error) {
	// if trapezoid is found, throw error
	if _, _, err := pl.findTrapeziod(*tr); err != nil {
		pl.Trs = append(pl.Trs, tr)
		return nil
	}
	return errors.New("trapezoid is found in the system")
}

func (pl *PointLocation) removeTrapezoid(tr trapezoid) (err error) {
	for trIndex := range pl.Trs {
		if pl.Trs[trIndex].leftp.sameCoordinate(tr.leftp) &&
			pl.Trs[trIndex].top == tr.top &&
			pl.Trs[trIndex].bottom == tr.bottom &&
			pl.Trs[trIndex].rightp.sameCoordinate(tr.rightp) {

			pl.Trs[trIndex].replaceLeftNeighborsWith(nil)
			pl.Trs[trIndex].replaceRightNeighborsWith(nil)
			copy(pl.Trs[trIndex:], pl.Trs[trIndex+1:])
			pl.Trs[len(pl.Trs)-1] = nil
			pl.Trs = pl.Trs[:len(pl.Trs)-1]
			return nil
		}
	}
	return nil
	// return fmt.Errorf("unable to find %v", tr)
}

func (pl PointLocation) findIntersection(s Segment) (trs []*trapezoid, err error) {
	startNode, err := pl.DAG.FindPoint(s.startPoint)
	if err != nil {
		err = errors.Wrapf(err, "findIntersection %v", s)
		return
	}
	startTr := startNode.(*trapezoidNode).tr
	trs = append(trs, startTr)

	// finding if the segment endpoint is still on right of trapeziod right point, if so,
	// look if the segment is upper or below right point to determine neighbor to return
	currTr := startTr
	if currTr == nil {
		fmt.Print("waaaaat")
	}
	for currTr.rightp.x < s.endPoint.x {
		var segmentY float64
		if segmentY, err = s.y(currTr.rightp.x); err != nil {
			err = errors.Wrapf(err, "unable to get y while finding Intersection segment: %v, tr: %v", s, currTr)
			return
		}
		if currTr.rightp.y > segmentY {
			trs = append(trs, currTr.lowerRightN)
			currTr = currTr.lowerRightN
		} else {
			trs = append(trs, currTr.upperRightN)
			currTr = currTr.upperRightN
		}
	}
	return
}

func (pl *PointLocation) addSegment(s Segment) (err error) {
	intersectedTrapeziods, err := pl.findIntersection(s)

	debugPl := PointLocation{
		Trs: intersectedTrapeziods,
	}

	debugPl.PlotTrs(fmt.Sprintf("intersected trapezoid %v", s))

	for plI := range pl.Trs {
		if pl.Trs[plI].lowerRightN != nil && *pl.Trs[plI].lowerRightN == (trapezoid{}) {
			fmt.Println("what happened before adding ", s)
		}
	}

	if err != nil {
		err = errors.Wrapf(err, "unable to find intersected node in dag, with segment %v", s)
		return
	}

	var trLeft, trRight trapezoid
	var trTop, trBot *trapezoid
	var newNode node
	var lastTop, lastBot *trapezoid
	var skipTop, skipBot bool

	for trIndex := range intersectedTrapeziods {
		skipTop, skipBot = false, false
		trLeft, trRight, trTop, trBot, err = intersectedTrapeziods[trIndex].addSegment(trTop, trBot, s)
		if err != nil {
			return errors.Wrapf(err, "fail adding segment %v", s)
		}

		if lastBot != nil &&
			lastBot.top == trBot.top &&
			lastBot.bottom == trBot.bottom {
			if err = lastBot.mergeWith(trBot); err != nil {
				return
			}
			skipBot = true
		}
		if lastTop != nil &&
			lastTop.top == trTop.top &&
			lastTop.bottom == trTop.bottom {
			if err = lastTop.mergeWith(trTop); err != nil {
				return
			}
			skipTop = true
		}
		if trLeft != (trapezoid{}) && trRight != (trapezoid{}) {
			if err = intersectedTrapeziods[trIndex].replaceLeftNeighborsWith(&trLeft); err != nil {
				return
			}
			if err = intersectedTrapeziods[trIndex].replaceRightNeighborsWith(&trRight); err != nil {
				return
			}
			if err = trLeft.assignRightNeighborToTrapezoid(trTop); err != nil {
				return
			}
			if err = trLeft.assignRightNeighborToTrapezoid(trBot); err != nil {
				return
			}
			if err = trRight.assignLeftNeighborToTrapezoid(trTop); err != nil {
				return
			}
			if err = trRight.assignLeftNeighborToTrapezoid(trBot); err != nil {
				return
			}
		} else if trLeft != (trapezoid{}) {
			if err = intersectedTrapeziods[trIndex].replaceLeftNeighborsWith(&trLeft); err != nil {
				return
			}
			if err = intersectedTrapeziods[trIndex].replaceRightNeighborsWith(trTop); err != nil {
				return
			}
			if err = intersectedTrapeziods[trIndex].replaceRightNeighborsWith(trBot); err != nil {
				return
			}
			if err = trLeft.assignRightNeighborToTrapezoid(trTop); err != nil {
				return
			}
			if err = trLeft.assignRightNeighborToTrapezoid(trBot); err != nil {
				return
			}
		} else if trRight != (trapezoid{}) {
			if err = intersectedTrapeziods[trIndex].replaceRightNeighborsWith(&trRight); err != nil {
				return
			}

			if trIndex > 0 {
				if err = lastTop.replaceLeftNeighborsWith(trTop); err != nil {
					return
				}
				if err = lastBot.replaceLeftNeighborsWith(trBot); err != nil {
					return
				}
			} else {
				if err = intersectedTrapeziods[trIndex].replaceLeftNeighborsWith(trTop); err != nil {
					return
				}
				if err = intersectedTrapeziods[trIndex].replaceLeftNeighborsWith(trBot); err != nil {
					return
				}
			}

			if err = trRight.assignLeftNeighborToTrapezoid(trTop); err != nil {
				return
			}
			if err = trRight.assignLeftNeighborToTrapezoid(trBot); err != nil {
				return
			}
		} else {
			if err = intersectedTrapeziods[trIndex].replaceRightNeighborsWith(trTop); err != nil {
				return
			}
			if err = intersectedTrapeziods[trIndex].replaceRightNeighborsWith(trBot); err != nil {
				return
			}
			if trIndex > 0 {
				if err = lastTop.replaceLeftNeighborsWith(trTop); err != nil {
					return
				}
				if err = lastBot.replaceLeftNeighborsWith(trBot); err != nil {
					return
				}
			} else {
				if err = intersectedTrapeziods[trIndex].replaceLeftNeighborsWith(trTop); err != nil {
					return
				}
				if err = intersectedTrapeziods[trIndex].replaceLeftNeighborsWith(trBot); err != nil {
					return
				}
			}
		}

		newNode = pl.createNode(&trLeft, trTop, &trRight, trBot)

		if intersectedTrapeziods[trIndex].dagRef == pl.DAG.root {
			pl.DAG.root = newNode
		} else {
			intersectedTrapeziods[trIndex].dagRef.replaceWith(newNode)
		}

		if trLeft != (trapezoid{}) {
			if err = pl.addTrapeziod(&trLeft); err != nil {
				return errors.Wrap(err, "unable to add trapeziod from pointlocation list")
			}
		}

		if trRight != (trapezoid{}) {
			if err = pl.addTrapeziod(&trRight); err != nil {
				return errors.Wrap(err, "unable to add trapeziod from pointlocation list")
			}
		}
		if !skipTop {
			if err = pl.addTrapeziod(trTop); err != nil {
				return errors.Wrap(err, "unable to add trapeziod from pointlocation list")
			}
			lastTop = trTop
		}
		if !skipBot {
			if err = pl.addTrapeziod(trBot); err != nil {
				return errors.Wrap(err, "unable to add trapeziod from pointlocation list")
			}
			lastBot = trBot
		}
		pl.PlotTrs(fmt.Sprintf("[%v]4afteradded", trIndex))

		debugPl := PointLocation{
			Trs: []*trapezoid{
				trTop,
				trBot,
			},
		}
		if trLeft != (trapezoid{}) {
			debugPl.Trs = append(debugPl.Trs, &trLeft)
		}

		if trRight != (trapezoid{}) {
			debugPl.Trs = append(debugPl.Trs, &trRight)
		}

		debugPl.PlotTrs(fmt.Sprintf("[%v]5added trapezoid", trIndex))

		for plI := range pl.Trs {
			if pl.Trs[plI].lowerRightN != nil && *pl.Trs[plI].lowerRightN == (trapezoid{}) {
				fmt.Println("what happened while adding ", s)
			}
		}

	}

	pl.PlotTrs(fmt.Sprintf("beforeremove %v", s))
	for intersectIndex := range intersectedTrapeziods {
		PlotTr(fmt.Sprintf("[%v]removing", intersectIndex), *intersectedTrapeziods[intersectIndex])
		if err = pl.removeTrapezoid(*intersectedTrapeziods[intersectIndex]); err != nil {
			return errors.Wrap(err, "unable to remove trapeziod from pointlocation list")
		}
	}
	pl.PlotTrs(fmt.Sprintf("afterremove %v", s))

	return
}

func (pl PointLocation) createNode(lt, ut, rt, bt *trapezoid) node {
	var ltNode, ltXNode, rtNode, rtXNode, yyNode node
	if lt != nil && *lt != (trapezoid{}) {
		ltNode = &trapezoidNode{tr: lt, parents: make([]node, 0)}
		ltXNode = &xNode{xCoordinate: lt.rightp.x}
		lt.dagRef = ltNode
		ltXNode.assignLeft(ltNode)
	}
	if rt != nil && *rt != (trapezoid{}) {
		rtNode = &trapezoidNode{tr: rt, parents: make([]node, 0)}
		rtXNode = &xNode{xCoordinate: rt.leftp.x}
		rt.dagRef = rtNode
		rtXNode.assignRight(rtNode)
	}

	utNode := &trapezoidNode{tr: ut, parents: make([]node, 0)}
	btNode := &trapezoidNode{tr: bt, parents: make([]node, 0)}
	yyNode = &yNode{s: ut.bottom}
	yyNode.assignLeft(utNode)
	yyNode.assignRight(btNode)
	ut.dagRef = utNode
	bt.dagRef = btNode

	if lt != nil && *lt != (trapezoid{}) && rt != nil && *rt != (trapezoid{}) {
		ltXNode.assignRight(rtXNode)
		rtXNode.assignLeft(yyNode)
		return ltXNode
	} else if lt != nil && *lt != (trapezoid{}) {
		ltXNode.assignRight(yyNode)
		return ltXNode
	} else if rt != nil && *rt != (trapezoid{}) {
		rtXNode.assignLeft(yyNode)
		return rtXNode
	}
	return yyNode
}

func (pl PointLocation) PlotTrs(filename string) (err error) {

	p, err := plot.New()
	if err != nil {
		return
	}

	p.Title.Text = "Trapezoidal Map"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"
	p.Add(plotter.NewGrid())

	for trIndex := range pl.Trs {
		// PlotTr(fmt.Sprintf("a%v", trIndex), pl.Trs[trIndex])
		if err = plotAddTr(p, *pl.Trs[trIndex]); err != nil {
			return
		}
	}

	if filename == "" {
		filename = "points"
	}
	// Save the plot to a PNG file.
	if err = p.Save(10*vg.Inch, 10*vg.Inch, fmt.Sprintf("%v.png", filename)); err != nil {
		return
	}
	return
}

func (pl PointLocation) PlotTrsWithSegment(filename string, s Segment) (err error) {
	p, err := plot.New()
	if err != nil {
		return
	}

	p.Title.Text = "Trapezoidal Map"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"
	p.Add(plotter.NewGrid())

	for trIndex := range pl.Trs {
		if err = plotAddTr(p, *pl.Trs[trIndex]); err != nil {
			return
		}
	}
	l, err := s.line()
	if err != nil {
		return
	}
	p.Add(l)

	if filename == "" {
		filename = "points"
	}

	// Save the plot to a PNG file.
	if err = p.Save(10*vg.Inch, 10*vg.Inch, fmt.Sprintf("%v.png", filename)); err != nil {
		return
	}
	return
}

func (pl PointLocation) PlotTrsWithPoint(filename string, po Point) (err error) {
	p, err := plot.New()
	if err != nil {
		return
	}

	p.Title.Text = "Trapezoidal Map"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"
	p.Add(plotter.NewGrid())

	for trIndex := range pl.Trs {
		// PlotTr(fmt.Sprintf("tr %v", trIndex), *pl.Trs[trIndex])
		if err = plotAddTr(p, *pl.Trs[trIndex]); err != nil {
			return
		}
	}
	l, err := po.scatter()
	if err != nil {
		return
	}
	p.Add(l)

	if filename == "" {
		filename = "points"
	}

	// Save the plot to a PNG file.
	if err = p.Save(10*vg.Inch, 10*vg.Inch, fmt.Sprintf("%v.png", filename)); err != nil {
		return
	}
	return
}
