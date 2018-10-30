package pointlocation

import "testing"

func Test_point_positionBySegment(t *testing.T) {
	type fields struct {
		x float64
		y float64
	}
	type args struct {
		s segment
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
				newSegment(
					point{-10, -4},
					point{-6, -2},
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
				newSegment(
					point{-10, -4},
					point{-6, -2},
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
			args:    args{newSegment(point{-10.34, 1.83}, point{-6.58, 2.01})},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := point{
				tt.fields.x,
				tt.fields.y,
			}
			gotPos, err := p.positionBySegment(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("point.positionBySegment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotPos != tt.wantPos {
				t.Errorf("point.positionBySegment() = %v, want %v", gotPos, tt.wantPos)
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
		args    segment
		wantPos int
	}{
		{
			name: "vertical to clockwise",
			fields: fields{
				x: 2,
				y: 2,
			},
			args:    newSegment(point{0, 0}, point{0, 5}),
			wantPos: clockwise,
		},
		{
			name: "vertical to counterclockwise",
			fields: fields{
				x: -2,
				y: -2,
			},
			args:    newSegment(point{0, 0}, point{0, 5}),
			wantPos: counterclockwise,
		},
		{
			name: "horizontal to clockwise",
			fields: fields{
				x: 3,
				y: 3,
			},
			args:    newSegment(point{0, 0}, point{5, 0}),
			wantPos: counterclockwise,
		},
		{
			name: "horizontal to counterclockwise",
			fields: fields{
				x: 2,
				y: -2,
			},
			args:    newSegment(point{0, 0}, point{5, 0}),
			wantPos: clockwise,
		},
		{
			name: "third point vertical clockwise",
			fields: fields{
				x: -1,
				y: 2,
			},
			args:    newSegment(point{0, 0}, point{-1, -1}),
			wantPos: clockwise,
		},
		{
			name: "third point vertical counterclockwise",
			fields: fields{
				x: -1,
				y: -2,
			},
			args:    newSegment(point{0, 0}, point{-1, -1}),
			wantPos: counterclockwise,
		},
		{
			name: "third point horizontal clockwise",
			fields: fields{
				x: 2,
				y: 1,
			},
			args:    newSegment(point{0, 0}, point{1, 1}),
			wantPos: clockwise,
		},
		{
			name: "third point horizontal counterclockwise",
			fields: fields{
				x: 0,
				y: 1,
			},
			args:    newSegment(point{0, 0}, point{1, 1}),
			wantPos: counterclockwise,
		},
		{
			name: "third point horizontal and segment slope < 0 clockwise",
			fields: fields{
				x: 2,
				y: -1,
			},
			args:    newSegment(point{0, 0}, point{1, -1}),
			wantPos: counterclockwise,
		},
		{
			name: "third point horizontal and segment slope < 0 counterclockwise",
			fields: fields{
				x: 0,
				y: -1,
			},
			args:    newSegment(point{0, 0}, point{1, -1}),
			wantPos: clockwise,
		},
		{
			name: "isIntersect fail case 1",
			fields: fields{
				x: 0,
				y: 0,
			},
			args:    newSegment(point{10, 0}, point{0, 10}),
			wantPos: clockwise,
		},
		{
			name: "isIntersect fail case 2",
			fields: fields{
				x: 10,
				y: 10,
			},
			args:    newSegment(point{10, 0}, point{0, 10}),
			wantPos: counterclockwise,
		},
		{
			name: "isIntersect fail case 3",
			fields: fields{
				x: 10,
				y: 0,
			},
			args:    newSegment(point{0, 0}, point{10, 10}),
			wantPos: clockwise,
		},
		{
			name: "isIntersect fail case 4",
			fields: fields{
				x: 0,
				y: 10,
			},
			args:    newSegment(point{0, 0}, point{10, 10}),
			wantPos: counterclockwise,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := point{
				x: tt.fields.x,
				y: tt.fields.y,
			}
			if gotPos := p.orientationFromSegment(tt.args); gotPos != tt.wantPos {
				t.Errorf("point.orientationFromSegment() = %v, want %v", gotPos, tt.wantPos)
			}
		})
	}
}
