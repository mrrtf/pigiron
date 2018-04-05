package mapping

import "fmt"

// A padGroupType is a collection of pads of unspecified size(s)
// organized in a certain way in a x-y rectilinear plane
type padGroupType struct {
	fastID      []int
	fastIndices []int
	nofPadsX    int
	nofPadsY    int
	nofPads     int
}

func (pgt padGroupType) NofPads() int {
	return pgt.nofPads
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
	pgt.fastID = ids
	pgt.fastIndices = validIndices(pgt.fastID)
	pgt.nofPads = len(pgt.fastIndices)
	pgt.nofPadsX = nofPadsX
	pgt.nofPadsY = nofPadsY
	return *pgt
}

func (pgt *padGroupType) String() string {
	s := fmt.Sprintf("n=%d nx=%d ny=%d\n", pgt.nofPads, pgt.nofPadsX, pgt.nofPadsY)
	s += "index "
	for i := 0; i < len(pgt.fastID); i++ {
		s += fmt.Sprintf("%2d ", pgt.fastID[i])
	}
	return s
}

func (pgt padGroupType) fastIndex(ix int, iy int) int {
	return ix + iy*pgt.nofPadsX
}

func (pgt padGroupType) idByFastIndex(fastIndex int) int {
	return fastIndex
}

// Return the index of the pad with indices = (ix,iy)
// or -1 if not found
func (pgt padGroupType) idByIndices(ix int, iy int) int {
	return pgt.idByFastIndex(pgt.fastIndex(ix, iy))
}
