package mapping

type segPair struct {
	Bending, NonBending Segmentation
}

// SegCache is the simple map of segmentation pairs
// (one segmentation for bending plane and one for non-bending plane)
type SegCache struct {
	segpairs map[int]segPair
}

// Segmentation returns the segmentation for given detection element id
// and given plane (true for bending plane).
// The segmentation for both planes of that detection element is created
// and cached if not already cached
func (sc *SegCache) Segmentation(deid int, bending bool) Segmentation {
	if sc.segpairs == nil {
		sc.segpairs = make(map[int]segPair)
	}
	seg := sc.segpairs[deid]
	if seg.Bending == nil {
		sc.segpairs[deid] = segPair{
			Bending:    NewSegmentation(deid, true),
			NonBending: NewSegmentation(deid, false),
		}
		seg = sc.segpairs[deid]
	}
	if bending {
		return seg.Bending
	}
	return seg.NonBending
}