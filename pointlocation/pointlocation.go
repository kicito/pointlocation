package pointlocation

import (
	"fmt"
	"image/color"
	"math/rand"

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

	// random input segment order
	rand.Shuffle(len(ss), func(i, j int) {
		ss[i], ss[j] = ss[j], ss[i]
	})

	for index := range ss {
		fmt.Println("adding segment", ss[index])

		pl.PlotTrsWithSegment(fmt.Sprintf("./steps/[%v]1before add segment %v", index, ss[index]), ss[index])
		if err = pl.addSegment(ss[index], index); err != nil {
			return pl, errors.Wrapf(err, "added segment %v\n adding segment %v\n", segmentAdded, ss[index])
		}
		pl.PlotTrs(fmt.Sprintf("./steps/[%v]4after add segment %v", index, ss[index]))

		segmentAdded = append(segmentAdded, ss[index])
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
			copy(pl.Trs[trIndex:], pl.Trs[trIndex+1:])
			pl.Trs[len(pl.Trs)-1] = nil
			pl.Trs = pl.Trs[:len(pl.Trs)-1]
			return nil
		}
	}
	return fmt.Errorf("unable to find %v", tr)
}

func (pl *PointLocation) removeTrapezoidByIndex(index int) (err error) {

	pl.Trs[index].replaceLeftNeighborsWith(nil)
	pl.Trs[index].replaceRightNeighborsWith(nil)
	copy(pl.Trs[index:], pl.Trs[index+1:])
	pl.Trs[len(pl.Trs)-1] = nil
	pl.Trs = pl.Trs[:len(pl.Trs)-1]
	return
}

func (pl PointLocation) findIntersection(s Segment) (trs []trapezoid, err error) {
	startNode, err := pl.DAG.FindPoint(s.startPoint)
	if err != nil {
		err = errors.Wrapf(err, "findIntersection %v", s)
		return
	}
	startTr := startNode.(*trapezoidNode).tr
	trs = append(trs, *startTr)

	// finding if the segment endpoint is still on right of trapeziod right point, if so,
	// look if the segment is upper or below right point to determine neighbor to return
	currTr := startTr
	for currTr.rightp.x < s.endPoint.x {
		var segmentY float64
		if segmentY, err = s.y(currTr.rightp.x); err != nil {
			err = errors.Wrapf(err, "unable to get y while finding Intersection segment: %v, tr: %v", s, currTr)
			return
		}
		if currTr.rightp.y > segmentY {
			trs = append(trs, *currTr.lowerRightN)
			currTr = currTr.lowerRightN
		} else {
			trs = append(trs, *currTr.upperRightN)
			currTr = currTr.upperRightN
		}
	}
	return
}

