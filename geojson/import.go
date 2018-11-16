package geojson

import (
	"io/ioutil"
	"log"

	"github.com/kicito/assignment-geo-2/pointlocation"

	"github.com/paulmach/go.geojson"
)

const datadir = "./raw_data/dmk_dump.geojson"

type GeoData struct {
	Segments     []pointlocation.Segment
	Community    string
	Municipality string
}

func ReadGeoJSON() *geojson.FeatureCollection {
	content, err := ioutil.ReadFile(datadir)
	if err != nil {
		log.Fatalf("unable to read file %v", datadir)
	}
	result, err := geojson.UnmarshalFeatureCollection(content)
	if err != nil {
		log.Fatalf("parse geojson %v", err)
	}
	return result
}

func FeatureCollectionToGeoData(f *geojson.FeatureCollection) (results []GeoData) {
	results = make([]GeoData, len(f.Features))
	for featureIndex := range f.Features {
		result := GeoData{
			Segments: pointlocation.GeometryToSegments(f.Features[featureIndex].Geometry),
		}
		if data, ok := f.Features[featureIndex].Properties["OpstilNav"].(string); ok {
			result.Community = data
		}
		if data, ok := f.Features[featureIndex].Properties["KommuneNav"].(string); ok {
			result.Municipality = data
		}

		results[featureIndex] = result
	}
	return
}
