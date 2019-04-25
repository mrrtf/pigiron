package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"

	"github.com/mrrtf/pigiron/mapping"
)

type Pad struct {
	DSID int     `json:"DSID"`
	DSCH int     `json:"DSCH"`
	X    float64 `json:"X"`
	Y    float64 `json:"Y"`
	SX   float64 `json:"SX"`
	SY   float64 `json:"SY"`
}

func jsonPadsInArea(w io.Writer, cseg mapping.CathodeSegmentation,
	xmin, ymin, xmax, ymax float64) {

	cseg.ForEachPadInArea(xmin, ymin, xmax, ymax, func(padcid mapping.PadCID) {
		b, err := json.Marshal(
			Pad{DSID: int(cseg.PadDualSampaID(padcid)),
				DSCH: int(cseg.PadDualSampaChannel(padcid)),
				X:    cseg.PadPositionX(padcid),
				Y:    cseg.PadPositionY(padcid),
				SX:   cseg.PadSizeX(padcid),
				SY:   cseg.PadSizeY(padcid),
			})

		if err != nil {
			log.Fatalf(err.Error())
		}

		fmt.Fprintf(w, string(b))
	})

}
