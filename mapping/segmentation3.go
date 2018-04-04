package mapping

import "fmt"

//   Segmentation(int segType, bool isBendingPlane, std::vector<PadGroup> padGroups,
//                std::vector<PadGroupType> padGroupTypes, std::vector<std::pair<float, float>> padSizes);

type segmentation3 struct {
}

// NewSegmentation creates a segmentation object for the given detection element plane
func NewSegmentation(detElemID int, isBendingPlane bool) Segmentation {

	segType, err := detElemID2SegType(detElemID)
	if err != nil {
		return nil
	}
	fmt.Printf("detElemID %d -> segType %d\n", detElemID, segType)
	builder := GetSegmentationBuilder(segType)
	if builder == nil {
		return nil
	}
	return builder.Build(isBendingPlane)
}

func newSegmentation(segType int, isBendingPlane bool, padGroups []PadGroup,
	padGroupTypes []PadGroupType, padSizes []PadSize) *segmentation3 {

	fmt.Println("segType", segType, "isBendingPlane", isBendingPlane)
	return &segmentation3{}
}

func (seg segmentation3) IsValid(padid int) bool {
	return false
}
func (seg segmentation3) FindPadByFEE(dualSampaID int, dualSampaChannel int) (int, error) {
	return 0, fmt.Errorf("invalid pad")
}
func (seg segmentation3) FindPadByPosition(x float64, y float64) (int, error) {
	return 0, fmt.Errorf("invalid pad")
}
