package pointlocation

import (
	"fmt"
	"testing"
)

func TestTrapezoidMap(t *testing.T) {
	type args struct {
		ss []Segment
		q  Point
	}
	floatZero := 0.0
	floatHalf := 0.5
	floatEight := 8.0
	floatMThree := -3.0
	tests := []struct {
		name   string
		args   args
		wantTr trapezoid
	}{
		{
			name: "simple 1",
			args: args{
				ss: []Segment{
					NewSegment(Point{x: 2, y: 2}, Point{x: 4, y: 3}),
					NewSegment(Point{x: 1, y: 1}, Point{x: 3, y: 2}),
				},
				q: Point{x: 0.5, y: 0.5},
			},
			wantTr: trapezoid{
				leftp:  &Point{x: -3, y: -3},
				rightp: &Point{x: 1, y: 1},
				top:    &Segment{startPoint: Point{x: -3, y: 8}, endPoint: Point{x: 9, y: 8}, slope: &floatZero, yIntercept: &floatEight},
				bottom: &Segment{startPoint: Point{x: -3, y: -3}, endPoint: Point{x: 9, y: -3}, slope: &floatZero, yIntercept: &floatMThree},
			},
		},
		{
			name: "simple 1.2",
			args: args{
				ss: []Segment{
					NewSegment(Point{x: 2, y: 2}, Point{x: 4, y: 3}),
					NewSegment(Point{x: 1, y: 1}, Point{x: 3, y: 2}),
				},
				q: Point{x: 5, y: 5},
			},
			wantTr: trapezoid{
				leftp:  &Point{x: 4, y: 3},
				rightp: &Point{x: 9, y: -3},
				top:    &Segment{startPoint: Point{x: -3, y: 8}, endPoint: Point{x: 9, y: 8}, slope: &floatZero, yIntercept: &floatEight},
				bottom: &Segment{startPoint: Point{x: -3, y: -3}, endPoint: Point{x: 9, y: -3}, slope: &floatZero, yIntercept: &floatMThree},
			},
		},
		{
			name: "simple 1.3",
			args: args{
				ss: []Segment{
					NewSegment(Point{x: 2, y: 2}, Point{x: 4, y: 3}),
					NewSegment(Point{x: 1, y: 1}, Point{x: 3, y: 2}),
				},
				q: Point{x: 2, y: -1},
			},
			wantTr: trapezoid{
				leftp:  &Point{x: 1, y: 1},
				rightp: &Point{x: 3, y: 2},
				top:    &Segment{startPoint: Point{x: 1, y: 1}, endPoint: Point{x: 3, y: 2}, slope: &floatHalf, yIntercept: &floatHalf},
				bottom: &Segment{startPoint: Point{x: -3, y: -3}, endPoint: Point{x: 9, y: -3}, slope: &floatZero, yIntercept: &floatMThree},
			},
		},
		{
			name: "simple 1.4",
			args: args{
				ss: []Segment{
					NewSegment(Point{x: 2, y: 2}, Point{x: 4, y: 3}),
					NewSegment(Point{x: 1, y: 1}, Point{x: 3, y: 2}),
				},
				q: Point{x: 3, y: 3},
			},
			wantTr: trapezoid{
				leftp:  &Point{x: 2, y: 2},
				rightp: &Point{x: 4, y: 3},
				top:    &Segment{startPoint: Point{x: -3, y: 8}, endPoint: Point{x: 9, y: 8}, slope: &floatZero, yIntercept: &floatEight},
				bottom: &Segment{startPoint: Point{x: 2, y: 2}, endPoint: Point{x: 4, y: 3}, slope: &floatHalf, yIntercept: &floatHalf},
			},
		},
		{
			name: "simple 1.5",
			args: args{
				ss: []Segment{
					NewSegment(Point{x: 2, y: 2}, Point{x: 4, y: 3}),
					NewSegment(Point{x: 1, y: 1}, Point{x: 3, y: 2}),
				},
				q: Point{x: 3.5, y: 2.5},
			},
			wantTr: trapezoid{
				leftp:  &Point{x: 3, y: 2},
				rightp: &Point{x: 4, y: 3},
				top:    &Segment{startPoint: Point{x: 2, y: 2}, endPoint: Point{x: 4, y: 3}, slope: &floatHalf, yIntercept: &floatHalf},
				bottom: &Segment{startPoint: Point{x: -3, y: -3}, endPoint: Point{x: 9, y: -3}, slope: &floatZero, yIntercept: &floatMThree},
			},
		},
		{
			name: "same startpoint",
			args: args{
				ss: []Segment{
					NewSegment(Point{x: 1, y: 1}, Point{x: 5, y: 2}),
					NewSegment(Point{x: 1, y: 1}, Point{x: 3, y: 3}),
				},
				q: Point{x: 2.5, y: 2},
			},
			wantTr: trapezoid{
				leftp:  &Point{x: 1, y: 1},
				rightp: &Point{x: 3, y: 3},
				top:    &Segment{startPoint: Point{x: 1, y: 1}, endPoint: Point{x: 3, y: 3}},
				bottom: &Segment{startPoint: Point{x: 1, y: 1}, endPoint: Point{x: 5, y: 2}},
			},
		},
		{
			name: "same startpoint 2",
			args: args{
				ss: []Segment{
					NewSegment(Point{x: 1, y: 1}, Point{x: 5, y: 2}),
					NewSegment(Point{x: 1, y: 1}, Point{x: 3, y: 3}),
				},
				q: Point{x: 3, y: 1},
			},
			wantTr: trapezoid{
				leftp:  &Point{x: 1, y: 1},
				rightp: &Point{x: 5, y: 2},
				top:    &Segment{startPoint: Point{x: 1, y: 1}, endPoint: Point{x: 5, y: 2}},
				bottom: &Segment{startPoint: Point{x: -4, y: -4}, endPoint: Point{x: 10, y: -4}},
			},
		},
		{
			name: "same startpoint 3",
			args: args{
				ss: []Segment{
					NewSegment(Point{x: 1, y: 1}, Point{x: 5, y: 2}),
					NewSegment(Point{x: 1, y: 1}, Point{x: 3, y: 3}),
				},
				q: Point{x: 4, y: 3},
			},
			wantTr: trapezoid{
				leftp:  &Point{x: 3, y: 3},
				rightp: &Point{x: 5, y: 2},
				top:    &Segment{startPoint: Point{x: -4, y: 8}, endPoint: Point{x: 10, y: 8}},
				bottom: &Segment{startPoint: Point{x: 1, y: 1}, endPoint: Point{x: 5, y: 2}},
			},
		},
		{
			name: "same startpoint 4",
			args: args{
				ss: []Segment{
					NewSegment(Point{x: 1, y: 1}, Point{x: 5, y: 2}),
					NewSegment(Point{x: 1, y: 1}, Point{x: 3, y: 3}),
				},
				q: Point{x: 2, y: 3},
			},
			wantTr: trapezoid{
				leftp:  &Point{x: 1, y: 1},
				rightp: &Point{x: 3, y: 3},
				top:    &Segment{startPoint: Point{x: -4, y: 8}, endPoint: Point{x: 10, y: 8}},
				bottom: &Segment{startPoint: Point{x: 1, y: 1}, endPoint: Point{x: 3, y: 3}},
			},
		},
		{
			name: "same startpoint 5",
			args: args{
				ss: []Segment{
					NewSegment(Point{x: 1, y: 1}, Point{x: 5, y: 2}),
					NewSegment(Point{x: 1, y: 1}, Point{x: 3, y: 3}),
				},
				q: Point{x: 0, y: 1},
			},
			wantTr: trapezoid{
				leftp:  &Point{x: -4, y: -4},
				rightp: &Point{x: 1, y: 1},
				top:    &Segment{startPoint: Point{x: -4, y: 8}, endPoint: Point{x: 10, y: 8}},
				bottom: &Segment{startPoint: Point{x: -4, y: -4}, endPoint: Point{x: 10, y: -4}},
			},
		},
		{
			name: "complex1",
			args: args{
				ss: []Segment{
					NewSegment(Point{x: 1, y: 1}, Point{x: 3, y: 3}),
					NewSegment(Point{x: 1, y: 1}, Point{x: 4, y: 0}),
					NewSegment(Point{x: 5, y: 4}, Point{x: 9, y: 2}),
					NewSegment(Point{x: 4, y: 0}, Point{x: 9, y: 2}),
					NewSegment(Point{x: 2, y: 1}, Point{x: 6, y: 2}),
				},
				q: Point{x: 2.5, y: 2},
			},
			wantTr: trapezoid{
				leftp:  &Point{x: 2, y: 1},
				rightp: &Point{x: 3, y: 3},
				top:    &Segment{startPoint: Point{x: 1, y: 1}, endPoint: Point{x: 3, y: 3}},
				bottom: &Segment{startPoint: Point{x: 2, y: 1}, endPoint: Point{x: 6, y: 2}},
			},
		},
		{
			name: "complex2",
			args: args{
				ss: []Segment{
					NewSegment(Point{x: 1, y: 1}, Point{x: 3, y: 3}),
					NewSegment(Point{x: 1, y: 1}, Point{x: 4, y: 0}),
					NewSegment(Point{x: 5, y: 4}, Point{x: 9, y: 2}),
					NewSegment(Point{x: 4, y: 0}, Point{x: 9, y: 2}),
					NewSegment(Point{x: 2, y: 1}, Point{x: 6, y: 2}),
				},
				q: Point{x: 2.5, y: 2},
			},
			wantTr: trapezoid{
				leftp:  &Point{x: 2, y: 1},
				rightp: &Point{x: 3, y: 3},
				top:    &Segment{startPoint: Point{x: 1, y: 1}, endPoint: Point{x: 3, y: 3}},
				bottom: &Segment{startPoint: Point{x: 2, y: 1}, endPoint: Point{x: 6, y: 2}},
			},
		},
		{
			name: "complex3",
			args: args{
				ss: []Segment{
					NewSegment(Point{x: 5, y: 4}, Point{x: 9, y: 2}),
					NewSegment(Point{x: 4, y: 0}, Point{x: 9, y: 2}),
				},
				q: Point{x: 7, y: 2},
			},
			wantTr: trapezoid{
				leftp:  &Point{x: 5, y: 4},
				rightp: &Point{x: 9, y: 2},
				top:    &Segment{startPoint: Point{x: 5, y: 4}, endPoint: Point{x: 9, y: 2}},
				bottom: &Segment{startPoint: Point{x: 4, y: 0}, endPoint: Point{x: 9, y: 2}},
			},
		},
		{
			name: "complex4",
			args: args{
				ss: []Segment{
					NewSegment(Point{x: 1, y: 1}, Point{x: 3, y: 3}),
					NewSegment(Point{x: 1, y: 1}, Point{x: 4, y: 0}),
					NewSegment(Point{x: 5, y: 4}, Point{x: 9, y: 2}),
					NewSegment(Point{x: 4, y: 0}, Point{x: 9, y: 2}),
					NewSegment(Point{x: 2, y: 1}, Point{x: 6, y: 2}),
				},
				q: Point{x: 4.5, y: 3},
			},
			wantTr: trapezoid{
				leftp:  &Point{x: 3, y: 3},
				rightp: &Point{x: 5, y: 4},
				top:    &Segment{startPoint: Point{x: 0, y: 9}, endPoint: Point{x: 14, y: 9}},
				bottom: &Segment{startPoint: Point{x: 2, y: 1}, endPoint: Point{x: 6, y: 2}},
			},
		},
		{
			name: "same rightp",
			args: args{
				ss: []Segment{
					NewSegment(Point{x: 1, y: -4}, Point{x: 4, y: -3}),
					NewSegment(Point{x: 6, y: -3}, Point{x: 8, y: -5}),
					NewSegment(Point{x: 4, y: -3}, Point{x: 6, y: -3}),
				},
				q: Point{x: 9, y: -4},
			},
			wantTr: trapezoid{
				leftp:  &Point{x: 8, y: -5},
				rightp: &Point{x: 13, y: -8},
				top:    &Segment{startPoint: Point{x: 1, y: 2}, endPoint: Point{x: 13, y: 2}},
				bottom: &Segment{startPoint: Point{x: 1, y: -8}, endPoint: Point{x: 13, y: -8}},
			},
		},
	}
	for _, tt := range tests {
		// for _, tt := range tests[len(tests)-1:] {
		t.Run(tt.name, func(t *testing.T) {
			gotD := TrapezoidMap(tt.args.ss)
			gotTr := *gotD.FindPoint(tt.args.q).(*trapezoidNode).tr
			fmt.Println("test result")
			fmt.Println(gotTr)
			fmt.Println("trapezoid right neighbors")
			fmt.Println(gotTr.lowerRightN, gotTr.upperRightN)
			fmt.Println("trapezoid left neighbors")
			fmt.Println(gotTr.lowerLeftN, gotTr.upperLeftN)
			if gotTr.leftp.x != tt.wantTr.leftp.x ||
				gotTr.leftp.y != tt.wantTr.leftp.y ||
				gotTr.rightp.x != tt.wantTr.rightp.x ||
				gotTr.rightp.y != tt.wantTr.rightp.y ||
				gotTr.top.startPoint.x != tt.wantTr.top.startPoint.x ||
				gotTr.bottom.startPoint.x != tt.wantTr.bottom.startPoint.x ||
				gotTr.top.startPoint.y != tt.wantTr.top.startPoint.y ||
				gotTr.bottom.startPoint.y != tt.wantTr.bottom.startPoint.y ||
				gotTr.top.endPoint.x != tt.wantTr.top.endPoint.x ||
				gotTr.bottom.endPoint.x != tt.wantTr.bottom.endPoint.x ||
				gotTr.top.endPoint.y != tt.wantTr.top.endPoint.y ||
				gotTr.bottom.endPoint.y != tt.wantTr.bottom.endPoint.y {
				t.Errorf("Point.positionBySegment() = %v, want %v", gotTr, tt.wantTr)
			}
		})
	}
}

