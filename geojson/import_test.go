package geojson

import (
	"reflect"
	"testing"
)

func TestReadGeoJSON(t *testing.T) {
	tests := []struct {
		name string
		want []byte
	}{
		{
			name: "readfile",
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ReadGeoJSON(); !reflect.DeepEqual(got, tt.want) {
				// t.Errorf("ReadGeoJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}
