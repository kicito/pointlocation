package pointlocation

import (
	"fmt"
	"image/color"

	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

type Segment struct {
	startPoint Point
	endPoint   Point
	slope      *float64
	yIntercept *float64
}

func (s Segment) String() string {
	// var slope string
	// var yint string
	// if s.slope == nil {
	// 	slope = "-"
	// } else {
	// 	slope = fmt.Sprintf("%.2f", *s.slope)
	// }
	// if s.yIntercept == nil {
	// 	yint = "-"
	// } else {
	// 	yint = fmt.Sprintf("%.2f", *s.yIntercept)
	// }
	return fmt.Sprintf("%+v -> %+v", s.startPoint, s.endPoint)
}

type outOfBoundError struct {
	err   string
	input float64
}

func (e *outOfBoundError) Error() string {
	return fmt.Sprintf("Out of bound %v: %s", e.input, e.err)
}

func (s Segment) y(x float64) (y float64, err error) {
	if !s.inBoundX(x) {
		return y, &outOfBoundError{
			input: x,
			err:   fmt.Sprintf("x value, Segment bound x = %v - %v", s.minX(), s.maxX()),
		}
	}

	if s.slope == nil {
		// vertical line
		err = &outOfBoundError{
			input: y,
			err:   fmt.Sprintf("vertical line, unable to find y"),
		}

		return
	}

	if *s.slope == 0 {
		y = s.startPoint.y
		return
	}

	// y = slope * x + yintercept
	y = *s.slope*x + *s.yIntercept

	return
}

func (s Segment) x(y float64) (x float64, err error) {
	if !s.inBoundY(y) {
		err = &outOfBoundError{
			input: y,
			err:   fmt.Sprintf("y value, Segment bound y = %v - %v", s.minY(), s.maxY()),
		}
		return
	}

	if s.slope == nil {
		// vertical line
		x = s.startPoint.x
		return
	}

	if *s.slope == 0 {
		err = &outOfBoundError{
			input: y,
			err:   fmt.Sprintf("horizontal line, unable to find x"),
		}
		return
	}

	// x = y - yintercept / slope
	x = (y - *s.yIntercept) / *s.slope

	return
}

func (s Segment) inBound(p Point) bool {
	if !s.inBoundY(p.y) {
		return false
	}
	if !s.inBoundX(p.x) {
		return false
	}
	return true
}

func (s Segment) inBoundX(x float64) bool {
	return x >= s.minX() && x <= s.maxX()
}

func (s Segment) inBoundY(y float64) bool {
	return y >= s.minY() && y <= s.maxY()
}

func (s Segment) minX() float64 {
	if s.startPoint.x < s.endPoint.x {
		return s.startPoint.x
	}
	return s.endPoint.x
}

func (s Segment) minY() float64 {
	if s.startPoint.y < s.endPoint.y {
		return s.startPoint.y
	}
	return s.endPoint.y
}

func (s Segment) maxX() float64 {
	if s.startPoint.x < s.endPoint.x {
		return s.endPoint.x
	}
	return s.startPoint.x
}

func (s Segment) maxY() float64 {
	if s.startPoint.y < s.endPoint.y {
		return s.endPoint.y
	}
	return s.startPoint.y
}

func (s Segment) isSegmentIntersect(so Segment) bool {
	// https://www.geeksforgeeks.org/check-if-two-given-line-segments-intersect/

	o1 := so.startPoint.orientationFromSegment(s)
	o2 := so.endPoint.orientationFromSegment(s)
	o3 := s.startPoint.orientationFromSegment(so)
	o4 := s.endPoint.orientationFromSegment(so)
	fmt.Println(o1, o2, o3, o4)
	if o1 != o2 && o3 != o4 {
		return true
	}
	fmt.Println(so.startPoint.isPointOnSegment(s), so.endPoint.isPointOnSegment(s), s.startPoint.isPointOnSegment(so), s.endPoint.isPointOnSegment(so))

	if o1 == colinear && so.startPoint.isPointOnSegment(s) {
		return true
	}
	if o2 == colinear && so.endPoint.isPointOnSegment(s) {
		return true
	}
	if o3 == colinear && s.startPoint.isPointOnSegment(so) {
		return true
	}
	if o4 == colinear && s.endPoint.isPointOnSegment(so) {
		return true
	}
	return false
}

func NewSegment(start Point, end Point) Segment {
	var slope *float64
	var yIntercept *float64
	// var swapFlag bool

	// sort by lexicography
	if start.x > end.x {
		start, end = end, start
	}
	// else if start.x == end.x && start.y < end.y {
	// 	start, end = end, start
	// }

	if end.x == start.x {
		// vertical case
		slope = nil
		yIntercept = nil
	} else if end.y == start.y {
		// horizontal case
		slope = new(float64)
		yIntercept = &start.y
	} else {
		tmpSlope := (end.y - start.y) / (end.x - start.x)
		slope = &tmpSlope
		tmpYInt := start.y - tmpSlope*start.x
		yIntercept = &tmpYInt
	}

	s := Segment{
		startPoint: start,
		endPoint:   end,
		slope:      slope,
		yIntercept: yIntercept,
	}
	s.startPoint.s = &s
	s.endPoint.s = &s
	return s
}

func (s Segment) line() (line *plotter.Line, err error) {
	scatter := make(plotter.XYs, 2)
	scatter[0].X = s.startPoint.x
	scatter[0].Y = s.startPoint.y
	scatter[1].X = s.endPoint.x
	scatter[1].Y = s.endPoint.y
	line, err = plotter.NewLine(scatter)
	if err != nil {
		return
	}
	line.LineStyle.Width = vg.Points(3)
	line.LineStyle.Color = color.RGBA{B: 255, A: 255}
	return
}

func (s Segment) lineWithXBound(start, end float64) (line *plotter.Line, err error) {
	scatter := make(plotter.XYs, 2)
	if start > s.startPoint.x {
		scatter[0].X = start
	} else {
		scatter[0].X = s.startPoint.x
	}
	scatter[0].Y = s.startPoint.y

	if end < s.endPoint.x {
		scatter[1].X = end
	} else {
		scatter[1].X = s.endPoint.x
	}
	scatter[1].Y = s.endPoint.y
	line, err = plotter.NewLine(scatter)
	if err != nil {
		return
	}
	line.LineStyle.Width = vg.Points(3)
	line.LineStyle.Color = color.RGBA{B: 255, A: 255}
	return
}
