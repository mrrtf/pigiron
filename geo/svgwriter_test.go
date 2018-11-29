package geo

import (
	"testing"
)

var (
	box BBox
)

func init() {
	var err error
	box, err = NewBBox(1, 2, 3, 4)
	if err != nil {
		panic(err)
	}
}

func TestSVGWriterWithOriginOptionSetHasBBoxShiftedToOrigin(t *testing.T) {
	svg := NewSVGWriter(1000, box, true)
	x := svg.BBox.Xmin()
	y := svg.BBox.Ymin()
	if !(x == 0.0 && y == 0.0) {
		t.Errorf("Want (x,y)=(%e,%e) - Got %e,%e", box.Xmin(), box.Ymin(), x, y)
	}
}

func TestCreateSVGWriterWithOriginNotSetAsOriginAtBBox(t *testing.T) {
	svg := NewSVGWriter(1000, box, false)
	x := svg.BBox.Xmin()
	y := svg.BBox.Ymin()
	if !(x == box.Xmin() && y == box.Ymin()) {
		t.Errorf("Want (x,y)=(%e,%e) - Got %e,%e", box.Xmin(), box.Ymin(), x, y)
	}
}
