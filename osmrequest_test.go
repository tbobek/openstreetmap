package openstreetmap

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func TestGetNodes(t *testing.T) {
	fmt.Println("vim-go")
	corner1 := GeoPoint{Lat: 50.84, Lon: 6.085}
	corner2 := GeoPoint{Lat: 50.87, Lon: 6.095}
	body := makeRequest(corner1, corner2)
	ioutil.WriteFile("response.xml", body, 0777)
	if len(body) < 1 {
		t.Errorf("no nodes returned")
	}
	// insert nods in influxdb
}

/*
func TestInsertNodes(t *testing.T) {
	err := InsertNodes(nodes)
	if err != nil {
		log.Panic("could not insert nodes", err)
	}
	fmt.Println("finished")
}
*/
