package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strconv"

	"github.com/mrrtf/pigiron/mapping"
	v2 "github.com/mrrtf/pigiron/mch-mapping-api/v2"
	"github.com/spf13/viper"

	// must include the specific implementation package of the mapping
	_ "github.com/mrrtf/pigiron/mapping/impl4"
)

var (
	ErrMissingDeId         = errors.New("Specifying a detection element id (deid=[number]) is required")
	ErrMissingDsId         = errors.New("Specifying a dualsampa id (dsid=[number]) is required")
	ErrMissingBending      = errors.New("Specifying a bending plane (bending=true or bending=false) is required")
	ErrDeIdShouldBeInteger = errors.New("deid should be an integer")
	ErrDsIdShouldBeInteger = errors.New("dsid should be an integer")
	ErrInvalidBending      = errors.New("bending should be true or false")
	ErrInvalidDeId         error
	ErrInvalidDsId         error
	validdeids             []int
)

func init() {
	mapping.ForEachDetectionElement(func(i mapping.DEID) {
		validdeids = append(validdeids, int(i))
	})
	sort.Ints(validdeids)
	s, _ := json.Marshal(validdeids)
	ErrInvalidDeId = errors.New("Invalid deid. Possible values are :" + string(s))
	ErrInvalidDsId = errors.New("Invalid dsid")
}

type Bending struct {
	present bool
	value   bool
}

// getDeIdBending decode the query part of the url, expecting it to be
// of the form : deid=[number]&bending=[true|false].
// the bending part is optional.
func getDeIdBending(u *url.URL) (int, Bending, error) {
	q := u.Query()
	de, ok := q["deid"]
	if !ok {
		return -1, Bending{}, ErrMissingDeId
	}
	deid, err := strconv.Atoi(de[0])
	if err != nil {
		return -1, Bending{}, ErrDeIdShouldBeInteger
	}
	l := sort.SearchInts(validdeids, deid)
	if l >= len(validdeids) || validdeids[l] != deid {
		return -1, Bending{}, ErrInvalidDeId
	}
	_, ok = q["bending"]
	if ok {
		b, err := strconv.ParseBool(q["bending"][0])
		if err != nil {
			return -1, Bending{}, ErrInvalidBending
		}
		return deid, Bending{present: true, value: b}, nil
	}
	return deid, Bending{present: false}, nil
}

// getDeDsBending decode the query part of the url, expecting it to be
// of the form : deid=[number]&dsid=[number]
// the bending part is optional.
func getDeIdDsId(u *url.URL) (int, int, error) {
	q := u.Query()
	de, ok := q["deid"]
	if !ok {
		return -1, -1, ErrMissingDeId
	}
	deid, err := strconv.Atoi(de[0])
	if err != nil {
		return -1, -1, ErrDeIdShouldBeInteger
	}
	l := sort.SearchInts(validdeids, deid)
	if l >= len(validdeids) || validdeids[l] != deid {
		return -1, -1, ErrInvalidDeId
	}
	ds, ok := q["dsid"]
	if !ok {
		return -1, -1, ErrMissingDsId
	}
	dsid, err := strconv.Atoi(ds[0])
	if err != nil {
		return -1, -1, ErrDsIdShouldBeInteger
	}
	bending := true
	if dsid >= 1024 {
		bending = false
	}
	cseg := mapping.NewCathodeSegmentation(mapping.DEID(deid), bending)
	validds := false
	for i := 0; i < cseg.NofDualSampas(); i++ {
		d, err := cseg.DualSampaID(i)
		if err != nil {
			fmt.Println(i, err)
		}
		if int(d) == dsid {
			validds = true
		}
	}
	if !validds {
		return -1, -1, ErrInvalidDsId
	}

	return deid, dsid, nil
}

func makeHandler(fn func(w http.ResponseWriter, r *http.Request, deid int, bending bool), isBendingRequired bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-type", "application/json")
		deid, bending, err := getDeIdBending(r.URL)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if isBendingRequired && !bending.present {
			http.Error(w, ErrMissingBending.Error(), http.StatusBadRequest)
			return
		}
		fn(w, r, deid, bending.value)
	}
}

func makeDEDSHandler(fn func(w http.ResponseWriter, r *http.Request, deid int, dsid int)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Content-type", "application/json")
		deid, dsid, err := getDeIdDsId(r.URL)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		fn(w, r, deid, dsid)
	}
}

func dualSampas(w http.ResponseWriter, r *http.Request, deid int, bending bool) {
	cseg := mapping.NewCathodeSegmentation(mapping.DEID(deid), bending)
	jsonDualSampas(w, cseg, bending)
}

func deGeo(w http.ResponseWriter, r *http.Request, deid int, bending bool) {
	cseg := mapping.NewCathodeSegmentation(mapping.DEID(deid), bending)
	jsonDEGeo(w, cseg, bending)
}

func handler() http.Handler {
	r := http.NewServeMux()
	r.HandleFunc("/", usage())
	bendingIsRequired := true
	r.HandleFunc("/dualsampas", makeHandler(dualSampas, bendingIsRequired))
	r.HandleFunc("/degeo", makeHandler(deGeo, bendingIsRequired))
	r.HandleFunc("/v2/dualsampas", makeHandler(v2.DualSampas, bendingIsRequired))
	r.HandleFunc("/v2/dualsampapads", makeDEDSHandler(v2.GetDualSampaPads))
	r.HandleFunc("/v2/degeo", makeHandler(v2.DeGeo, bendingIsRequired))
	r.HandleFunc("/v2/padlist", v2.PadList)
	return r
}

func main() {
	viper.SetEnvPrefix("MCH")
	viper.BindEnv("MAPPING_API_PORT")
	viper.SetDefault("MAPPING_API_PORT", 8080)
	port := viper.GetInt("MAPPING_API_PORT")
	fmt.Println("Started server to listen on port", port)
	if err := http.ListenAndServe(":"+strconv.Itoa(port), handler()); err != nil {
		panic(err)
	}
}
