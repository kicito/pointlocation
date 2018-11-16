package pointlocation

import (
	"reflect"
	"testing"
)

func Test_newSegment(t *testing.T) {
	type args struct {
		start Point
		end   Point
	}
	floatTwo := float64(2)
	floatOne := float64(1)
	tests := []struct {
		name string
		args args
		want Segment
	}{
		{
			name: "vertical Segment",
			args: args{
				start: Point{x: 2, y: 2},
				end:   Point{x: 2, y: 2},
			},
			want: Segment{
				startPoint: Point{x: 2, y: 2},
				endPoint:   Point{x: 2, y: 2},
				slope:      nil,
				yIntercept: nil,
			},
		},
		{
			name: "horizontal Segment",
			args: args{
				start: Point{x: 2, y: 2},
				end:   Point{x: 5, y: 2},
			},
			want: Segment{
				startPoint: Point{x: 2, y: 2},
				endPoint:   Point{x: 5, y: 2},
				slope:      new(float64),
				yIntercept: &floatTwo,
			},
		},
		{
			name: "normal Segment",
			args: args{
				start: Point{x: 1, y: 1},
				end:   Point{x: 5, y: 5},
			},
			want: Segment{
				startPoint: Point{x: 1, y: 1},
				endPoint:   Point{x: 5, y: 5},
				yIntercept: new(float64),
				slope:      &floatOne,
			},
		},
		{
			name: "swap Point Segment",
			args: args{
				start: Point{x: 5, y: 5},
				end:   Point{x: 1, y: 1},
			},
			want: Segment{
				startPoint: Point{x: 1, y: 1},
				endPoint:   Point{x: 5, y: 5},
				yIntercept: new(float64),
				slope:      &floatOne,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewSegment(tt.args.start, tt.args.end)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Segment.inBoundX() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func Test_segment_inBoundX(t *testing.T) {
	type args struct {
		x float64
	}
	tests := []struct {
		name   string
		fields Segment
		args   args
		want   bool
	}{
		{
			name:   "normal case",
			fields: NewSegment(Point{x: 1.0, y: 1.0}, Point{x: 3.0, y: 2.0}),
			args: args{
				x: 2.0,
			},
			want: true,
		},
		{
			name:   "horizontal case",
			fields: NewSegment(Point{x: 1.0, y: 2.0}, Point{x: 1.0, y: 4.0}),
			args: args{
				x: 2.0,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Segment{
				startPoint: tt.fields.startPoint,
				endPoint:   tt.fields.endPoint,
				slope:      tt.fields.slope,
				yIntercept: tt.fields.yIntercept,
			}
			if got := s.inBoundX(tt.args.x); got != tt.want {
				t.Errorf("Segment.inBoundX() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_segment_inBoundY(t *testing.T) {
	type args struct {
		y float64
	}
	tests := []struct {
		name   string
		fields Segment
		args   args
		want   bool
	}{
		{
			name:   "normal case",
			fields: NewSegment(Point{x: 1.0, y: 1.0}, Point{x: 3.0, y: 2.0}),
			args: args{
				y: 2.0,
			},
			want: true,
		},
		{
			name:   "vertical case",
			fields: NewSegment(Point{x: -1.0, y: 4.0}, Point{x: -1.0, y: 3.0}),
			args: args{
				y: 2.0,
			},
			want: false,
		},
		{
			name:   "horizontal case",
			fields: NewSegment(Point{x: -1.0, y: 4.0}, Point{x: -1.0, y: 3.0}),
			args: args{
				y: 2.0,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Segment{
				startPoint: tt.fields.startPoint,
				endPoint:   tt.fields.endPoint,
				slope:      tt.fields.slope,
				yIntercept: tt.fields.yIntercept,
			}
			if got := s.inBoundY(tt.args.y); got != tt.want {
				t.Errorf("Segment.inBoundY() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_segment_y(t *testing.T) {
	type args struct {
		x float64
	}
	tests := []struct {
		name    string
		fields  Segment
		args    args
		wantY   float64
		wantErr bool
	}{
		{
			name:    "horizontal case",
			fields:  NewSegment(Point{x: 1.0, y: 2.0}, Point{x: 3.0, y: 2.0}),
			args:    args{2.0},
			wantY:   2.0,
			wantErr: false,
		},
		{
			name:    "vertical case",
			fields:  NewSegment(Point{x: 1.0, y: 2.0}, Point{x: 1.0, y: 5.0}),
			args:    args{2.0},
			wantY:   0.0,
			wantErr: true,
		},
		{
			name:    "vertical case inbound",
			fields:  NewSegment(Point{x: 1.0, y: 2.0}, Point{x: 1.0, y: 5.0}),
			args:    args{1.0},
			wantY:   0.0,
			wantErr: true,
		},
		{
			name:    "normal case",
			fields:  NewSegment(Point{x: -10.46, y: -2.49}, Point{x: -6.58, y: 2.01}),
			args:    args{-8.52},
			wantY:   -0.24,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := tt.fields
			gotY, err := s.y(tt.args.x)
			if (err != nil) != tt.wantErr {
				t.Errorf("Segment.y() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !almostEqual(gotY, tt.wantY) {
				t.Errorf("Segment.y() = %v, want %v", gotY, tt.wantY)
			}
		})
	}
}

func Test_segment_x(t *testing.T) {
	type args struct {
		y float64
	}
	tests := []struct {
		name    string
		fields  Segment
		args    args
		wantX   float64
		wantErr bool
	}{
		{
			name:    "horizontal case",
			fields:  NewSegment(Point{x: 1.0, y: 2.0}, Point{x: 3.0, y: 2.0}),
			args:    args{10.0},
			wantErr: true,
		},
		{
			name:    "horizontal case inbound",
			fields:  NewSegment(Point{x: 1.0, y: 2.0}, Point{x: 3.0, y: 2.0}),
			args:    args{2.0},
			wantErr: true,
		},
		{
			name:   "vertical case",
			fields: NewSegment(Point{x: 1.0, y: 2.0}, Point{x: 1.0, y: 5.0}),
			args:   args{2.0},
			wantX:  1.0,
		},
		{
			name:    "normal case",
			fields:  NewSegment(Point{x: -1.0, y: -4}, Point{x: -6, y: -2}),
			args:    args{-3},
			wantX:   -8,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := tt.fields
			gotX, err := s.x(tt.args.y)
			if (err != nil) != tt.wantErr {
				t.Errorf("Segment.x() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotX, tt.wantX) {
				t.Errorf("Segment.x() = %v, want %v", gotX, tt.wantX)
			}
		})
	}
}

func Test_segment_inBound(t *testing.T) {
	type args struct {
		p Point
	}
	tests := []struct {
		name   string
		fields Segment
		args   args
		want   bool
	}{
		{
			name:   "outbound x",
			fields: NewSegment(Point{x: -10, y: -4}, Point{x: -6, y: -2}),
			args:   args{Point{x: -5, y: -3}},
			want:   false,
		},
		{
			name:   "outbound y",
			fields: NewSegment(Point{x: -10, y: -4}, Point{x: -6, y: -2}),
			args:   args{Point{x: -7, y: -1}},
			want:   false,
		},
		{
			name:   "inbound",
			fields: NewSegment(Point{x: -10, y: -4}, Point{x: -6, y: -2}),
			args:   args{Point{x: -9, y: -3}},
			want:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Segment{
				startPoint: tt.fields.startPoint,
				endPoint:   tt.fields.endPoint,
				slope:      tt.fields.slope,
				yIntercept: tt.fields.yIntercept,
			}
			if got := s.inBound(tt.args.p); got != tt.want {
				t.Errorf("Segment.inBound() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_segment_isSegmentIntersect(t *testing.T) {
	type args struct {
		so Segment
	}
	tests := []struct {
		name   string
		fields Segment
		args   args
		want   bool
	}{
		{
			name:   "simple parallel",
			fields: NewSegment(Point{x: 1, y: 1}, Point{x: 10, y: 1}),
			args:   args{NewSegment(Point{x: 1, y: 2}, Point{x: 10, y: 2})},
			want:   false,
		},
		{
			name:   "simple intersect",
			fields: NewSegment(Point{x: 10, y: 0}, Point{x: 0, y: 10}),
			args:   args{NewSegment(Point{x: 0, y: 0}, Point{x: 10, y: 10})},
			want:   true,
		},
		{
			name:   "simple no intersect",
			fields: NewSegment(Point{x: -5, y: -5}, Point{x: 0, y: 0}),
			args:   args{NewSegment(Point{x: 1, y: 1}, Point{x: 10, y: 10})},
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Segment{
				startPoint: tt.fields.startPoint,
				endPoint:   tt.fields.endPoint,
				slope:      tt.fields.slope,
				yIntercept: tt.fields.yIntercept,
			}
			if got := s.isSegmentIntersect(tt.args.so); got != tt.want {
				t.Errorf("Segment.isSegmentIntersect() = %v, want %v", got, tt.want)
			}
		})
	}
}
