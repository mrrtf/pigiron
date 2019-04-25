package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strconv"

	"github.com/hashicorp/go-multierror"
	"github.com/mrrtf/pigiron/mapping"
	"github.com/spf13/viper"

	// must include the specific implementation package of the mapping
	_ "github.com/mrrtf/pigiron/mapping/impl4"
)

var (
	ErrMissingDeId         = errors.New("Specifying a detection element id (deid=[number]) is required")
	ErrMissingBending      = errors.New("Specifying a bending plane (bending=true or bending=false) is required")
	ErrDeIdShouldBeInteger = errors.New("deid should be an integer")
	ErrInvalidBending      = errors.New("bending should be true or false")
	ErrInvalidDeId         error
	validdeids             []int
)

func init() {
	mapping.ForEachDetectionElement(func(i mapping.DEID) {
		validdeids = append(validdeids, int(i))
	})
	sort.Ints(validdeids)
	s, _ := json.Marshal(validdeids)
	ErrInvalidDeId = errors.New("Invalid deid. Possible values are :" + string(s))
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

func makeHandler(fn func(w http.ResponseWriter, r *http.Request, deid int, bending bool)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		deid, bending, err := getDeIdBending(r.URL)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if !bending.present {
			http.Error(w, ErrMissingBending.Error(), http.StatusBadRequest)
			return
		}
		fn(w, r, deid, bending.value)
	}
}

func getAreaLimit(e *error, limit string, u *url.URL) float64 {
	q := u.Query()
	s, ok := q[limit]
	if !ok {
		*e = multierror.Append(*e, fmt.Errorf("missing %s", limit))
		return 0.0
	}
	l, err := strconv.ParseFloat(s[0], 64)
	if err != nil {
		*e = multierror.Append(*e, fmt.Errorf("%s is not a float", limit))
		return 0.0
	}
	return l
}

func makePadHandler(fn func(w http.ResponseWriter, r *http.Request, deid int, bending bool, xmin, ymin, xmax, ymax float64)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		deid, bending, err := getDeIdBending(r.URL)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if !bending.present {
			http.Error(w, ErrMissingBending.Error(), http.StatusBadRequest)
			return
		}
		u := r.URL
		var e error

		xmin := getAreaLimit(&e, "xmin", u)
		ymin := getAreaLimit(&e, "ymin", u)
		xmax := getAreaLimit(&e, "xmax", u)
		ymax := getAreaLimit(&e, "ymax", u)
		if e != nil {
			http.Error(w, e.Error(), http.StatusBadRequest)
			return
		}
		fn(w, r, deid, bending.value, xmin, ymin, xmax, ymax)
	}
}

func dualSampas(w http.ResponseWriter, r *http.Request, deid int, bending bool) {
	cseg := mapping.NewCathodeSegmentation(mapping.DEID(deid), bending)
	jsonDualSampas(w, cseg)
}

func deGeo(w http.ResponseWriter, r *http.Request, deid int, bending bool) {
	cseg := mapping.NewCathodeSegmentation(mapping.DEID(deid), bending)
	jsonDEGeo(w, cseg)
}

func padsInArea(w http.ResponseWriter, r *http.Request, deid int, bending bool,
	xmin, ymin, xmax, ymax float64) {
	cseg := mapping.NewCathodeSegmentation(mapping.DEID(deid), bending)
	jsonPadsInArea(w, cseg, xmin, ymin, xmax, ymax)
}

func handler() http.Handler {
	r := http.NewServeMux()
	r.HandleFunc("/", usage())
	r.HandleFunc("/dualsampas", makeHandler(dualSampas))
	r.HandleFunc("/degeo", makeHandler(deGeo))
	r.HandleFunc("/padsinarea", makePadHandler(padsInArea))
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
