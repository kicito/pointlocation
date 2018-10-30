package pointlocation

import (
	"reflect"
	"testing"
)

func Test_trapezoid_topSegment(t *testing.T) {
	type fields struct {
		leftp  *point
		rightp *point
		top    *segment
		bottom *segment
	}
	bounderyTopSegment := newSegment(point{0, 5}, point{5, 5})
	bounderyBotSegment := newSegment(point{0, 0}, point{5, 0})
	tests := []struct {
		name   string
		fields fields
		want   segment
	}{
		{
			name: "boundery",
			fields: fields{
				leftp:  &point{0, 0},
				rightp: &point{5, 0},
				top:    &bounderyTopSegment,
				bottom: &bounderyBotSegment,
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
		leftp  *point
		rightp *point
		top    *segment
		bottom *segment
	}
	bounderyTopSegment := newSegment(point{0, 5}, point{5, 5})
	bounderyBotSegment := newSegment(point{0, 0}, point{5, 0})
	tests := []struct {
		name   string
		fields fields
		want   segment
	}{
		{
			name: "boundery",
			fields: fields{
				leftp:  &point{0, 0},
				rightp: &point{5, 0},
				top:    &bounderyTopSegment,
				bottom: &bounderyBotSegment,
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
		leftp  *point
		rightp *point
		top    *segment
		bottom *segment
	}
	bounderyTopSegment := newSegment(point{0, 5}, point{5, 5})
	bounderyBotSegment := newSegment(point{0, 0}, point{5, 0})
	tests := []struct {
		name   string
		fields fields
		want   segment
	}{
		{
			name: "boundery",
			fields: fields{
				leftp:  &point{0, 0},
				rightp: &point{5, 0},
				top:    &bounderyTopSegment,
				bottom: &bounderyBotSegment,
			},
			want: newSegment(point{0, 5}, point{0, 0}),
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
		leftp  *point
		rightp *point
		top    *segment
		bottom *segment
	}
	bounderyTopSegment := newSegment(point{0, 5}, point{5, 5})
	bounderyBotSegment := newSegment(point{0, 0}, point{5, 0})
	tests := []struct {
		name   string
		fields fields
		want   segment
	}{
		{
			name: "boundery",
			fields: fields{
				leftp:  &point{0, 0},
				rightp: &point{5, 0},
				top:    &bounderyTopSegment,
				bottom: &bounderyBotSegment,
			},
			want: newSegment(point{5, 5}, point{5, 0}),
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
		leftp  *point
		rightp *point
		top    *segment
		bottom *segment
	}
	type args struct {
		p point
	}
	bounderyTopSegment := newSegment(point{0, 5}, point{5, 5})
	bounderyBotSegment := newSegment(point{0, 0}, point{5, 0})
	tests := []struct {
		name   string
		fields fields
		args   args
		wantB  bool
	}{
		{
			name: "boundery",
			fields: fields{
				leftp:  &point{0, 0},
				rightp: &point{5, 0},
				top:    &bounderyTopSegment,
				bottom: &bounderyBotSegment,
			},
			args:  args{point{1, 1}},
			wantB: true,
		},
		{
			name: "boundery 2",
			fields: fields{
				leftp:  &point{0, 0},
				rightp: &point{5, 0},
				top:    &bounderyTopSegment,
				bottom: &bounderyBotSegment,
			},
			args:  args{point{2, 2}},
			wantB: true,
		},
		{
			name: "out of boundery",
			fields: fields{
				leftp:  &point{0, 0},
				rightp: &point{5, 0},
				top:    &bounderyTopSegment,
				bottom: &bounderyBotSegment,
			},
			args:  args{point{6, 6}},
			wantB: false,
		},
		{
			name: "out of boundery 2",
			fields: fields{
				leftp:  &point{0, 0},
				rightp: &point{5, 0},
				top:    &bounderyTopSegment,
				bottom: &bounderyBotSegment,
			},
			args:  args{point{-1, -1}},
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

func Test_trapezoid_addSegmentInside(t *testing.T) {
	type fields struct {
		leftp       *point
		rightp      *point
		top         *segment
		bottom      *segment
		upperLeftN  *trapezoid
		lowerLeftN  *trapezoid
		upperRightN *trapezoid
		lowerRightN *trapezoid
	}
	type args struct {
		s segment
	}
	lp := point{0, 0}
	rp := point{5, 0}
	inputStart := point{2, 2}
	inputEnd := point{4, 3}
	bounderyTopSegment := newSegment(lp, point{5, 5})
	bounderyBotSegment := newSegment(point{0, 0}, rp)
	input := newSegment(inputStart, inputEnd)

	lt := trapezoid{
		leftp:  &lp,
		rightp: &inputStart,
		top:    &bounderyTopSegment,
		bottom: &bounderyBotSegment,
	}
	ut := trapezoid{
		leftp:  &inputStart,
		rightp: &inputEnd,
		top:    &bounderyTopSegment,
		bottom: &input,
	}
	rt := trapezoid{
		leftp:  &inputEnd,
		rightp: &rp,
		top:    &bounderyTopSegment,
		bottom: &bounderyBotSegment,
	}
	bt := trapezoid{
		leftp:  &inputStart,
		rightp: &inputEnd,
		top:    &input,
		bottom: &bounderyBotSegment,
	}
	// assign neighbor
	lt.upperRightN = &ut
	lt.lowerRightN = &bt

	ut.upperLeftN = &lt
	ut.upperRightN = &rt

	rt.upperLeftN = &ut
	rt.lowerLeftN = &bt

	bt.upperLeftN = &lt
	bt.upperRightN = &rt
	tests := []struct {
		name   string
		fields fields
		args   args
		wantTs []trapezoid
	}{
		{
			name: "first add",
			fields: fields{
				leftp:  &point{0, 0},
				rightp: &point{5, 0},
				top:    &bounderyTopSegment,
				bottom: &bounderyBotSegment,
			},
			args:   args{input},
			wantTs: []trapezoid{lt, ut, rt, bt},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := trapezoid{
				leftp:       tt.fields.leftp,
				rightp:      tt.fields.rightp,
				top:         tt.fields.top,
				bottom:      tt.fields.bottom,
				upperLeftN:  tt.fields.upperLeftN,
				lowerLeftN:  tt.fields.lowerLeftN,
				upperRightN: tt.fields.upperRightN,
				lowerRightN: tt.fields.lowerRightN,
			}
			if gotTs := tr.addSegmentInside(tt.args.s); !reflect.DeepEqual(gotTs, tt.wantTs) {
				t.Errorf("trapezoid.addSegmentInside() = %v, want %v", gotTs, tt.wantTs)
			}
		})
	}
}
