package v2

import (
	"net/http"

	"github.com/mrrtf/pigiron/mapping"
)

func DualSampas(w http.ResponseWriter, r *http.Request, deid int, bending bool) {
	cseg := mapping.NewCathodeSegmentation(mapping.DEID(deid), bending)
	jsonDualSampas(w, cseg, bending)
}
