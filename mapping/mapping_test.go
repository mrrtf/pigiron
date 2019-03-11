package mapping_test

import (
	"fmt"
	"io"

	"github.com/mrrtf/pigiron/mapping"
)

// PrintPad prints all known information about a pad.
func PrintCathodePad(out io.Writer, cseg mapping.CathodeSegmentation, padcid mapping.PadCID) {
	if !cseg.IsValid(padcid) {
		fmt.Printf("invalid pad")
		return
	}
	fmt.Fprintf(out, "DE %4d DSID %4d CH %2d X %7.2f Y %7.2f DX %7.2f DY %7.2f\n",
		cseg.DetElemID(),
		cseg.PadDualSampaID(padcid),
		cseg.PadDualSampaChannel(padcid),
		cseg.PadPositionX(padcid),
		cseg.PadPositionY(padcid),
		cseg.PadSizeX(padcid),
		cseg.PadSizeY(padcid))
}

// PrintPad prints all known information about a pad.
func PrintPad(out io.Writer, seg mapping.Segmentation, paduid mapping.PadUID) {
	if !seg.IsValid(paduid) {
		fmt.Printf("invalid pad")
		return
	}
	fmt.Fprintf(out, "DE %4d DSID %4d CH %2d X %7.2f Y %7.2f DX %7.2f DY %7.2f\n",
		seg.DetElemID(),
		seg.PadDualSampaID(paduid),
		seg.PadDualSampaChannel(paduid),
		seg.PadPositionX(paduid),
		seg.PadPositionY(paduid),
		seg.PadSizeX(paduid),
		seg.PadSizeY(paduid))

}
