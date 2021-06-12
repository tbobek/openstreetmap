package openstreetmap

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestGetNodes(t *testing.T) {
	corner1 := GeoPoint{Lat: 50.84, Lon: 6.085}
	corner2 := GeoPoint{Lat: 50.87, Lon: 6.095}
	body := makeRequest(corner1, corner2)
	ioutil.WriteFile("response.xml", body, 0777)
	if len(body) < 1 {
		t.Errorf("no nodes returned")
	}
	// insert nods in influxdb
}

func TestGetOverpassNodes(t *testing.T) {
	corner1 := GeoPoint{Lat: 50.84, Lon: 6.07}
	corner2 := GeoPoint{Lat: 50.87, Lon: 6.09}
	body := makeOverpassRequest(corner1, corner2)
	ioutil.WriteFile("response.xml", body, 0777)
	if len(body) < 1 {
		t.Errorf("no nodes returned")
	}
}

func TestGetOverpassNodes2(t *testing.T) {
	corner1 := GeoPoint{Lat: 50.79147736514272, Lon: 6.05222225189209}
	corner2 := GeoPoint{Lat: 50.81550721541137, Lon: 6.097497940063477}
	body := makeOverpassRequest(corner1, corner2)
	ioutil.WriteFile("response.xml", body, 0777)
	if len(body) < 1 {
		t.Errorf("no nodes returned")
	}
}

func TestParseXml(t *testing.T) {
	buffer, err := ioutil.ReadFile("response.xml")
	if err != nil {
		t.Errorf("could not read file response.xml")
	}
	if len(buffer) == 0 {
		t.Errorf("no contents")
	}
	ways, nodes, err := parseOsm(buffer)
	if err != nil {
		log.Panic("could not parse")
	}
	if len(ways) == 0 || len(nodes) == 0 {
		t.Errorf("no result")
	}
}

func TestExtractPoints(t *testing.T) {
	ways, nodes, err := getWaysNodes("response.xml")
	if err != nil {
		panic(err)
	}
	fmt.Printf("#ways returned: %d\n", len(ways))
	fmt.Printf("#nodes returned: %d\n", len(nodes))
	points := extractPoints(nodes)
	if len(points) == 0 {
		t.Errorf("no points extracted")
	}
}

func TestTags(t *testing.T) {
	ways, _, err := getWaysNodes("response.xml")
	if err != nil {
		panic(err)
	}
	if len(ways) == 0 {
		t.Error("no ways present")
	} else {
		for i, way := range ways {
			for j, osmtag := range way.Tags {
				if osmtag.Key == "" {
					t.Errorf("way #%d, tag #%d: Key is empty", i, j)
				}
				if osmtag.Value == "" {
					t.Errorf("way #%d, tag #%d: Value is empty", i, j)
				}
			}
		}
	}
}

func getWaysNodes(filename string) ([]Way, []Node, error) {
	buffer, err := ioutil.ReadFile(filename)
	if err != nil {
		panic("could not read file response.xml")
	}
	if len(buffer) == 0 {
		panic("no contents")
	}
	return parseOsm(buffer)
}

func TestInsertPoints(t *testing.T) {
	os.Setenv("OSM_INFLUX_HOST", "192.168.178.34")
	os.Setenv("OSM_INFLUX_DB", "street_data")
	ways, nodes, err := getWaysNodes("response.xml")
	if err != nil {
		panic(err)
	}
	fmt.Printf("#ways returned: %d\n", len(ways))
	fmt.Printf("#nodes returned: %d\n", len(nodes))
	points := extractPoints(nodes)
	if len(points) == 0 {
		t.Errorf("no points extracted")
	}
	n := insertPoints(points)
	if n == 0 {
		t.Error("no points written to influx")
	}
}
