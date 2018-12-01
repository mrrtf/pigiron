package mapping

import (
	"fmt"
)

type cathodeSegmentationBuilder interface {
	Build(isBendingPlane bool) CathodeSegmentation
}

var builderRegistry map[int]cathodeSegmentationBuilder

func registerCathodeSegmentationBuilder(segType int, builder cathodeSegmentationBuilder) {
	if builderRegistry == nil {
		builderRegistry = make(map[int]cathodeSegmentationBuilder)
	}
	_, alreadyThere := builderRegistry[segType]
	if alreadyThere {
		fmt.Printf("already got a build for segType %d : will not override it", segType)
		return
	}
	builderRegistry[segType] = builder
}

func getCathodeSegmentationBuilder(segType int) cathodeSegmentationBuilder {
	return builderRegistry[segType]
}
