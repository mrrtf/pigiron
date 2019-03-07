package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const sampleDualSampaAnswer = `{"DualSampas":[{"ID":3,"Vertices":[{"X":0,"Y":-4},{"X":0,"Y":-20},{"X":20,"Y":-20},{"X":20,"Y":-4},{"X":0,"Y":-4}],"Value":0},{"ID":4,"Vertices":[{"X":20,"Y":-4},{"X":20,"Y":-20},{"X":40,"Y":-20},{"X":40,"Y":-4},{"X":20,"Y":-4}],"Value":0},{"ID":7,"Vertices":[{"X":-40,"Y":-4},{"X":-40,"Y":-20},{"X":-20,"Y":-20},{"X":-20,"Y":-4},{"X":-40,"Y":-4}],"Value":0},{"ID":8,"Vertices":[{"X":-20,"Y":-4},{"X":-20,"Y":-20},{"X":0,"Y":-20},{"X":0,"Y":-4},{"X":-20,"Y":-4}],"Value":0},{"ID":101,"Vertices":[{"X":20,"Y":4},{"X":20,"Y":-4},{"X":40,"Y":-4},{"X":40,"Y":20},{"X":30,"Y":20},{"X":30,"Y":4},{"X":20,"Y":4}],"Value":0},{"ID":102,"Vertices":[{"X":10,"Y":20},{"X":10,"Y":4},{"X":30,"Y":4},{"X":30,"Y":20},{"X":10,"Y":20}],"Value":0},{"ID":103,"Vertices":[{"X":0,"Y":20},{"X":0,"Y":-4},{"X":20,"Y":-4},{"X":20,"Y":4},{"X":10,"Y":4},{"X":10,"Y":20},{"X":0,"Y":20}],"Value":0},{"ID":106,"Vertices":[{"X":-20,"Y":4},{"X":-20,"Y":-4},{"X":0,"Y":-4},{"X":0,"Y":20},{"X":-10,"Y":20},{"X":-10,"Y":4},{"X":-20,"Y":4}],"Value":0},{"ID":107,"Vertices":[{"X":-30,"Y":20},{"X":-30,"Y":4},{"X":-10,"Y":4},{"X":-10,"Y":20},{"X":-30,"Y":20}],"Value":0},{"ID":108,"Vertices":[{"X":-40,"Y":20},{"X":-40,"Y":-4},{"X":-20,"Y":-4},{"X":-20,"Y":4},{"X":-30,"Y":4},{"X":-30,"Y":20},{"X":-40,"Y":20}],"Value":0}]}`

const sampleDEGeoAnswer = `{"ID":706,"Bending":true,"X":0,"Y":0,"SX":80,"SY":40}`

func TestDualSampaEndPoint(t *testing.T) {
	endPointRequiringDeIdBending(t, makeHandler(dualSampas, true), "dualsampas", sampleDualSampaAnswer)
}

func TestDEGeoEndPoint(t *testing.T) {
	endPointRequiringDeIdBending(t, makeHandler(deGeo, true), "degeo", sampleDEGeoAnswer)
}

func endPointRequiringDeIdBending(t *testing.T, h http.HandlerFunc, endpoint string, sampleAnswer string) {

	tt := []struct {
		name   string
		query  string
		answer string
		status int
		errmsg string
	}{
		{"happy path", "deid=706&bending=true", sampleAnswer, http.StatusOK, ""},
		{"missing bending", "deid=501", "", http.StatusBadRequest, ErrMissingBending.Error()},
		{"deid not an integer", "deid=xx", "", http.StatusBadRequest, ErrDeIdShouldBeInteger.Error()},
		{"bending not a boolean", "deid=706&bending=x", "", http.StatusBadRequest, ErrInvalidBending.Error()},
		{"missing deid", "", "", http.StatusBadRequest, ErrMissingDeId.Error()},
		{"missing deid (bis)", "bending=false", "", http.StatusBadRequest, ErrMissingDeId.Error()},
		{"invalid deid", "deid=104&bending=false", "", http.StatusBadRequest, ErrInvalidDeId.Error()},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/"+endpoint+"?"+tc.query, nil)
			if err != nil {
				t.Fatalf("Could not create request: %v", err)
			}
			rec := httptest.NewRecorder()
			h(rec, req)
			resp := rec.Result()
			defer resp.Body.Close()
			if resp.StatusCode != tc.status {
				t.Fatalf("Expected status code %d. Got %d", tc.status, resp.StatusCode)
			}
			body, _ := ioutil.ReadAll(resp.Body)
			s := strings.Trim(string(body), " \n")
			if tc.errmsg != "" && s != tc.errmsg {
				t.Fatalf("Expected error message %q. Got %q", tc.errmsg, s)
			}
			if tc.answer != "" && s != tc.answer {
				t.Fatalf("Expected answer %q. Got %q", tc.answer, s)
			}
		})
	}
}
