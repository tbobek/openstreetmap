package openstreetmap

import "encoding/xml"

type GeoPoint struct {
	Lat  float64 `json:"lat"`
	Lon  float64 `json:"lon"`
	Type string  `json:"type"`
}

type Osm struct {
	XMLName xml.Name `xml:"osm"`
	Nodes   []Node   `xml:"node"`
	Ways    []Way    `xml:"way"`
}

type Way struct {
	XMLName xml.Name `xml:"way"`
	Nds     []Nd     `xml:"nd"`
	Tags    []Tag    `xml:"tag"`
}

type Tag struct {
	XMLName xml.Name `xml:"tag"`
	Key     string   `xml:"k,attr"`
	Value   string   `xml:"v,attr"`
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