func (pl *PointLocation) addSegment(s Segment, index int) (err error) {
	intersectedTrapeziods, err := pl.findIntersection(s)
	deleteIndices := make([]int, len(intersectedTrapeziods))
	for i := range intersectedTrapeziods {
		for j := range pl.Trs {
			if *pl.Trs[j] == intersectedTrapeziods[i] {
				deleteIndices = append(deleteIndices, j)
			}
		}
	}
	intersetedLen := len(intersectedTrapeziods)
	debugPl := PointLocation{}
	for i := range intersectedTrapeziods {
		debugPl.Trs = append(debugPl.Trs, &intersectedTrapeziods[i])
	}

	debugPl.PlotTrsWithSegment(fmt.Sprintf("./steps/[%v]2intersected trapezoid %v", index, s), s)

	if err != nil {
		err = errors.Wrapf(err, "unable to find intersected node in dag, with segment %v", s)
		return
	}

	var trLeft, trRight *trapezoid
	var newNode node
	var trTopList, trBotList []*trapezoid

	// create tr
	for trIndex := range intersectedTrapeziods {
		var tl, tr, ut, tb *trapezoid
		tl, tr, ut, tb, err = intersectedTrapeziods[trIndex].addSegment(s)
		if err != nil {
			return errors.Wrapf(err, "fail adding segment %v", s)
		}
		if tl != nil {
			trLeft = tl
		}
		if tr != nil {
			trRight = tr
		}

		trTopList = append(trTopList, ut)
		trBotList = append(trBotList, tb)
	}

	// fixing neighbors && merging if necessary
	if trLeft != nil {
		if err = intersectedTrapeziods[0].replaceLeftNeighborsWith(trLeft); err != nil {
			return
		}
		if err = trLeft.assignRightNeighborToTrapezoid(trTopList[0]); err != nil {
			return
		}
		if err = trLeft.assignRightNeighborToTrapezoid(trBotList[0]); err != nil {
			return
		}
	} else {
		if err = intersectedTrapeziods[0].replaceLeftNeighborsWith(trTopList[0]); err != nil {
			return
		}
		if err = intersectedTrapeziods[0].replaceLeftNeighborsWith(trBotList[0]); err != nil {
			return
		}
	}

	if trRight != nil {
		if err = intersectedTrapeziods[intersetedLen-1].replaceRightNeighborsWith(trRight); err != nil {
			return
		}
		if err = trRight.assignLeftNeighborToTrapezoid(trTopList[intersetedLen-1]); err != nil {
			return
		}
		if err = trRight.assignLeftNeighborToTrapezoid(trBotList[intersetedLen-1]); err != nil {
			return
		}
	} else {
		if err = intersectedTrapeziods[intersetedLen-1].replaceRightNeighborsWith(trTopList[intersetedLen-1]); err != nil {
			return
		}
		if err = intersectedTrapeziods[intersetedLen-1].replaceRightNeighborsWith(trBotList[intersetedLen-1]); err != nil {
			return
		}
	}

	if intersetedLen > 1 {
		for i := 1; i < intersetedLen; i++ {
			lastUT := trTopList[i-1]
			lastBT := trBotList[i-1]
			ut := trTopList[i]
			bt := trBotList[i]
			if lastUT.top == ut.top && lastUT.bottom == ut.bottom {
				if err = lastUT.mergeWith(ut); err != nil {
					return
				}
			} else {
				if err = intersectedTrapeziods[i].replaceLeftNeighborsWith(ut); err != nil {
					return
				}
				if err = ut.assignLeftNeighborToTrapezoid(lastUT); err != nil {
					return
				}
			}
			if lastBT.top == bt.top && lastBT.bottom == bt.bottom {
				if err = lastBT.mergeWith(bt); err != nil {
					return
				}
			} else {
				if err = intersectedTrapeziods[i].replaceLeftNeighborsWith(bt); err != nil {
					return
				}
				if err = bt.assignLeftNeighborToTrapezoid(lastBT); err != nil {
					return
				}
			}
		}
	}

	if intersetedLen == 1 {
		newNode = pl.createNode(trLeft, trTopList[0], trRight, trBotList[0])
		if intersectedTrapeziods[0].dagRef == pl.DAG.root {
			pl.DAG.root = newNode
		} else {
			intersectedTrapeziods[0].dagRef.replaceWith(newNode)
		}
	} else {
		for intersectedIndex := range intersectedTrapeziods {
			if intersectedIndex == 0 {
				if trLeft != nil {
					newNode = pl.createNode(trLeft, trTopList[0], nil, trBotList[0])
				} else {
					newNode = pl.createNode(nil, trTopList[0], nil, trBotList[0])
				}
			} else if intersectedIndex == intersetedLen-1 {
				if trRight != nil {
					newNode = pl.createNode(nil, trTopList[intersetedLen-1], trRight, trBotList[intersetedLen-1])
				} else {
					newNode = pl.createNode(nil, trTopList[intersetedLen-1], nil, trBotList[intersetedLen-1])
				}
			} else {
				newNode = pl.createNode(nil, trTopList[intersectedIndex], nil, trBotList[intersectedIndex])
			}

			intersectedTrapeziods[intersectedIndex].dagRef.replaceWith(newNode)
		}
	}

	for intersectIndex := range intersectedTrapeziods {
		PlotTr(fmt.Sprintf("./steps/[%v]removing trapez[%v]", index, intersectIndex), intersectedTrapeziods[intersectIndex])
		if err = pl.removeTrapezoid(intersectedTrapeziods[intersectIndex]); err != nil {
			return errors.Wrap(err, "unable to remove trapeziod from pointlocation list")
		}
	}

	if trLeft != nil {
		if err = pl.addTrapeziod(trLeft); err != nil {
			return errors.Wrap(err, "unable to add trapeziod from pointlocation list")
		}
	}
	if trRight != nil {
		if err = pl.addTrapeziod(trRight); err != nil {
			return errors.Wrap(err, "unable to add trapeziod from pointlocation list")
		}
	}

	debugPl = PointLocation{
		Trs: []*trapezoid{},
	}
	if trLeft != nil {
		debugPl.Trs = append(debugPl.Trs, trLeft)
	}
	if trRight != nil {
		debugPl.Trs = append(debugPl.Trs, trRight)
	}
	debugPl.Trs = append(debugPl.Trs, trTopList...)
	debugPl.Trs = append(debugPl.Trs, trTopList...)
	debugPl.Trs = append(debugPl.Trs, trBotList...)

	debugPl.PlotTrsWithSegment(fmt.Sprintf("./steps/[%v]3adding trapezoid %v", index, s), s)

	for intersectedIndex := range intersectedTrapeziods {
		if err = pl.addTrapeziod(trTopList[intersectedIndex]); err != nil {
			return errors.Wrap(err, "unable to add trapeziod from pointlocation list")
		}
		if err = pl.addTrapeziod(trBotList[intersectedIndex]); err != nil {
			return errors.Wrap(err, "unable to add trapeziod from pointlocation list")
		}
	}
	return
}

