package v2

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

// type RequestPadId struct {
// 	Des []struct {
// 		ID     int   `json:"id"`
// 		Padids []int `json:"padids"`
// 	} `json:"des"`
// }

const input = `{"des":[{"id":102,"padids":[1,2,3]},{"id":1025,"padids":[10,20,33]}]}`

func TestUnmarshal(t *testing.T) {
	r := RequestPadId{}
	err := json.Unmarshal([]byte(input), &r)
	if err != nil {
		t.Errorf("Could not Unmarshal : %s", err.Error())
	} else {
		fmt.Println("input string=", input)
		fmt.Printf("output struct=%+v\n", r)
	}
}

func TestDecoder(t *testing.T) {
	r := RequestPadId{}
	decoder := json.NewDecoder(strings.NewReader(input))
	err := decoder.Decode(&r)
	if err != nil {
		t.Errorf("Could not Decode : %s", err.Error())
	} else {
		fmt.Println(r)
	}
}
