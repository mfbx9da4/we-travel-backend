package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Property struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Highway  string `json:"highway"`
	Access   string `json:"access"`
	Lit      string `json:"lit"`
	Sidewalk string `json:"sidewalk"`
}

type Coordinate = [2]float64

type Geometry struct {
	Type        string       `json:"type"`
	Coordinates []Coordinate `json:"coordinates"`
}

type Feature struct {
	Type       string   `json:"type"`
	ID         string   `json:"id"`
	Properties Property `json:"properties"`
	Geometry   Geometry `json:"geometry"`
}

type GeoJson struct {
	Type     string    `json:"type"`
	Features []Feature `json:"features"`
}

var graph Graph

func loadGeoJSON() {
	// geoJsonDownloadLink := "https://ucb7e1be7e59700bb615fc052d06.dl.dropboxusercontent.com/cd/0/get/ApeoomlSroMi4LLrd88j2O1YyfZcz-fnOcR-BMu7Ca3F-aclMpnyLmlzJPZtgze6QSfiGh_SZAcCl-TzGSrcNR14iFsaOBl-vs7CsUzWnL6UbsaH7V_CR-apDThjG8fUH78/file?dl=1DownloadLink"
	// resp, err := http.Get(geoJsonDownloadLink)
	// if err != nil {
	// 	// handle error
	// }
	// defer resp.Body.Close()
	// jsonFile := resp.Body
	// GeoJSON for central london around highbury islington
	jsonFile, err := os.Open("./data/central.geojson")
	// GeoJSON for Greater London
	// from http://download.geofabrik.de/europe/great-britain/england/greater-london.html
	// jsonFile, err := os.Open("./data/greater-london-latest.geojson")
	// fetch from
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened geojson")
	byteValue, _ := ioutil.ReadAll(jsonFile)
	fmt.Println("Successfully ReadAll geojson")
	var geojson GeoJson
	json.Unmarshal(byteValue, &geojson)

	for i := 0; i < len(geojson.Features); i++ {
		isLineString := geojson.Features[i].Geometry.Type == "LineString"
		isHighway := geojson.Features[i].Properties.Highway != ""
		hasSidewalk := geojson.Features[i].Properties.Sidewalk != "" || geojson.Features[i].Properties.Sidewalk != "none"
		isPath := geojson.Features[i].Properties.Highway == "path"
		isValidPath := geojson.Features[i].Properties.Highway == "path" && (geojson.Features[i].Properties.Access == "no" || geojson.Features[i].Properties.Access == "private")
		isNotPathOrIsValidPath := !isPath || isValidPath
		isLit := geojson.Features[i].Properties.Lit != "" || geojson.Features[i].Properties.Lit == "yes"
		if isHighway && isLineString && hasSidewalk && isNotPathOrIsValidPath && isLit {
			var feature = geojson.Features[i]
			var prev *Node
			for j := 0; j < len(feature.Geometry.Coordinates); j++ {
				var coords = feature.Geometry.Coordinates[j]
				node := CreateNode(coords)
				graph.AddNode(&node)
				if j != 0 {
					graph.AddEdge(&node, prev)
					graph.AddEdge(prev, &node)
				}
				prev = &node
			}
		}
	}

	defer jsonFile.Close()
	fmt.Println("geojson Graph created with", len(graph.nodes), "nodes")
}

func calculatePath(startCoords Coordinate, endCoords Coordinate) Route {
	nodeStart := graph.FindNode(startCoords)
	nodeEnd := graph.FindNode(endCoords)
	pathFound := graph.FindPath(graph.nodes[nodeStart], graph.nodes[nodeEnd])
	fmt.Println(pathFound)
	return pathFound
}