func TestTrapezoidMapPolygon(t *testing.T) {
	type args struct {
		ss []Segment
		q  []Point
	}
	tests := []struct {
		name   string
		args   args
		wantTr []trapezoid
	}{
		{
			name: "simple 1",
			args: args{
				ss: []Segment{
					NewSegment(Point{x: 2, y: 2}, Point{x: 5, y: 5}),
					NewSegment(Point{x: 5, y: 5}, Point{x: 8, y: 2}),
					NewSegment(Point{x: 2, y: 2}, Point{x: 8, y: 2}),
				},
				q: []Point{
					{x: 0.55, y: 3.61},
					{x: 3.61, y: 0.63},
					{x: 6.57, y: 0.49},
					{x: 6.85, y: 4.43},
					{x: 2.63, y: 4.45},
					{x: 9.49, y: 3.35},
					{x: 4, y: 3},
					{x: 6, y: 3},
				},
			},
			wantTr: []trapezoid{
				{
					leftp:  &Point{x: 0, y: -3},
					rightp: &Point{x: 2, y: 2},
					top:    &Segment{startPoint: Point{x: 0, y: 10}, endPoint: Point{x: 13, y: 10}},
					bottom: &Segment{startPoint: Point{x: 0, y: -3}, endPoint: Point{x: 13, y: -3}},
				},
				{
					leftp:  &Point{x: 1, y: 2},
					rightp: &Point{x: 8, y: 2},
					top:    &Segment{startPoint: Point{x: 2, y: 2}, endPoint: Point{x: 8, y: 2}},
					bottom: &Segment{startPoint: Point{x: 0, y: -3}, endPoint: Point{x: 13, y: -3}},
				},
				{
					leftp:  &Point{x: 2, y: 2},
					rightp: &Point{x: 8, y: 2},
					top:    &Segment{startPoint: Point{x: 2, y: 2}, endPoint: Point{x: 8, y: 2}},
					bottom: &Segment{startPoint: Point{x: 0, y: -3}, endPoint: Point{x: 13, y: -3}},
				},
				{
					leftp:  &Point{x: 5, y: 5},
					rightp: &Point{x: 8, y: 2},
					top:    &Segment{startPoint: Point{x: 0, y: 10}, endPoint: Point{x: 13, y: 10}},
					bottom: &Segment{startPoint: Point{x: 5, y: 5}, endPoint: Point{x: 8, y: 2}},
				},
				{
					leftp:  &Point{x: 2, y: 2},
					rightp: &Point{x: 5, y: 5},
					top:    &Segment{startPoint: Point{x: 0, y: 10}, endPoint: Point{x: 13, y: 10}},
					bottom: &Segment{startPoint: Point{x: 2, y: 2}, endPoint: Point{x: 5, y: 5}},
				},
				{
					leftp:  &Point{x: 8, y: 2},
					rightp: &Point{x: 13, y: -3},
					top:    &Segment{startPoint: Point{x: 0, y: 10}, endPoint: Point{x: 13, y: 10}},
					bottom: &Segment{startPoint: Point{x: 0, y: -3}, endPoint: Point{x: 13, y: -3}},
				},
				{
					leftp:  &Point{x: 2, y: 2},
					rightp: &Point{x: 5, y: 5},
					top:    &Segment{startPoint: Point{x: 2, y: 2}, endPoint: Point{x: 5, y: 5}},
					bottom: &Segment{startPoint: Point{x: 2, y: 2}, endPoint: Point{x: 8, y: 2}},
				},
				{
					leftp:  &Point{x: 5, y: 5},
					rightp: &Point{x: 8, y: 2},
					top:    &Segment{startPoint: Point{x: 5, y: 5}, endPoint: Point{x: 8, y: 2}},
					bottom: &Segment{startPoint: Point{x: 2, y: 2}, endPoint: Point{x: 8, y: 2}},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotD := TrapezoidMap(tt.args.ss)
			for i := range tt.args.q {
				gotTr := *gotD.FindPoint(tt.args.q[i]).(*trapezoidNode).tr
				fmt.Printf("test result for %v\n", tt.args.q[i])
				fmt.Println(gotTr)
				fmt.Println("trapezoid right neighbors")
				fmt.Println(gotTr.lowerRightN, gotTr.upperRightN)
				fmt.Println("trapezoid left neighbors")
				fmt.Println(gotTr.lowerLeftN, gotTr.upperLeftN)
				if gotTr.leftp.x != tt.wantTr[i].leftp.x ||
					gotTr.leftp.y != tt.wantTr[i].leftp.y ||
					gotTr.rightp.x != tt.wantTr[i].rightp.x ||
					gotTr.rightp.y != tt.wantTr[i].rightp.y ||
					gotTr.top.startPoint.x != tt.wantTr[i].top.startPoint.x ||
					gotTr.bottom.startPoint.x != tt.wantTr[i].bottom.startPoint.x ||
					gotTr.top.startPoint.y != tt.wantTr[i].top.startPoint.y ||
					gotTr.bottom.startPoint.y != tt.wantTr[i].bottom.startPoint.y ||
					gotTr.top.endPoint.x != tt.wantTr[i].top.endPoint.x ||
					gotTr.bottom.endPoint.x != tt.wantTr[i].bottom.endPoint.x ||
					gotTr.top.endPoint.y != tt.wantTr[i].top.endPoint.y ||
					gotTr.bottom.endPoint.y != tt.wantTr[i].bottom.endPoint.y {
					t.Errorf("Point.positionBySegment(%v) = %v, want %v", tt.args.q[i], gotTr, tt.wantTr[i])
				}
			}
		})
	}
}
