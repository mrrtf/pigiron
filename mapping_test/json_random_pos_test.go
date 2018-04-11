package mapping_test

import "encoding/json"

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
