package mapping

import "fmt"

// A padGroupType is a collection of pads of unspecified size(s)
// organized in a certain way in a x-y rectilinear plane
type padGroupType struct {
	FastID      []int
	FastIndices []int
	NofPadsX    int
	NofPadsY    int
	NofPads     int
}

func validIndices(v []int) []int {
	valid := []int{}
	for i := 0; i < len(v); i++ {
		if v[i] >= 0 {
			valid = append(valid, v[i])
		}
	}
	return valid
}

// NewPadGroupType returns a pad group type
func NewPadGroupType(nofPadsX int, nofPadsY int, ids []int) padGroupType {
	pgt := new(padGroupType)
	pgt.FastID = ids
	pgt.FastIndices = validIndices(pgt.FastID)
	pgt.NofPads = len(pgt.FastIndices)
	pgt.NofPadsX = nofPadsX
	pgt.NofPadsY = nofPadsY
	return *pgt
}

func (pgt *padGroupType) String() string {
	s := fmt.Sprintf("n=%d nx=%d ny=%d\n", pgt.NofPads, pgt.NofPadsX, pgt.NofPadsY)
	s += "index "
	for i := 0; i < len(pgt.FastID); i++ {
		s += fmt.Sprintf("%2d ", pgt.FastID[i])
	}
	return s
}

func (pgt *padGroupType) fastIndex(ix int, iy int) int {
	return ix + iy*pgt.NofPadsX
}

func (pgt *padGroupType) idByFastIndex(fastIndex int) int {
	if fastIndex >= 0 && fastIndex < len(pgt.FastID) {
		return pgt.FastID[fastIndex]
	}
	return -1
}

// Return the index of the pad with indices = (ix,iy)
// or -1 if not found
func (pgt *padGroupType) idByIndices(ix int, iy int) int {
	return pgt.idByFastIndex(pgt.fastIndex(ix, iy))
}

func (pgt *padGroupType) iy(fastIndex int) int {
	return fastIndex / pgt.NofPadsX
}

func (pgt *padGroupType) ix(fastIndex int) int {
	return fastIndex - pgt.iy(fastIndex)*pgt.NofPadsX
}

func (pgt *padGroupType) areIndicesPossible(ix, iy int) bool {
	return ix >= 0 && ix < pgt.NofPadsX && iy >= 0 && iy < pgt.NofPadsY
}
