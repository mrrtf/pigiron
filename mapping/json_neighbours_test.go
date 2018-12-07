package mapping_test

import (
	"encoding/json"

	"github.com/aphecetche/pigiron/mapping"
)

type TestNeighbourStruct struct {
	Neighbours []struct {
		Deid mapping.DEID `json:"deid"`
		Ds   []struct {
			ID       mapping.DualSampaID `json:"id"`
			Channels []struct {
				Ch  mapping.DualSampaChannelID `json:"ch"`
				Nei []struct {
					Dsid mapping.DualSampaID        `json:"dsid"`
					Dsch mapping.DualSampaChannelID `json:"dsch"`
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
