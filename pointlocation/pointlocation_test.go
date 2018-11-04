package pointlocation

import (
	"reflect"
	"testing"
)

func TestTrapezoidMap(t *testing.T) {
	type args struct {
		ss []segment
	}
	floatZero := 0.0
	floatEight := 8.0
	floatMThree := -3.0
	tests := []struct {
		name   string
		args   args
		wantTr trapezoid
		wantD  dag
	}{
		{
			name: "simple 1",
			args: args{
				ss: []segment{
					newSegment(point{2, 2}, point{4, 3}),
				},
			},
			wantTr: trapezoid{
				leftp:  &point{-3, -3},
				rightp: &point{9, -3},
				top:    &segment{startPoint: point{-3, 8}, endPoint: point{9, 8}, slope: &floatZero, yIntercept: &floatEight},
				bottom: &segment{startPoint: point{-3, -3}, endPoint: point{9, -3}, slope: &floatZero, yIntercept: &floatMThree},
			},
			wantD: dag{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTr, gotD := TrapezoidMap(tt.args.ss)
			if !reflect.DeepEqual(gotTr, tt.wantTr) {
				t.Errorf("TrapezoidMap() gotTr = %v, want %v", gotTr, tt.wantTr)
			}
			if !reflect.DeepEqual(gotD, tt.wantD) {
				t.Errorf("TrapezoidMap() gotD = %v, want %v", gotD, tt.wantD)
			}
		})
	}
}
