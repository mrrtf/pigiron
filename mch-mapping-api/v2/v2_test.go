package v2

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

const input = `{"padlist":[{"deid":100,"padid":1},{"deid":100,"padid":2},{"deid":100,"padid":3},{"deid":1025,"padid":1},{"deid":1025,"padid":2},{"deid":1025,"padid":3},{"deid":503,"padid":1589}]}`

func TestUnmarshal(t *testing.T) {
	r := Padlist{}
	err := json.Unmarshal([]byte(input), &r)
	if err != nil {
		t.Errorf("Could not Unmarshal : %s", err.Error())
	} else {
		fmt.Println("input string=", input)
		fmt.Printf("output struct=%+v\n", r)
	}
}

func TestDecoder(t *testing.T) {
	r := Padlist{}
	decoder := json.NewDecoder(strings.NewReader(input))
	err := decoder.Decode(&r)
	if err != nil {
		t.Errorf("Could not Decode : %s", err.Error())
	} else {
		fmt.Println(r)
	}
}

type Vertices struct {
	Vertices []Vertex
}

func TestMarshalMap(t *testing.T) {
	m := make(map[string]Vertices)
	var v []Vertex
	v = append(v, Vertex{12, 34})
	m["1-2-3"] = Vertices{v}
	jsonString, err := json.Marshal(m)
	if err != nil {
		t.Error("could not marshal map")
	}
	const expected = `{"1-2-3":{"Vertices":[{"x":12,"y":34}]}}`
	got := string(jsonString)
	if expected != got {
		t.Errorf("got %s instead of expected %s", got, expected)
	}

}
