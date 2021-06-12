package openstreetmap

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

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
