package pointlocation

import (
	"fmt"
	"testing"
)

func TestPointLocation_findIntersection(t *testing.T) {
	dummyTrs1.assignRightNeighborToTrapezoid(&dummyTrs2)
	pl, err := NewPointLocation(
		[]Segment{
			NewSegment(Point{x: -11, y: -6}, Point{x: -6, y: -3}),
			NewSegment(Point{x: -6, y: -3}, Point{x: 2, y: -3}),
		},
	)
	if err != nil {
		t.Error(err)
	}
	type args struct {
		s Segment
	}
	tests := []struct {
		name    string
		fields  PointLocation
		args    args
		wantTrs []*trapezoid
		wantErr bool
	}{
		{
			name:   "simple bottom intersect",
			fields: pl,
			args: args{
				NewSegment(Point{x: -8, y: -8}, Point{x: -2, y: -6}),
			},
			wantTrs: []*trapezoid{
				&dummyTrs1,
				&dummyTrs2,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pl = tt.fields
			gotTrs, err := pl.findIntersection(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("PointLocation.findIntersection() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for trIndex := range tt.wantTrs {

				pl.PlotTrsWithSegment(fmt.Sprintf("%v", trIndex), tt.args.s)
				PlotTr(fmt.Sprintf("got %v", trIndex), gotTrs[trIndex])

				if !gotTrs[trIndex].equalTrapezoid(*tt.wantTrs[trIndex], false) {
					t.Errorf("PointLocation.findIntersection() gotTrs = %v, wantTrs %v", gotTrs[trIndex], *tt.wantTrs[trIndex])
				}
			}
		})
	}
}
