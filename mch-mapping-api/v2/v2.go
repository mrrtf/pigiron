package v2

import (
	"net/http"

	"github.com/mrrtf/pigiron/mapping"
)

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
