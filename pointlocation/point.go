package pointlocation

import (
	"fmt"
	"math"
)

const float64EqualityThreshold = 1e-13

func almostEqual(a, b float64) bool {
	return math.Abs(a-b) <= float64EqualityThreshold
}

type Point struct {
	x float64
	y float64
	s *Segment
}

func NewPoint(x, y float64) Point {
	return Point{x: x, y: y}
}

func (p Point) String() string {
	return fmt.Sprintf("(%v, %v)", p.x, p.y)
}

const (
	upper     = iota // upper = 0
	lower            // lower = 1
	clockwise        // clockwise = 2
	counterclockwise
	colinear
)

func (p Point) positionBySegment(s Segment) (pos int, err error) {
	var segmentY float64

	if p.x == s.startPoint.x && p.y == s.startPoint.y {
		if p.s != nil && s.slope != nil && *p.s.slope > *s.slope {
			return upper, nil
		} else {
			return lower, nil
		}
	}
	if segmentY, err = s.y(p.x); err != nil {
		return
	}
	if p.y > segmentY {
		pos = upper
	} else {
		pos = lower
	}
	return
}

func (p Point) sameCoordinate(pp Point) bool {
	if p.x == pp.x && p.y == pp.y {
		return true
	}
	return false
}

//  the function checks if Point p lies on line Segment 's'
func (p Point) isPointOnSegment(s Segment) bool {

	if p.x <= s.maxX() && p.x >= s.minX() &&
		p.y <= s.maxY() && p.y >= s.minY() {
		return true
	}
	return false
}

/**
Slope of line Segment (p1, p2): σ = (y2 - y1)/(x2 - x1)
Slope of line Segment (p2, p3): τ = (y3 - y2)/(x3 - x2)

If  σ < τ, the orientation is counterclockwise (left turn)
If  σ = τ, the orientation is collinear
If  σ > τ, the orientation is clockwise (right turn)
*/
func (p Point) orientationFromSegment(s Segment) (pos int) {
	// fmt.Println("comparing Point", p, "with Segment", s)
	slopeSegmentEndPointToPoint := NewSegment(s.endPoint, p).slope
	if s.slope == nil {
		if p.x == s.startPoint.x {
			return colinear
		} else if p.x > s.startPoint.x {
			return clockwise
		}
		return counterclockwise
	}

	if slopeSegmentEndPointToPoint == nil {
		if p.y == s.endPoint.y {
			return colinear
		} else if p.y > s.endPoint.y {
			return counterclockwise
		}
		return clockwise
	}

	// horizontal case where s.slope == 0
	if *s.slope == 0 {
		if p.y == s.endPoint.y {
			return colinear
		} else if p.y > s.endPoint.y {
			return counterclockwise
		}
		return clockwise
	}

	// horizontal case where slopeSegmentEndPointToPoint == 0
	if *slopeSegmentEndPointToPoint == 0 {
		// fmt.Println(s, p, p.x > s.endPoint.x)
		if p.x == s.endPoint.x {
			return colinear
		} else if p.x > s.endPoint.x {
			if *s.slope < 0 {
				return counterclockwise
			}
			return clockwise
		}
		if *s.slope < 0 {
			return clockwise
		}
		return counterclockwise
	}

	// fmt.Printf("Segment %+v endsegment to Point %+v\n", s, NewSegment(s.endPoint, p))
	// fmt.Printf("slope Segment %+v, slope endsegment to Point %+v\n", *s.slope, *slopeSegmentEndPointToPoint)
	if *s.slope == *slopeSegmentEndPointToPoint {
		return colinear
	} else if *s.slope > *slopeSegmentEndPointToPoint {
		return clockwise
	} else {
		return counterclockwise
	}

}
