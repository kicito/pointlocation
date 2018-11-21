package pointlocation

import "testing"

func Test_plotTr(t *testing.T) {
	type args struct {
		tr trapezoid
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "simple",
			args: args{
				tr: dummyTrs1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			PlotTr("", tt.args.tr)
		})
	}
}

func Test_Plot_PointLocation(t *testing.T) {

	plTest1, _ := NewPointLocation(
		[]Segment{
			NewSegment(Point{x: -11, y: -6}, Point{x: -6, y: -3}),
			NewSegment(Point{x: -6, y: -3}, Point{x: 2, y: -3}),
		},
	)

	err := plTest1.PlotTrs("")
	if err != nil {
		t.Error(err)
	}
}

func Test_Plot_PointLocationWithSegment(t *testing.T) {
	s := NewSegment(Point{x: -8, y: -8}, Point{x: -2, y: 6})

	var plTest1, _ = NewPointLocation(
		[]Segment{
			NewSegment(Point{x: -11, y: -6}, Point{x: -6, y: -3}),
			NewSegment(Point{x: -6, y: -3}, Point{x: 2, y: -3}),
		},
	)

	err := plTest1.PlotTrsWithSegment("", s)
	if err != nil {
		t.Error(err)
	}
}
