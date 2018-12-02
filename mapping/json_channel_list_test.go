package mapping_test

import (
	"encoding/json"
	"sort"

	"github.com/aphecetche/pigiron/mapping"
)

func UnmarshalTestChannelList(data []byte) (TestChannelList, error) {
	var r TestChannelList
	err := json.Unmarshal(data, &r)
	if err == nil {
		r.init()
	}
	return r, err
}

func (r *TestChannelList) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type TestChannelList struct {
	DetectionElements []DetectionElement `json:"detection_elements"`
}

func (m *TestChannelList) init() {
	for i, _ := range m.DetectionElements {
		m.DetectionElements[i].init()
	}
}

type DetectionElement struct {
	ID   int   `json:"id"`
	FECs []FEC `json:"manus"`
}

func (de *DetectionElement) init() {
	for i, _ := range de.FECs {
		de.FECs[i].init()
	}
}

func (de DetectionElement) channels(dsid mapping.DualSampaID) ChannelInfoSlice {
	var channels ChannelInfoSlice
	for _, f := range de.FECs {
		if f.ID == dsid {
			for _, c := range f.Channels {
				channels = append(channels, ChannelInfo{dsid, int(c)})
			}
			break
		}
	}
	return channels
}

type FEC struct {
	ID       mapping.DualSampaID `json:"id"`
	Channels []int64             `json:"channels,omitempty"`
}

func (f *FEC) init() {
	if f.Channels == nil || len(f.Channels) == 0 {
		f.Channels = make([]int64, 64)
		for i := 0; i < 64; i++ {
			f.Channels[i] = int64(i)
		}
	}
}

type ChannelInfo struct {
	fecID     mapping.DualSampaID
	channelID int
}
type ChannelInfoSlice []ChannelInfo

func (ci ChannelInfoSlice) Len() int {
	return len(ci)
}
func (ci ChannelInfoSlice) Swap(i, j int) {
	ci[i], ci[j] = ci[j], ci[i]
}
func (ci ChannelInfoSlice) Less(i, j int) bool {
	if ci[i].fecID == ci[j].fecID {
		return ci[i].channelID < ci[j].channelID
	}
	return ci[i].fecID < ci[j].fecID
}

func containSameChannelElements(a ChannelInfoSlice, b ChannelInfoSlice) bool {
	if len(a) != len(b) {
		return false
	}

	if (a == nil) != (b == nil) {
		return false
	}

	sort.Sort(a)
	sort.Sort(b)

	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

// fecID returns a slice of all the front-end card ids of this detection element
func (de DetectionElement) fecIDs() []mapping.DualSampaID {
	var fecids = make([]mapping.DualSampaID, len(de.FECs))
	for i, f := range de.FECs {
		fecids[i] = f.ID
	}
	return fecids
}

// fecID returns a slice of all the front-end card ids of this detection element
func (de DetectionElement) channelIDs() ChannelInfoSlice {
	var channels ChannelInfoSlice
	for _, fec := range de.FECs {
		channels = append(channels, de.channels(fec.ID)...)
	}
	sort.Sort(channels)
	return channels
}

func containSameDualSampaIDs(a []mapping.DualSampaID, b []mapping.DualSampaID) bool {
	if len(a) != len(b) {
		return false
	}

	if (a == nil) != (b == nil) {
		return false
	}

	sort.Slice(a, func(i, j int) bool {
		return int(a[i]) > int(a[j])
	})
	sort.Slice(b, func(i, j int) bool {
		return int(b[i]) > int(b[j])
	})

	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
