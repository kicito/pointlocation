package pointlocation

import (
	"fmt"

	"github.com/pkg/errors"
)

type PointLocation struct {
	Trs []trapezoid
	DAG dag
}

func NewPointLocation(ss []Segment) (pl PointLocation, err error) {

	// Determine a bounding box R that contains all segments of S, and initialize
	// the trapezoidal map structure T and search structure D for it.
	tr := boundingBox(ss)
	pl.DAG = newDAG(&tr)
	pl.addTrapeziod(tr)
	segmentAdded := make([]Segment, 0)

	for index := range ss {
		if err = pl.addSegment(ss[index]); err != nil {
			return pl, errors.Wrapf(err, "added segment %v\n adding segment %v\n", segmentAdded, ss[index])
		}
		segmentAdded = append(segmentAdded, ss[index])
	}
	return
}

func (pl PointLocation) findTrapeziod(tr trapezoid) (result trapezoid, index int, err error) {
	for trIndex := range pl.Trs {
		if pl.Trs[trIndex] == tr {
			index = trIndex
			result = pl.Trs[trIndex]
			return
		}
	}
	err = errors.New("trapezoid not found")
	return
}

func (pl *PointLocation) addTrapeziod(tr trapezoid) (err error) {
	// if trapezoid is found, throw error
	if _, _, err := pl.findTrapeziod(tr); err != nil {
		pl.Trs = append(pl.Trs, tr)
		return nil
	}
	return errors.New("trapezoid is found in the system")
}

func (pl *PointLocation) removeTrapezoid(tr trapezoid) (err error) {
	for trIndex := range pl.Trs {
		if pl.Trs[trIndex] == tr {
			copy(pl.Trs[trIndex:], pl.Trs[trIndex+1:])
			pl.Trs[len(pl.Trs)-1] = trapezoid{}
			pl.Trs = pl.Trs[:len(pl.Trs)-1]
			return nil
		}
	}
	return errors.New(fmt.Sprintf("unable to find %v", tr))
}

func (pl PointLocation) findIntersection(s Segment) (trs []trapezoid, err error) {
	startNode := pl.DAG.FindPoint(s.startPoint)
	startTr := startNode.(*trapezoidNode).tr
	trs = append(trs, *startTr)
	// finding if the segment endpoint is still on right of trapeziod right point, if so,
	// look if the segment is upper or below right point to determine neighbor to return
	currTr := startTr
	for currTr.rightp.x < s.endPoint.x {
		var segmentY float64
		if segmentY, err = s.y(currTr.rightp.x); err != nil {
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

func (pl *PointLocation) addSegment(s Segment) (err error) {
	intersectedTrapeziods, err := pl.findIntersection(s)
	if err != nil {
		err = errors.Wrapf(err, "unable to find intersected node in dag, with segment %v", s)
	}

	var trTop, trBot *trapezoid
	var newNode node
	var newTrapeziodals []trapezoid

	for trIndex := range intersectedTrapeziods {
		trTop, trBot, newNode, newTrapeziodals, err = intersectedTrapeziods[trIndex].addSegment(trTop, trBot, s)
		if err != nil {
			return errors.Wrapf(err, "fail adding segment %v", s)
		}
		if err = pl.removeTrapezoid(intersectedTrapeziods[trIndex]); err != nil {
			return errors.Wrap(err, "unable to remove trapeziod from pointlocation list")
		}
		for newTrIndex := range newTrapeziodals {
			if err = pl.addTrapeziod(newTrapeziodals[newTrIndex]); err != nil {
				return errors.Wrap(err, "unable to add trapeziod from pointlocation list")
			}
		}

		if intersectedTrapeziods[trIndex].dagRef == pl.DAG.root {
			pl.DAG.root = newNode
		} else {
			intersectedTrapeziods[trIndex].dagRef = intersectedTrapeziods[trIndex].dagRef.replaceWith(newNode)
		}
	}

	return
}
