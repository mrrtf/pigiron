package mapping_test

import (
	"encoding/json"
	"sort"
)

func UnmarshalMch(data []byte) (Mch, error) {
	var r Mch
	err := json.Unmarshal(data, &r)
	if err == nil {
		r.init()
	}
	return r, err
}

func (r *Mch) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type Mch struct {
	DetectionElements []DetectionElement `json:"detection_elements"`
}

func (m *Mch) init() {
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

func (de DetectionElement) channels(dualSampaID int) ChannelInfoSlice {
	var channels ChannelInfoSlice
	for _, f := range de.FECs {
		if f.ID == dualSampaID {
			for _, c := range f.Channels {
				channels = append(channels, ChannelInfo{dualSampaID, int(c)})
			}
			break
		}
	}
	return channels
}

type FEC struct {
	ID       int     `json:"id"`
	Channels []int64 `json:"channels,omitempty"`
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
	fecID     int
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
func (de DetectionElement) fecIDs() []int {
	var fecids = make([]int, len(de.FECs))
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

func containSameIntElements(a []int, b []int) bool {
	if len(a) != len(b) {
		return false
	}

	if (a == nil) != (b == nil) {
		return false
	}

	sort.Ints(a)
	sort.Ints(b)

	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
