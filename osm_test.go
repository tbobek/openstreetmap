package openstreetmap

import (
	"io/ioutil"
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
	ways, nodes := parseOsm(buffer)
	if len(ways) == 0 || len(nodes) == 0 {
		t.Errorf("no result")
	}

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
