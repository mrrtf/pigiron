package mapping_test

import (
	"encoding/json"
	"fmt"

	"github.com/aphecetche/pigiron/mapping"
)

func UnmarshalTestRandomPos(data []byte) (TestRandomPos, error) {
	var r TestRandomPos
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *TestRandomPos) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type TestRandomPos struct {
	Testpositions []Testposition `json:"testpositions"`
}

type Testposition struct {
	De      int64      `json:"de"`
	Bending BoolString `json:"bending"`
	X       float64    `json:"x"`
	Y       float64    `json:"y"`
	Outside BoolString `json:"isoutside,omitempty"`
	Dsid    int64      `json:"dsid"`
	Dsch    int64      `json:"dsch"`
}

func (tp Testposition) String() string {
	return fmt.Sprintf("DE %4d %s X %v Y %v Outside %v fecID %d fecChannel %d",
		tp.De, mapping.PlaneAbbreviation(tp.isBendingPlane()), tp.X, tp.Y, tp.Outside == "true",
		tp.Dsid, tp.Dsch)
}

func (tp Testposition) isBendingPlane() bool {
	return tp.Bending == "true"
}

func (tp Testposition) isOutside() bool {
	return tp.Outside == "true"
}

type BoolString string

const (
	False BoolString = "false"
	True  BoolString = "true"
)
