package mapping

import (
	"fmt"
)

type SegmentationBuilder interface {
	Build(isBendingPlane bool) Segmentation
}

var builderRegistry map[int]SegmentationBuilder

func RegisterSegmentationBuilder(segType int, builder SegmentationBuilder) error {

	if builderRegistry == nil {
		builderRegistry = make(map[int]SegmentationBuilder)
	}
	_, alreadyThere := builderRegistry[segType]
	if alreadyThere {
		return fmt.Errorf("already got a build for segType %d : will not override it", segType)
	}
	builderRegistry[segType] = builder
	return nil
}

func GetSegmentationBuilder(segType int) SegmentationBuilder {
	return builderRegistry[segType]
}
