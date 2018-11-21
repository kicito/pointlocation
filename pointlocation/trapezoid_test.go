package pointlocation

import (
	"fmt"
	"reflect"
	"testing"
)

func Test_trapezoid_topSegment(t *testing.T) {
	type fields struct {
		leftp  Point
		rightp Point
		top    Segment
		bottom Segment
	}
	bounderyTopSegment := NewSegment(Point{x: 0, y: 5}, Point{x: 5, y: 5})
	bounderyBotSegment := NewSegment(Point{x: 0, y: 0}, Point{x: 5, y: 0})
	tests := []struct {
		name   string
		fields fields
		want   Segment
	}{
		{
			name: "boundery",
			fields: fields{
				leftp:  Point{x: 0, y: 0},
				rightp: Point{x: 5, y: 0},
				top:    bounderyTopSegment,
				bottom: bounderyBotSegment,
			},
			want: bounderyTopSegment,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := trapezoid{
				leftp:  tt.fields.leftp,
				rightp: tt.fields.rightp,
				top:    tt.fields.top,
				bottom: tt.fields.bottom,
			}
			if got := tr.topSegment(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("trapezoid.topSegment() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_trapezoid_bottomSegment(t *testing.T) {
	type fields struct {
		leftp  Point
		rightp Point
		top    Segment
		bottom Segment
	}
	bounderyTopSegment := NewSegment(Point{x: 0, y: 5}, Point{x: 5, y: 5})
	bounderyBotSegment := NewSegment(Point{x: 0, y: 0}, Point{x: 5, y: 0})
	tests := []struct {
		name   string
		fields fields
		want   Segment
	}{
		{
			name: "boundery",
			fields: fields{
				leftp:  Point{x: 0, y: 0},
				rightp: Point{x: 5, y: 0},
				top:    bounderyTopSegment,
				bottom: bounderyBotSegment,
			},
			want: bounderyBotSegment,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := trapezoid{
				leftp:  tt.fields.leftp,
				rightp: tt.fields.rightp,
				top:    tt.fields.top,
				bottom: tt.fields.bottom,
			}
			if got := tr.bottomSegment(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("trapezoid.bottomSegment() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_trapezoid_leftSegment(t *testing.T) {
	type fields struct {
		leftp  Point
		rightp Point
		top    Segment
		bottom Segment
	}
	bounderyTopSegment := NewSegment(Point{x: 0, y: 5}, Point{x: 5, y: 5})
	bounderyBotSegment := NewSegment(Point{x: 0, y: 0}, Point{x: 5, y: 0})
	tests := []struct {
		name   string
		fields fields
		want   Segment
	}{
		{
			name: "boundery",
			fields: fields{
				leftp:  Point{x: 0, y: 0},
				rightp: Point{x: 5, y: 0},
				top:    bounderyTopSegment,
				bottom: bounderyBotSegment,
			},
			want: NewSegment(Point{x: 0, y: 5}, Point{x: 0, y: 0}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := trapezoid{
				leftp:  tt.fields.leftp,
				rightp: tt.fields.rightp,
				top:    tt.fields.top,
				bottom: tt.fields.bottom,
			}
			if got := tr.leftSegment(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("trapezoid.leftSegment() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_trapezoid_rightSegment(t *testing.T) {
	type fields struct {
		leftp  Point
		rightp Point
		top    Segment
		bottom Segment
	}
	bounderyTopSegment := NewSegment(Point{x: 0, y: 5}, Point{x: 5, y: 5})
	bounderyBotSegment := NewSegment(Point{x: 0, y: 0}, Point{x: 5, y: 0})
	tests := []struct {
		name   string
		fields fields
		want   Segment
	}{
		{
			name: "boundery",
			fields: fields{
				leftp:  Point{x: 0, y: 0},
				rightp: Point{x: 5, y: 0},
				top:    bounderyTopSegment,
				bottom: bounderyBotSegment,
			},
			want: NewSegment(Point{x: 5, y: 5}, Point{x: 5, y: 0}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := trapezoid{
				leftp:  tt.fields.leftp,
				rightp: tt.fields.rightp,
				top:    tt.fields.top,
				bottom: tt.fields.bottom,
			}
			if got := tr.rightSegment(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("trapezoid.rightSegment() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_trapezoid_pointInTrapezoid(t *testing.T) {
	type fields struct {
		leftp  Point
		rightp Point
		top    Segment
		bottom Segment
	}
	type args struct {
		p Point
	}
	bounderyTopSegment := NewSegment(Point{x: 0, y: 5}, Point{x: 5, y: 5})
	bounderyBotSegment := NewSegment(Point{x: 0, y: 0}, Point{x: 5, y: 0})
	tests := []struct {
		name   string
		fields fields
		args   args
		wantB  bool
	}{
		{
			name: "boundery",
			fields: fields{
				leftp:  Point{x: 0, y: 0},
				rightp: Point{x: 5, y: 0},
				top:    bounderyTopSegment,
				bottom: bounderyBotSegment,
			},
			args:  args{Point{x: 1, y: 1}},
			wantB: true,
		},
		{
			name: "boundery 2",
			fields: fields{
				leftp:  Point{x: 0, y: 0},
				rightp: Point{x: 5, y: 0},
				top:    bounderyTopSegment,
				bottom: bounderyBotSegment,
			},
			args:  args{Point{x: 2, y: 2}},
			wantB: true,
		},
		{
			name: "out of boundery",
			fields: fields{
				leftp:  Point{x: 0, y: 0},
				rightp: Point{x: 5, y: 0},
				top:    bounderyTopSegment,
				bottom: bounderyBotSegment,
			},
			args:  args{Point{x: 6, y: 6}},
			wantB: false,
		},
		{
			name: "out of boundery 2",
			fields: fields{
				leftp:  Point{x: 0, y: 0},
				rightp: Point{x: 5, y: 0},
				top:    bounderyTopSegment,
				bottom: bounderyBotSegment,
			},
			args:  args{Point{x: -1, y: -1}},
			wantB: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := trapezoid{
				leftp:  tt.fields.leftp,
				rightp: tt.fields.rightp,
				top:    tt.fields.top,
				bottom: tt.fields.bottom,
			}
			if gotB := tr.pointInTrapezoid(tt.args.p); gotB != tt.wantB {
				t.Errorf("trapezoid.pointInTrapezoid() = %v, want %v", gotB, tt.wantB)
			}
		})
	}
}

// func Test_trapezoid_addSegmentInside(t *testing.T) {
// 	type fields struct {
// 		leftp       Point
// 		rightp      Point
// 		top         Segment
// 		bottom      Segment
// 		upperLeftN  *trapezoid
// 		lowerLeftN  *trapezoid
// 		upperRightN *trapezoid
// 		lowerRightN *trapezoid
// 	}
// 	type args struct {
// 		s Segment
// 	}
// 	lp := Point{x: 0, y: 0}
// 	rp := Point{x: 5, y: 0}
// 	inputStart := Point{x: 2, y: 2}
// 	inputEnd := Point{x: 4, y: 3}
// 	bounderyTopSegment := NewSegment(Point{x: lp.x, y: 5}, Point{x: 5, y: 5})
// 	bounderyBotSegment := NewSegment(lp, rp)
// 	input := NewSegment(inputStart, inputEnd)

// 	lt := trapezoid{
// 		leftp:  lp,
// 		rightp: inputStart,
// 		top:    bounderyTopSegment,
// 		bottom: bounderyBotSegment,
// 	}
// 	ut := trapezoid{
// 		leftp:  inputStart,
// 		rightp: inputEnd,
// 		top:    bounderyTopSegment,
// 		bottom: input,
// 	}
// 	rt := trapezoid{
// 		leftp:  inputEnd,
// 		rightp: rp,
// 		top:    bounderyTopSegment,
// 		bottom: bounderyBotSegment,
// 	}
// 	bt := trapezoid{
// 		leftp:  inputStart,
// 		rightp: inputEnd,
// 		top:    input,
// 		bottom: bounderyBotSegment,
// 	}
// 	// assign neighbor
// 	lt.upperRightN = &ut
// 	lt.lowerRightN = &bt

// 	ut.upperLeftN = &lt
// 	ut.upperRightN = &rt

// 	rt.upperLeftN = &ut
// 	rt.lowerLeftN = &bt

// 	bt.upperLeftN = &lt
// 	bt.upperRightN = &rt
// 	tests := []struct {
// 		name   string
// 		fields fields
// 		args   args
// 		wantLt trapezoid
// 		wantUt trapezoid
// 		wantRt trapezoid
// 		wantBt trapezoid
// 	}{
// 		{
// 			name: "first add",
// 			fields: fields{
// 				leftp:  Point{x: 0, y: 0},
// 				rightp: Point{x: 5, y: 0},
// 				top:    bounderyTopSegment,
// 				bottom: bounderyBotSegment,
// 			},
// 			args:   args{input},
// 			wantLt: lt,
// 			wantUt: ut,
// 			wantRt: rt,
// 			wantBt: bt,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			tr := trapezoid{
// 				leftp:       tt.fields.leftp,
// 				rightp:      tt.fields.rightp,
// 				top:         tt.fields.top,
// 				bottom:      tt.fields.bottom,
// 				upperLeftN:  tt.fields.upperLeftN,
// 				lowerLeftN:  tt.fields.lowerLeftN,
// 				upperRightN: tt.fields.upperRightN,
// 				lowerRightN: tt.fields.lowerRightN,
// 			}
// 			gotLt, gotUt, gotRt, gotBt := tr.addSegmentInside(&t.args.s)
// 			if !reflect.DeepEqual(gotLt, tt.wantLt) {
// 				t.Errorf("trapezoid.addSegmentInside() gotLt = %v, want %v", gotLt, tt.wantLt)
// 			}
// 			if !reflect.DeepEqual(gotUt, tt.wantUt) {
// 				t.Errorf("trapezoid.addSegmentInside() gotUt = %v, want %v", gotUt, tt.wantUt)
// 			}
// 			if !reflect.DeepEqual(gotRt, tt.wantRt) {
// 				t.Errorf("trapezoid.addSegmentInside() gotRt = %v, want %v", gotRt, tt.wantRt)
// 			}
// 			if !reflect.DeepEqual(gotBt, tt.wantBt) {
// 				t.Errorf("trapezoid.addSegmentInside() gotBt = %v, want %v", gotBt, tt.wantBt)
// 			}
// 		})
// 	}
// }

func Test_trapezoid_assignRightNeighborToTrapezoid(t *testing.T) {
	type args struct {
		tr *trapezoid
	}
	tests := []struct {
		name    string
		fields  trapezoid
		args    args
		wantErr bool
	}{
		{
			name:   "connected trapezoid",
			fields: dummyTrs1,
			args: args{
				tr: &dummyTrs2,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := tt.fields
			if err := tr.assignRightNeighborToTrapezoid(tt.args.tr); (err != nil) != tt.wantErr {
				t.Errorf("trapezoid.assignRightNeighborToTrapezoid() error = %v, wantErr %v", err, tt.wantErr)
			}
			fmt.Println(tr)
		})
	}
}
