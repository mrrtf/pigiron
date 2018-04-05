package mapping

import (
	"fmt"
)

type segmentationBuilder interface {
	Build(isBendingPlane bool) Segmentation
}

var builderRegistry map[int]segmentationBuilder

func registerSegmentationBuilder(segType int, builder segmentationBuilder) {
	if builderRegistry == nil {
		builderRegistry = make(map[int]segmentationBuilder)
	}
	_, alreadyThere := builderRegistry[segType]
	if alreadyThere {
		fmt.Printf("already got a build for segType %d : will not override it", segType)
		return
	}
	builderRegistry[segType] = builder
}

func getSegmentationBuilder(segType int) segmentationBuilder {
	return builderRegistry[segType]
}
