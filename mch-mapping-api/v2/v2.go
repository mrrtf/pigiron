package v2

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mrrtf/pigiron/mapping"
)

type Padlist struct {
	Padlist []PadRef `json:"padlist"`
}

type PadRef struct {
	DeId  int `json:"deid"`
	PadId int `json:"padid"`
}

func DualSampas(w http.ResponseWriter, r *http.Request, deid int, bending bool) {
	cseg := mapping.NewCathodeSegmentation(mapping.DEID(deid), bending)
	jsonDualSampas(w, cseg, bending)
}

func GetDualSampaPads(w http.ResponseWriter, r *http.Request, deid int, dsid int) {
	bending := true
	if dsid > 1024 {
		bending = false
	}
	cseg := mapping.NewCathodeSegmentation(mapping.DEID(deid), bending)
	jsonDualSampaPads(w, cseg, dsid)
}

func DeGeo(w http.ResponseWriter, r *http.Request, deid int, bending bool) {
	cseg := mapping.NewCathodeSegmentation(mapping.DEID(deid), bending)
	jsonDEGeo(w, cseg, bending)
}

func PadList(w http.ResponseWriter, r *http.Request) {
	// expected form of the request is :
	// "padlist": [ "X-Y-Z", "a-b-c" ... ]
	// where
	// X=DeId
	// Y=DsId
	// Z=DsCh
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Content-type", "application/json")
	var padlist Padlist
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&padlist)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	jsonPadList(w, padlist.Padlist)
}
