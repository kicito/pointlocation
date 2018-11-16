package pointlocation

import (
	"log"

	geojson "github.com/paulmach/go.geojson"
)

func GeometryToSegments(f *geojson.Geometry) (results []Segment) {
	results = make([]Segment, 0)
	if f.IsPolygon() {
		var startPoint, endPoint Point

		for polygonIndex := range f.Polygon {

			for pointIndex := range f.Polygon[polygonIndex][:len(f.Polygon[polygonIndex])-1] {
				startPoint = Point{
					x: f.Polygon[polygonIndex][pointIndex][0],
					y: f.Polygon[polygonIndex][pointIndex][1],
				}
				endPoint = Point{
					x: f.Polygon[polygonIndex][pointIndex+1][0],
					y: f.Polygon[polygonIndex][pointIndex+1][1],
				}
				results = append(results, NewSegment(startPoint, endPoint))
				startPoint = endPoint
			}
		}
	} else if f.IsMultiPolygon() {
		var startPoint, endPoint Point

		for multiPolygonIndex := range f.MultiPolygon {

			for polygonIndex := range f.MultiPolygon[multiPolygonIndex] {

				for pointIndex := range f.MultiPolygon[multiPolygonIndex][polygonIndex][:len(f.MultiPolygon[multiPolygonIndex][polygonIndex])-1] {
					startPoint = Point{
						x: f.MultiPolygon[multiPolygonIndex][polygonIndex][pointIndex][0],
						y: f.MultiPolygon[multiPolygonIndex][polygonIndex][pointIndex][1],
					}
					endPoint = Point{
						x: f.MultiPolygon[multiPolygonIndex][polygonIndex][pointIndex+1][0],
						y: f.MultiPolygon[multiPolygonIndex][polygonIndex][pointIndex+1][1],
					}

					results = append(results, NewSegment(startPoint, endPoint))
					startPoint = endPoint
				}
			}
		}
	} else if f.IsLineString() {

		var startPoint, endPoint Point

		for lineStringIndex := range f.LineString[:len(f.LineString)-1] {

			startPoint = Point{
				x: f.LineString[lineStringIndex][0],
				y: f.LineString[lineStringIndex][1],
			}
			endPoint = Point{
				x: f.LineString[lineStringIndex+1][0],
				y: f.LineString[lineStringIndex+1][1],
			}
			results = append(results, NewSegment(startPoint, endPoint))
			startPoint = endPoint
		}
	} else {
		log.Fatalf("Geometry type not support %v", f)

	}

	return
}
