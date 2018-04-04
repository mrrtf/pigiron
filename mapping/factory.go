package mapping

import (
	"fmt"
)

type segmentationBuilder interface {
	Build(isBendingPlane bool) Segmentation
}

var builderRegistry map[int]segmentationBuilder

func registerSegmentationBuilder(segType int, builder segmentationBuilder) error {

	if builderRegistry == nil {
		builderRegistry = make(map[int]segmentationBuilder)
	}
	_, alreadyThere := builderRegistry[segType]
	if alreadyThere {
		return fmt.Errorf("already got a build for segType %d : will not override it", segType)
	}
	builderRegistry[segType] = builder
	return nil
}

func getSegmentationBuilder(segType int) segmentationBuilder {
	return builderRegistry[segType]
}
