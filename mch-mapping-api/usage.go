package main

import (
	"fmt"
	"net/http"
)

func usage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, usageMsg)
	}
}

const usageMsg = `<h1>MCH Mapping API</h1>

<p>This API gives 2D geometric information about detection elements, dual sampas and pads.</p>

<h2>Detection element plane envelop</h2>

<pre>/degeo?deid=[number]&bending=[true|false]</pre>

<p>both deid and bending options are required</p>

<h2>Dual sampa envelops</h2>

<pre>/dualsampas?deid=[number]&bending=[true|false]</pre>

<p>Returns the vertices of the polygons describing the outline of all the dual sampas 
of a given detection element plane</p>

<h2>Pads in area</h2>

<pre>/padsinarea?deid=[integer](&bending=[true|false])&xmin=[float]&ymin=[float]&xmax=[float]&ymax=[float]</pre>

<p>Returns the pads which surface intersects the area specified by xmin,ymin,xmax,ymax.</p>
`
