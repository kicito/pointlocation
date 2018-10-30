package pointlocation

import (
	"reflect"
	"testing"
)

func Test_newSegment(t *testing.T) {
	type args struct {
		start point
		end   point
	}
	floatTwo := float64(2)
	floatOne := float64(1)
	tests := []struct {
		name string
		args args
		want segment
	}{
		{
			name: "vertical segment",
			args: args{
				start: point{2, 2},
				end:   point{2, 2},
			},
			want: segment{
				startPoint: point{2, 2},
				endPoint:   point{2, 2},
				slope:      nil,
				yIntercept: nil,
			},
		},
		{
			name: "horizontal segment",
			args: args{
				start: point{2, 2},
				end:   point{5, 2},
			},
			want: segment{
				startPoint: point{2, 2},
				endPoint:   point{5, 2},
				slope:      new(float64),
				yIntercept: &floatTwo,
			},
		},
		{
			name: "normal segment",
			args: args{
				start: point{1, 1},
				end:   point{5, 5},
			},
			want: segment{
				startPoint: point{1, 1},
				endPoint:   point{5, 5},
				yIntercept: new(float64),
				slope:      &floatOne,
			},
		},
		{
			name: "swap point segment",
			args: args{
				start: point{5, 5},
				end:   point{1, 1},
			},
			want: segment{
				startPoint: point{1, 1},
				endPoint:   point{5, 5},
				yIntercept: new(float64),
				slope:      &floatOne,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newSegment(tt.args.start, tt.args.end)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("segment.inBoundX() = %+v, want %+v", got, tt.want)
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
		fields segment
		args   args
		want   bool
	}{
		{
			name:   "normal case",
			fields: newSegment(point{1.0, 1.0}, point{3.0, 2.0}),
			args: args{
				x: 2.0,
			},
			want: true,
		},
		{
			name:   "horizontal case",
			fields: newSegment(point{1.0, 2.0}, point{1.0, 4.0}),
			args: args{
				x: 2.0,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := segment{
				startPoint: tt.fields.startPoint,
				endPoint:   tt.fields.endPoint,
				slope:      tt.fields.slope,
				yIntercept: tt.fields.yIntercept,
			}
			if got := s.inBoundX(tt.args.x); got != tt.want {
				t.Errorf("segment.inBoundX() = %v, want %v", got, tt.want)
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
		fields segment
		args   args
		want   bool
	}{
		{
			name:   "normal case",
			fields: newSegment(point{1.0, 1.0}, point{3.0, 2.0}),
			args: args{
				y: 2.0,
			},
			want: true,
		},
		{
			name:   "vertical case",
			fields: newSegment(point{-1.0, 4.0}, point{-1.0, 3.0}),
			args: args{
				y: 2.0,
			},
			want: false,
		},
		{
			name:   "horizontal case",
			fields: newSegment(point{-1.0, 4.0}, point{-1.0, 3.0}),
			args: args{
				y: 2.0,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := segment{
				startPoint: tt.fields.startPoint,
				endPoint:   tt.fields.endPoint,
				slope:      tt.fields.slope,
				yIntercept: tt.fields.yIntercept,
			}
			if got := s.inBoundY(tt.args.y); got != tt.want {
				t.Errorf("segment.inBoundY() = %v, want %v", got, tt.want)
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
		fields  segment
		args    args
		wantY   float64
		wantErr bool
	}{
		{
			name:    "horizontal case",
			fields:  newSegment(point{1.0, 2.0}, point{3.0, 2.0}),
			args:    args{2.0},
			wantY:   2.0,
			wantErr: false,
		},
		{
			name:    "vertical case",
			fields:  newSegment(point{1.0, 2.0}, point{1.0, 5.0}),
			args:    args{2.0},
			wantY:   0.0,
			wantErr: true,
		},
		{
			name:    "vertical case inbound",
			fields:  newSegment(point{1.0, 2.0}, point{1.0, 5.0}),
			args:    args{1.0},
			wantY:   0.0,
			wantErr: true,
		},
		{
			name:    "normal case",
			fields:  newSegment(point{-10.46, -2.49}, point{-6.58, 2.01}),
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
				t.Errorf("segment.y() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !almostEqual(gotY, tt.wantY) {
				t.Errorf("segment.y() = %v, want %v", gotY, tt.wantY)
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
		fields  segment
		args    args
		wantX   float64
		wantErr bool
	}{
		{
			name:    "horizontal case",
			fields:  newSegment(point{1.0, 2.0}, point{3.0, 2.0}),
			args:    args{10.0},
			wantErr: true,
		},
		{
			name:    "horizontal case inbound",
			fields:  newSegment(point{1.0, 2.0}, point{3.0, 2.0}),
			args:    args{2.0},
			wantErr: true,
		},
		{
			name:   "vertical case",
			fields: newSegment(point{1.0, 2.0}, point{1.0, 5.0}),
			args:   args{2.0},
			wantX:  1.0,
		},
		{
			name:    "normal case",
			fields:  newSegment(point{-10, -4}, point{-6, -2}),
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
				t.Errorf("segment.x() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotX, tt.wantX) {
				t.Errorf("segment.x() = %v, want %v", gotX, tt.wantX)
			}
		})
	}
}

func Test_segment_inBound(t *testing.T) {
	type args struct {
		p point
	}
	tests := []struct {
		name   string
		fields segment
		args   args
		want   bool
	}{
		{
			name:   "outbound x",
			fields: newSegment(point{-10, -4}, point{-6, -2}),
			args:   args{point{-5, -3}},
			want:   false,
		},
		{
			name:   "outbound y",
			fields: newSegment(point{-10, -4}, point{-6, -2}),
			args:   args{point{-7, -1}},
			want:   false,
		},
		{
			name:   "inbound",
			fields: newSegment(point{-10, -4}, point{-6, -2}),
			args:   args{point{-9, -3}},
			want:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := segment{
				startPoint: tt.fields.startPoint,
				endPoint:   tt.fields.endPoint,
				slope:      tt.fields.slope,
				yIntercept: tt.fields.yIntercept,
			}
			if got := s.inBound(tt.args.p); got != tt.want {
				t.Errorf("segment.inBound() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_segment_isSegmentIntersect(t *testing.T) {
	type args struct {
		so segment
	}
	tests := []struct {
		name   string
		fields segment
		args   args
		want   bool
	}{
		{
			name:   "simple parallel",
			fields: newSegment(point{1, 1}, point{10, 1}),
			args:   args{newSegment(point{1, 2}, point{10, 2})},
			want:   false,
		},
		{
			name:   "simple intersect",
			fields: newSegment(point{10, 0}, point{0, 10}),
			args:   args{newSegment(point{0, 0}, point{10, 10})},
			want:   true,
		},
		{
			name:   "simple no intersect",
			fields: newSegment(point{-5, -5}, point{0, 0}),
			args:   args{newSegment(point{1, 1}, point{10, 10})},
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := segment{
				startPoint: tt.fields.startPoint,
				endPoint:   tt.fields.endPoint,
				slope:      tt.fields.slope,
				yIntercept: tt.fields.yIntercept,
			}
			if got := s.isSegmentIntersect(tt.args.so); got != tt.want {
				t.Errorf("segment.isSegmentIntersect() = %v, want %v", got, tt.want)
			}
		})
	}
}
