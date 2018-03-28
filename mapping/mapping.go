package mapping

/*
#cgo CFLAGS: -I/Users/laurent/alice/sw/osx_x86-64/O2/clion-1/include/MCHMappingInterface
 #cgo LDFLAGS: -L/Users/laurent/alice/sw/osx_x86-64/O2/clion-1/lib
 #cgo LDFLAGS: -lMCHMappingImpl3
 #include "SegmentationCInterface.h"
*/
import "C"

import (
	"fmt"
)

// Segmentation is the main entry point to the MCH mapping
type Segmentation struct {
	handle C.MchSegmentationHandle
}

// NewSegmentation creates a segmentation object for the given detection element plane
func NewSegmentation(deid int, bending int) (Segmentation, error) {
	d := C.int(deid)
	b := C.int(bending)
	var seg Segmentation
	seg.handle = C.mchSegmentationConstruct(d, b)
	if seg.handle == nil {
		return seg, fmt.Errorf("zob")
	}
	return seg, nil
}

func (seg *Segmentation) FindPadByFEE(dualSampaID int, dualSampaChannel int) (int, error) {
	padid := C.mchSegmentationFindPadByFEE(seg.handle, C.int(dualSampaID), C.int(dualSampaChannel))
	if int(C.mchSegmentationIsPadValid(seg.handle, C.int(padid))) > 0 {
		return int(padid), nil
	}
	return 0, fmt.Errorf("Could not get a pad for dualSampaId %d channel %d", dualSampaID, dualSampaChannel)
}
