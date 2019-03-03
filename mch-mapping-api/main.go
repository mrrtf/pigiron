package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/aphecetche/pigiron/mapping"
	"github.com/aphecetche/pigiron/mch-mapping-api/output"
	"github.com/aphecetche/pigiron/segcontour"
	"github.com/spf13/viper"

	// must include the specific implementation package of the mapping
	_ "github.com/aphecetche/pigiron/mapping/impl4"
)

func showSegmentation() http.HandlerFunc {
	// TODO : initialize a global cseg cache here (beware mutex)
	// to be used in the lambda below
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		deid, err := strconv.Atoi(q.Get("deid"))
		if err != nil {
			w.Write([]byte("Malformed deid"))
			return
		}
		bending, err := strconv.ParseBool(q.Get("bending"))
		if err != nil {
			w.Write([]byte("Malformed bending"))
			return
		}
		showDEs, err := strconv.ParseBool(q.Get("des"))
		if err != nil {
			showDEs = false
		}
		showPads, err := strconv.ParseBool(q.Get("pads"))
		if err != nil {
			showPads = false
		}
		showDualSampas, err := strconv.ParseBool(q.Get("dualsampas"))
		if err != nil {
			showDualSampas = false
		}
		cseg := mapping.NewCathodeSegmentation(mapping.DEID(deid), bending)

		formats, ok := q["format"]
		format := "svg"
		if ok {
			format = formats[0]
		}

		showflags := segcontour.ShowFlags{DEs: showDEs, DualSampas: showDualSampas, Pads: showPads, PadChannels: true}

		switch strings.ToLower(format) {
		case "svg":
			output.ToSVG(w, cseg, bending, showflags)
		case "d3":
			output.ToD3(w, cseg, bending, showflags)
		default:
			fmt.Fprintf(w, "Don't know how to deal with format %s\n", format)
		}
	}
}

func DualSampas() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		q := r.URL.Query()
		deid, err := strconv.Atoi(q.Get("deid"))
		if err != nil {
			w.Write([]byte("Malformed deid"))
			return
		}
		bending, err := strconv.ParseBool(q.Get("bending"))
		if err != nil {
			w.Write([]byte("Malformed bending"))
			return
		}
		cseg := mapping.NewCathodeSegmentation(mapping.DEID(deid), bending)
		output.JSONDualSampas(w, cseg, bending)
	}
}

func DualSampaPads() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		deid, err := strconv.Atoi(q.Get("deid"))
		if err != nil {
			w.Write([]byte("Malformed deid"))
			return
		}
		dsid, err := strconv.Atoi(q.Get("dsid"))
		if err != nil {
			w.Write([]byte("Malformed dsid"))
			return
		}
		seg := mapping.NewSegmentation(mapping.DEID(deid))
		output.JSONDualSampaPads(w, seg, dsid)
	}
}

func DEGeo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		q := r.URL.Query()
		deid, err := strconv.Atoi(q.Get("deid"))
		if err != nil {
			w.Write([]byte("Malformed deid"))
			return
		}
		bending, err := strconv.ParseBool(q.Get("bending"))
		if err != nil {
			w.Write([]byte("Malformed bending"))
			return
		}
		cseg := mapping.NewCathodeSegmentation(mapping.DEID(deid), bending)
		output.JSONDEGeo(w, cseg, bending)
	}
}

func main() {

	viper.SetEnvPrefix("MCH")
	viper.BindEnv("MAPPING_API_PORT")
	viper.SetDefault("MAPPING_API_PORT", 8080)
	fs := http.FileServer(http.Dir("assets"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", showSegmentation())
	http.HandleFunc("/dualsampas", DualSampas())
	http.HandleFunc("/dualsampapads", DualSampaPads())
	http.HandleFunc("/degeo", DEGeo())
	port := viper.GetInt("MAPPING_API_PORT")
	fmt.Println("Started server to listen on port", port)
	if err := http.ListenAndServe(":"+strconv.Itoa(port), nil); err != nil {
		panic(err)
	}
}
