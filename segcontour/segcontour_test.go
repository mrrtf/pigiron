package segcontour

import (
	"os"
	"testing"

	"github.com/aphecetche/pigiron/geo"
	"github.com/aphecetche/pigiron/mapping"
)

func TestSegmentationBBox(t *testing.T) {
	for _, test := range []struct {
		detElemID int
		isBending bool
		want      geo.BBox
	}{
		{100, true, geo.NewBBoxUnchecked(0, 0, 89.04, 89.46)},
		{300, true, geo.NewBBoxUnchecked(-1, -0.75, 116, 117.25)},
		{500, true, geo.NewBBoxUnchecked(-75, -20, 57.5, 20)},
		{501, true, geo.NewBBoxUnchecked(-75, -20, 80, 20)},
		{502, true, geo.NewBBoxUnchecked(-80, -20, 75, 20)},
		{503, true, geo.NewBBoxUnchecked(-60, -20, 60, 20)},
		{504, true, geo.NewBBoxUnchecked(-40, -20, 40, 20)},
		{600, true, geo.NewBBoxUnchecked(-80, -20, 57.5, 20)},
		{601, true, geo.NewBBoxUnchecked(-80, -20, 80, 20)},
		{602, true, geo.NewBBoxUnchecked(-80, -20, 80, 20)},
		{700, true, geo.NewBBoxUnchecked(-100, -20, 100, 20)},
		{701, true, geo.NewBBoxUnchecked(-120, -20, 120, 20)},
		{702, true, geo.NewBBoxUnchecked(-100, -20, 100, 20)},
		{703, true, geo.NewBBoxUnchecked(-100, -20, 100, 20)},
		{704, true, geo.NewBBoxUnchecked(-80, -20, 80, 20)},
		{705, true, geo.NewBBoxUnchecked(-60, -20, 60, 20)},
		{706, true, geo.NewBBoxUnchecked(-40, -20, 40, 20)},
		{902, true, geo.NewBBoxUnchecked(-120, -20, 120, 20)},
		{903, true, geo.NewBBoxUnchecked(-120, -20, 120, 20)},
		{904, true, geo.NewBBoxUnchecked(-100, -20, 100, 20)},
		{905, true, geo.NewBBoxUnchecked(-80, -20, 80, 20)},

		{100, false, geo.NewBBoxUnchecked(-0.315, 0.21, 89.145, 89.25)},
		{300, false, geo.NewBBoxUnchecked(-0.625, -0.5, 115.625, 117.5)},
		{500, false, geo.NewBBoxUnchecked(-74.2857, -20, 58.5714, 20)},
		{501, false, geo.NewBBoxUnchecked(-74.2857, -20, 80, 20)},
		{502, false, geo.NewBBoxUnchecked(-80, -20, 74.2857, 20)},
		{503, false, geo.NewBBoxUnchecked(-60, -20, 60, 20)},
		{504, false, geo.NewBBoxUnchecked(-40, -20, 40, 20)},
		{600, false, geo.NewBBoxUnchecked(-80, -20, 58.5714, 20)},
		{601, false, geo.NewBBoxUnchecked(-80, -20, 80, 20)},
		{602, false, geo.NewBBoxUnchecked(-80, -20, 80, 20)},
		{700, false, geo.NewBBoxUnchecked(-100, -20, 100, 20)},
		{701, false, geo.NewBBoxUnchecked(-120, -20, 120, 20)},
		{702, false, geo.NewBBoxUnchecked(-100, -20, 100, 20)},
		{703, false, geo.NewBBoxUnchecked(-100, -20, 100, 20)},
		{704, false, geo.NewBBoxUnchecked(-80, -20, 80, 20)},
		{705, false, geo.NewBBoxUnchecked(-60, -20, 60, 20)},
		{706, false, geo.NewBBoxUnchecked(-40, -20, 40, 20)},
		{902, false, geo.NewBBoxUnchecked(-120, -20, 120, 20)},
		{903, false, geo.NewBBoxUnchecked(-120, -20, 120, 20)},
		{904, false, geo.NewBBoxUnchecked(-100, -20, 100, 20)},
		{905, false, geo.NewBBoxUnchecked(-80, -20, 80, 20)},
	} {
		if test.detElemID != 706 && test.detElemID != 500 {
			continue
		}

		if test.isBending == true {
			continue
		}
		seg := mapping.NewCathodeSegmentation(test.detElemID, test.isBending)
		bbox := BBox(seg)
		if !geo.EqualBBox(bbox, test.want) {
			t.Errorf("segmentation %3d - %v : wrong bbox got\n%v but want\n%v", test.detElemID,
				mapping.PlaneAbbreviation(test.isBending), bbox.String(), test.want.String())
		}
	}
}

type padSize struct {
	x, y float64
}

func TestPadSizes(t *testing.T) {

	padsizes := make(map[padSize]int)

	mapping.ForOneDetectionElementOfEachSegmentationType(func(detElemID int) {
		for _, isBending := range []bool{true, false} {
			seg := mapping.NewCathodeSegmentation(detElemID, isBending)
			seg.ForEachPad(func(paduid mapping.PadUID) {
				ps := &padSize{seg.PadSizeX(paduid), seg.PadSizeY(paduid)}
				padsizes[*ps]++
			})
		}
	})

	if len(padsizes) != 18 {
		t.Errorf("wanted 18 padsizes - got %d", len(padsizes))
	}

	b, _ := geo.NewBBox(0, 0, 12, 12)
	svg := geo.SVGWriter{Width: 1024, BBox: b}

	svg.Style(`
rect {
stroke: red;
stroke-width: 0.02;
fill: none;
}`)

	for ps := range padsizes {
		svg.Rect(1.0, 1.0, ps.x, ps.y)
	}

	f, _ := os.Create("padsizes.html")
	svg.WriteHTML(f)

}