func (pl PointLocation) createNode(lt, ut, rt, bt *trapezoid) node {
	var ltNode, ltXNode, rtNode, rtXNode, yyNode node
	if lt != nil {
		ltNode = &trapezoidNode{tr: lt, parents: make([]node, 0)}
		ltXNode = &xNode{xCoordinate: lt.rightp.x}
		lt.dagRef = ltNode
		ltXNode.assignLeft(ltNode)
	}
	if rt != nil {
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

	if lt != nil && rt != nil {
		ltXNode.assignRight(rtXNode)
		rtXNode.assignLeft(yyNode)
		return ltXNode
	} else if lt != nil {
		ltXNode.assignRight(yyNode)
		return ltXNode
	} else if rt != nil {
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
	l.LineStyle.Color = color.RGBA{G: 255, A: 255}
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

func (pl PointLocation) PlotTrsWithPoint(filename string, po Point) (plo plot.Plotter, err error) {
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

func (pl PointLocation) PlotTrsDetail() {
	for trIndex := range pl.Trs {
		debugPl := PointLocation{}
		currTr := *pl.Trs[trIndex]
		currTr.label = "curr tr"
		debugPl.Trs = append(debugPl.Trs, &currTr)
		if currTr.lowerLeftN != nil {
			ll := *currTr.lowerLeftN
			ll.label = "lower left"
			debugPl.Trs = append(debugPl.Trs, currTr.lowerLeftN)
		}
		if currTr.lowerRightN != nil {
			lr := *currTr.lowerRightN
			lr.label = "lower right"
			debugPl.Trs = append(debugPl.Trs, &lr)
		}
		if currTr.upperRightN != nil {
			ur := *currTr.upperRightN
			ur.label = "upper right"
			debugPl.Trs = append(debugPl.Trs, &ur)
		}
		if currTr.upperLeftN != nil {
			ul := *currTr.upperLeftN
			ul.label = "upper left"
			debugPl.Trs = append(debugPl.Trs, &ul)
		}
		debugPl.PlotTrs(fmt.Sprintf("wtf trapezoid[%v]", trIndex))
	}
}
