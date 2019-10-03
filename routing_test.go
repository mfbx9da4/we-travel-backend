package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

type TestData struct {
	GeoJSON      GeoJson      `json:"geojson"`
	From         Coordinate   `json:"from"`
	To           Coordinate   `json:"to"`
	Distance     float64      `json:"distance"`
	ShortestPath []Coordinate `json:"shortestPath"`
}

func LoadFile(filename string) TestData {
	jsonFile, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var testData TestData
	json.Unmarshal(byteValue, &testData)
	defer jsonFile.Close()
	return testData
}

func equal(coords1 []Coordinate, coords2 []Coordinate) bool {
	if len(coords1) != len(coords2) {
		return false
	}
	for i := 0; i < len(coords1); i++ {
		if coords1[i] != coords2[i] {
			return false
		}
	}
	return true
}

func TestSum(t *testing.T) {
	filenames := []string{
		"./tests/two_points.json",
		"./tests/straight_line.json",
		"./tests/three_roads.json",
		"./tests/shortest_hops.json",
		"./tests/intersection.json",
	}

	for i := 0; i < len(filenames); i++ {
		fmt.Println("===>", filenames[i])
		var expected = LoadFile(filenames[i])
		var graph = createGraph(expected.GeoJSON)
		graph.String()
		var route = graph.CalculatePath(expected.From, expected.To)
		var path = getCoordinates(route.Path)
		if !equal(path, expected.ShortestPath) {
			t.Errorf("Incorrect path. Got: %v Expected: %v", path, expected.ShortestPath)
		}
		if route.Distance != expected.Distance {
			t.Errorf("Incorrect distance. Got: %f Want: %f", route.Distance, expected.Distance)
		}
	}

}
