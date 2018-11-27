package main

import (
	"bytes"
	"encoding/csv"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/kicito/assignment-geo-2/pointlocation"
)

func main() {
	var filename string
	flag.StringVar(&filename, "file", "./test.csv", "csv file to read")
	flag.Parse()

	fmt.Println("clearing image files")
	dir, err := ioutil.ReadDir("./")
	if err != nil {
		log.Fatal(err)
	}
	for _, d := range dir {
		if strings.Contains(d.Name(), "png") {
			os.RemoveAll(path.Join([]string{"./", d.Name()}...))
		}
	}
	_ = os.MkdirAll("./steps", os.ModePerm)

	dir, err = ioutil.ReadDir("./steps")
	if err != nil {
		log.Fatal(err)
	}
	for _, d := range dir {
		if strings.Contains(d.Name(), "png") {
			os.RemoveAll(path.Join([]string{"./steps", d.Name()}...))
		}
	}

	fmt.Println("reading file", filename)

	content, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	csvReader := csv.NewReader(bytes.NewReader(content))

	coorList, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	sList := make([]pointlocation.Segment, len(coorList)-2)
	indexNext := 1
	for coorIndex := range coorList[:len(coorList)-2] {
		latitude, err := strconv.ParseFloat(strings.TrimSpace(coorList[coorIndex][0]), 64)
		logitude, err := strconv.ParseFloat(strings.TrimSpace(coorList[coorIndex][1]), 64)
		latitudeNext, err := strconv.ParseFloat(strings.TrimSpace(coorList[indexNext][0]), 64)
		logitudeNext, err := strconv.ParseFloat(strings.TrimSpace(coorList[indexNext][1]), 64)
		if err != nil {
			log.Fatal(err)
		}

		sList[coorIndex] = pointlocation.NewSegment(
			pointlocation.NewPoint(latitude, logitude),
			pointlocation.NewPoint(latitudeNext, logitudeNext),
			strconv.Itoa(coorIndex),
		)
		indexNext++
	}
	pointX, err := strconv.ParseFloat(strings.TrimSpace(coorList[indexNext][0]), 64)
	pointY, err := strconv.ParseFloat(strings.TrimSpace(coorList[indexNext][1]), 64)
	if err != nil {
		log.Fatal(err)
	}

	point := pointlocation.NewPoint(pointX, pointY)
	fmt.Println("search for ", point)
	// result := geojson.ReadGeoJSON()
	// geodatas := geojson.FeatureCollectionToGeoData(result)

	var pl pointlocation.PointLocation
	// fmt.Println("----------------------")
	// fmt.Println(geodatas[0].Segments)
	// fmt.Println("----------------------")
	if pl, err = pointlocation.NewPointLocation(sList); err != nil {
		log.Fatal(err)
		return
	}

	tr, err := pl.DAG.FindPoint(point)
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Println("----------------------")
	fmt.Println("-------result---------")
	fmt.Println("----------------------")
	fmt.Println(tr)
	fmt.Println("----------------------")
	_, err = pl.PlotTrsWithPoint("result", point)
	if err != nil {
		log.Fatal(err)
	}
}
