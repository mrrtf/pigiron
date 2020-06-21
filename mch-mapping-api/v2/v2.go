package v2

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mrrtf/pigiron/mapping"
)

type PadRef struct {
	DeId  int
	PadId int
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

type RequestPadId struct {
	Des []struct {
		ID     int   `json:"id"`
		Padids []int `json:"padids"`
	} `json:"des"`
}

func PadList(w http.ResponseWriter, r *http.Request) {
	// expected form of the request is :
	// {
	//         "des": [
	//                 {
	//                 "id": 102,
	//                 "padids": [
	//                         1,
	//                         3,
	//                         150,
	//                         14678
	//                 ]
	//                 }
	//         ]
	// }
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Content-type", "application/json")
	var p RequestPadId
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&p)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	var padlist []PadRef
	for _, de := range p.Des {
		for _, pad := range de.Padids {
			padlist = append(padlist, PadRef{de.ID, pad})
		}
	}
	jsonPadList(w, padlist)
}
