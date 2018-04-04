package mapping

import "fmt"

type segType0 struct{}

func (seg segType0) Build(isBendingPlane bool) Segmentation {
	if isBendingPlane {
		return newSegmentation(0, true, nil, nil, nil)
	}
	return newSegmentation(0, false, nil, nil, nil)
}

func init() {
	fmt.Println("could create seg0 maker here")
	registerSegmentationBuilder(0, segType0{})
}
