package output

import (
	"testing"

	"github.com/aphecetche/pigiron/mapping"
	// must include the specific implementation package of the mapping
	_ "github.com/aphecetche/pigiron/mapping/impl4"
)

func TestJSONDualSampas(t *testing.T) {

	expected := []int{3, 4, 7, 8, 101, 102, 103, 106, 107, 108}

	var ds []int
	cseg := mapping.NewCathodeSegmentation(706, true)
	for i := 0; i < cseg.NofDualSampas(); i++ {
		dsid, err := cseg.DualSampaID(i)
		if err != nil {
			panic(err)
		}
		ds = append(ds, int(dsid))
	}

	if len(expected) != len(ds) {
		t.Errorf("Want %d dual sampas. Got %d\n", len(expected), len(ds))
	}

	for i, e := range expected {
		if e != ds[i] {
			t.Errorf("Want ds[%d]=%d. Got %d\n", i, e, ds[i])
		}
	}
}
