package pointlocation

import "testing"

func Test_point_positionBySegment(t *testing.T) {
	type fields struct {
		x float64
		y float64
	}
	type args struct {
		s Segment
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantPos int
		wantErr bool
	}{
		{
			name: "upper",
			fields: fields{
				x: -8,
				y: -2,
			},
			args: args{
				NewSegment(
					Point{x: -10, y: -4},
					Point{x: -6, y: -2},
				),
			},
			wantPos: upper,
			wantErr: false,
		},
		{
			name: "lower",
			fields: fields{
				x: -8,
				y: -10,
			},
			args: args{
				NewSegment(
					Point{x: -10, y: -4},
					Point{x: -6, y: -2},
				),
			},
			wantPos: lower,
			wantErr: false,
		},
		{
			name: "out of bound",
			fields: fields{
				x: -11.68,
				y: -3.99,
			},
			args:    args{NewSegment(Point{x: -10.34, y: 1.83}, Point{x: -6.58, y: 2.01})},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Point{x: tt.fields.x, y: tt.fields.y}
			gotPos, err := p.positionBySegment(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("Point.x:positionBySegment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotPos != tt.wantPos {
				t.Errorf("Point.x:positionBySegment() = %v, want %v", gotPos, tt.wantPos)
			}
		})
	}
}

func Test_point_orientationFromSegment(t *testing.T) {
	type fields struct {
		x float64
		y float64
	}
	tests := []struct {
		name    string
		fields  fields
		args    Segment
		wantPos int
	}{
		{
			name: "vertical to clockwise",
			fields: fields{
				x: 2,
				y: 2,
			},
			args:    NewSegment(Point{x: 0, y: 0}, Point{x: 0, y: 5}),
			wantPos: clockwise,
		},
		{
			name: "vertical to counterclockwise",
			fields: fields{
				x: -2,
				y: -2,
			},
			args:    NewSegment(Point{x: 0, y: 0}, Point{x: 0, y: 5}),
			wantPos: counterclockwise,
		},
		{
			name: "horizontal to clockwise",
			fields: fields{
				x: 3,
				y: 3,
			},
			args:    NewSegment(Point{x: 0, y: 0}, Point{x: 5, y: 0}),
			wantPos: counterclockwise,
		},
		{
			name: "horizontal to counterclockwise",
			fields: fields{
				x: 2,
				y: -2,
			},
			args:    NewSegment(Point{x: 0, y: 0}, Point{x: 5, y: 0}),
			wantPos: clockwise,
		},
		{
			name: "third Point x:vertical clockwise",
			fields: fields{
				x: -1,
				y: 2,
			},
			args:    NewSegment(Point{x: 0, y: 0}, Point{x: -1, y: -1}),
			wantPos: clockwise,
		},
		{
			name: "third Point x:vertical counterclockwise",
			fields: fields{
				x: -1,
				y: -2,
			},
			args:    NewSegment(Point{x: 0, y: 0}, Point{x: -1, y: -1}),
			wantPos: counterclockwise,
		},
		{
			name: "third Point x:horizontal clockwise",
			fields: fields{
				x: 2,
				y: 1,
			},
			args:    NewSegment(Point{x: 0, y: 0}, Point{x: 1, y: 1}),
			wantPos: clockwise,
		},
		{
			name: "third Point x:horizontal counterclockwise",
			fields: fields{
				x: 0,
				y: 1,
			},
			args:    NewSegment(Point{x: 0, y: 0}, Point{x: 1, y: 1}),
			wantPos: counterclockwise,
		},
		{
			name: "third Point x:horizontal and Segment slope < 0 clockwise",
			fields: fields{
				x: 2,
				y: -1,
			},
			args:    NewSegment(Point{x: 0, y: 0}, Point{x: 1, y: -1}),
			wantPos: counterclockwise,
		},
		{
			name: "third Point x:horizontal and Segment slope < 0 counterclockwise",
			fields: fields{
				x: 0,
				y: -1,
			},
			args:    NewSegment(Point{x: 0, y: 0}, Point{x: 1, y: -1}),
			wantPos: clockwise,
		},
		{
			name: "isIntersect fail case 1",
			fields: fields{
				x: 0,
				y: 0,
			},
			args:    NewSegment(Point{x: 10, y: 0}, Point{x: 0, y: 10}),
			wantPos: clockwise,
		},
		{
			name: "isIntersect fail case 2",
			fields: fields{
				x: 10,
				y: 10,
			},
			args:    NewSegment(Point{x: 10, y: 0}, Point{x: 0, y: 10}),
			wantPos: counterclockwise,
		},
		{
			name: "isIntersect fail case 3",
			fields: fields{
				x: 10,
				y: 0,
			},
			args:    NewSegment(Point{x: 0, y: 0}, Point{x: 10, y: 10}),
			wantPos: clockwise,
		},
		{
			name: "isIntersect fail case 4",
			fields: fields{
				x: 0,
				y: 10,
			},
			args:    NewSegment(Point{x: 0, y: 0}, Point{x: 10, y: 10}),
			wantPos: counterclockwise,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Point{
				x: tt.fields.x,
				y: tt.fields.y,
			}
			if gotPos := p.orientationFromSegment(tt.args); gotPos != tt.wantPos {
				t.Errorf("Point.x:orientationFromSegment() = %v, want %v", gotPos, tt.wantPos)
			}
		})
	}
}
