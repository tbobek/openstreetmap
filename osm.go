package openstreetmap

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type OsmConfig struct {
	InfluxHost string
	InfluxDb   string
}

func GetOsmConfig() OsmConfig {
	osmConf := OsmConfig{
		InfluxHost: os.Getenv("OSM_INFLUX_HOST"),
		InfluxDb:   os.Getenv("OSM_INFLUX_DB"),
	}
	return osmConf
}

// getNodes retrieves street locations from openstreetmap
// p1 is the lower left point on the map
// p2 is the upper right point on the map
func GetNodes(p1 GeoPoint, p2 GeoPoint) []GeoPoint {
	body := makeRequest(p1, p2)
	_, nodes, err := parseOsm(body)
	if err != nil {
		return nil
	}
	return extractPoints(nodes)
}

func makeRequest(p1 GeoPoint, p2 GeoPoint) []byte {
	// api call is like /api/0.6/map?bbox=left, bottom, right, top
	url := "http://api.openstreetmap.org/api/0.6/map?bbox="
	request := fmt.Sprintf("%s%f,%f,%f,%f,%s", url, p1.Lon, p1.Lat, p2.Lon, p2.Lat, "highway=primary")
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

func makeOverpassRequest(p1 GeoPoint, p2 GeoPoint) []byte {
	url := "http://overpass-api.de/api/interpreter"
	data := fmt.Sprintf(`[out:xml][timeout:25];(way["highway"](%f,%f,%f,%f););out body;>;out skel qt;`, p1.Lat, p1.Lon, p2.Lat, p2.Lon)
	body := strings.NewReader(string(data))
	log.Printf("making overpass request: \nurl=%s\ndata=%s", url, data)
	res, err := http.Post(url, "application/string", body)
	if err != nil {
		panic(err)
	}
	respData, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()

	return respData
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
		points = append(points, GeoPoint{lat, lon, "street"})
	}
	return points
}

func parseOsm(input []byte) ([]Way, []Node, error) {
	var osm Osm

	err := xml.Unmarshal(input, &osm)
	if err != nil {
		return nil, nil, err
	}
	return osm.Ways, osm.Nodes, nil
}

func insertPoints(points []GeoPoint) int {
	osmConf := GetOsmConfig()
	// url
	url := fmt.Sprintf("http://%s:8086/write?db=%s", osmConf.InfluxHost, osmConf.InfluxDb)
	counter := 0
	for _, p := range points {

		pj := fmt.Sprintf("street_points,value=%d,name=test lat=%f,lon=%f", 1, p.Lat, p.Lon)
		//log.Printf("influx POST request: %s, body: %s", url, string(pj))
		body := strings.NewReader(string(pj))
		res, err := http.Post(url, "application/string", body)

		if err != nil {
			panic(err)
		}
		//data, _ := ioutil.ReadAll(res.Body)
		res.Body.Close()
		//log.Printf("response from influx POST insert: %s", data)
		counter += 1
	}
	return counter
}
