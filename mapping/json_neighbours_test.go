package mapping_test

import (
	"encoding/json"

	"github.com/aphecetche/pigiron/mapping"
)

type TestNeighbourStruct struct {
	Neighbours []struct {
		Deid int `json:"deid"`
		Ds   []struct {
			ID       mapping.DualSampaID `json:"id"`
			Channels []struct {
				Ch  int `json:"ch"`
				Nei []struct {
					Dsid mapping.DualSampaID `json:"dsid"`
					Dsch int                 `json:"dsch"`
				} `json:"nei"`
			} `json:"channels"`
		} `json:"ds"`
	} `json:"neighbours"`
}

func UnmarshalTestNeighbours(data []byte) (TestNeighbourStruct, error) {
	var n TestNeighbourStruct
	err := json.Unmarshal(data, &n)
	return n, err
}

func (n *TestNeighbourStruct) Marshal() ([]byte, error) {
	return json.Marshal(n)
}
