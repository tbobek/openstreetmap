package openstreetmap

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type GeoPoint struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type Osm struct {
	XMLName xml.Name `xml:"osm"`
	Nodes   []Node   `xml:"node"`
	Way     []Way    `xml:"way"`
}

type Way struct {
	XMLName xml.Name `xml:"way"`
	Nds     []Nd     `xml:"nd"`
	Tags    []Tag    `xml:"tag"`
}

type Tag struct {
	XMLName xml.Name `xml:"tag"`
	Key     string   `xml:"k"`
	Value   string   `xml:"v"`
}

type Nd struct {
	XMLName xml.Name `xml:"nd"`
	Ref     string   `xml:"ref,attr"`
}

type Node struct {
	XMLName xml.Name `xml:"node"`
	Id      string   `xml:"id,attr"`
	Lat     string   `xml:"lat,attr"`
	Lon     string   `xml:"lon,attr"`
}

// getNodes retrieves street locations from openstreetmap
// p1 is the lower left point on the map
// p2 is the upper right point on the map
func GetNodes(p1 GeoPoint, p2 GeoPoint) []GeoPoint {
	body := makeRequest(p1, p2)
	_, nodes := parseOsm(body)

	return extractPoints(nodes)
}

func makeRequest(p1 GeoPoint, p2 GeoPoint) []byte {
	// api call is like /api/0.6/map?bbox=left, bottom, right, top
	url := "http://api.openstreetmap.org/api/0.6/map?bbox="
	request := fmt.Sprintf("%s%f,%f,%f,%f", url, p1.Lon, p1.Lat, p2.Lon, p2.Lat)
	log.Printf("making request: %s", request)
	resp, err := http.Get(request)
	if err != nil {
		log.Panic("error during request", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Panic("error reading response body", err)
	}
	log.Printf("response (first 200 chars, total length: %d):\n%s", len(body), body[0:200])
	return body
}

func extractPoints(nodes []Node) []GeoPoint {
	points := make([]GeoPoint, 0)
	for _, node := range nodes {
		lat, err := strconv.ParseFloat(node.Lat, 64)
		if err != nil {
			log.Panic(err)
		}
		lon, err := strconv.ParseFloat(node.Lon, 64)
		if err != nil {
			log.Panic(err)
		}
		points = append(points, GeoPoint{lat, lon})
	}
	return points
}

func parseOsm(input []byte) ([]Way, []Node) {
	var way Way
	nodes := make([]Node, 0)
	ways := make([]Way, 0)
	xml.Unmarshal(input, &way)
	return ways, nodes
}
