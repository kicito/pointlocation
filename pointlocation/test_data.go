package pointlocation

var dummyTrs1 = trapezoid{
	leftp:  Point{x: -11, y: -6},
	rightp: Point{x: -6, y: -3},
	top:    NewSegment(Point{x: -11, y: -6}, Point{x: -6, y: -3}),
	bottom: NewSegment(Point{x: -16, y: -11}, Point{x: -1, y: -11}),
}
var dummyTrs2 = trapezoid{
	leftp:  Point{x: -6, y: -3},
	rightp: Point{x: -1, y: -11},
	top:    NewSegment(Point{x: -16, y: 2}, Point{x: -1, y: 2}),
	bottom: NewSegment(Point{x: -16, y: -11}, Point{x: -1, y: -11}),
}
