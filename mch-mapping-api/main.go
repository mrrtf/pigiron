package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/aphecetche/pigiron/mapping"
	"github.com/spf13/viper"

	// must include the specific implementation package of the mapping
	_ "github.com/aphecetche/pigiron/mapping/impl4"
)

func usage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `<h1>MCH Mapping API</h1>

<p>This API gives 2D geometric information about detection elements, dual sampas and pads.</p>

<h2>Detection element plane envelop</h2>

<pre>/degeo?deid=[number]&bending=[true|false]</pre>

<p>both deid and bending options are required</p>

<h2>Dual sampa envelops</h2>

<pre>/dualsampas?deid=[number]&bending=[true|false]</pre>

<p>Returns the vertices of the polygons describing the outline of all the dual sampas 
of a given detection element plane</p>
`)
	}
}

func dualSampas() http.HandlerFunc {
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
		jsonDualSampas(w, cseg, bending)
	}
}

func dualSampaPads() http.HandlerFunc {
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
		jsonDualSampaPads(w, seg, dsid)
	}
}

func deGeo() http.HandlerFunc {
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
		jsonDEGeo(w, cseg, bending)
	}
}

func main() {
	viper.SetEnvPrefix("MCH")
	viper.BindEnv("MAPPING_API_PORT")
	viper.SetDefault("MAPPING_API_PORT", 8080)
	port := viper.GetInt("MAPPING_API_PORT")
	fmt.Println("Started server to listen on port", port)

	http.HandleFunc("/", usage())
	http.HandleFunc("/dualsampas", dualSampas())
	http.HandleFunc("/dualsampapads", dualSampaPads())
	http.HandleFunc("/degeo", deGeo())
	if err := http.ListenAndServe(":"+strconv.Itoa(port), nil); err != nil {
		panic(err)
	}
}
