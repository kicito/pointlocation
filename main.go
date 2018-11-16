package main

import (
	"fmt"
	"log"

	"github.com/kicito/assignment-geo-2/pointlocation"
)

func main() {

	// result := geojson.ReadGeoJSON()
	// _ = geojson.FeatureCollectionToGeoData(result)
	// geodatas := geojson.FeatureCollectionToGeoData(result)

	geodatas := []pointlocation.Segment{
		pointlocation.NewSegment(pointlocation.NewPoint(-8, 2), pointlocation.NewPoint(-5, 1)),
		pointlocation.NewSegment(pointlocation.NewPoint(-8, 2), pointlocation.NewPoint(-6, 4)),
		pointlocation.NewSegment(pointlocation.NewPoint(-5, 1), pointlocation.NewPoint(0, 4)),
		pointlocation.NewSegment(pointlocation.NewPoint(2, 7), pointlocation.NewPoint(0, 4)),
		pointlocation.NewSegment(pointlocation.NewPoint(-10, 6), pointlocation.NewPoint(2, 7)),
		pointlocation.NewSegment(pointlocation.NewPoint(-10, 6), pointlocation.NewPoint(-4, 5)),
		pointlocation.NewSegment(pointlocation.NewPoint(-6, 4), pointlocation.NewPoint(-4, 5)),
	}
	var pl pointlocation.PointLocation
	var err error
	// fmt.Println("----------------------")
	// fmt.Println(geodatas[0].Segments)
	// fmt.Println("----------------------")
	// if pl, err = pointlocation.NewPointLocation(geodatas[0].Segments); err != nil {
	// 	log.Fatal(err)
	// }
	if pl, err = pointlocation.NewPointLocation(geodatas); err != nil {
		log.Fatal(err)
	}

	// for trIndex := range pl.Trs {
	// 	fmt.Println("----------------------")
	// 	fmt.Println(pl.Trs[trIndex])
	// 	fmt.Println("----------------------")
	// }

	// // // tr := pl.DAG.FindPoint(pointlocation.NewPoint(5, 3))
	tr := pl.DAG.FindPoint(pointlocation.NewPoint(-7, 2))
	// tr := pl.DAG.FindPoint(pointlocation.NewPoint(11.677410006523132, 56.20442736478735))

	fmt.Println("----------------------")
	fmt.Println("-------result---------")
	fmt.Println("----------------------")
	fmt.Println(tr)
	fmt.Println("----------------------")

}
