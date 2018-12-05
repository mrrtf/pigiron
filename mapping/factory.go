package mapping

import (
	"fmt"
)

type cathodeSegmentationBuilder interface {
	Build(isBendingPlane bool, deid int) CathodeSegmentation
}

var builderRegistry map[int]cathodeSegmentationBuilder

func RegisterCathodeSegmentationBuilder(segType int, builder cathodeSegmentationBuilder) {
	if builderRegistry == nil {
		builderRegistry = make(map[int]cathodeSegmentationBuilder)
	}
	_, alreadyThere := builderRegistry[segType]
	if alreadyThere {
		fmt.Printf("already got a build for segType %d : will not override it", segType)
		return
	}
	fmt.Printf("register builder for segType %d ", segType)
	builderRegistry[segType] = builder
}

func getCathodeSegmentationBuilder(segType int) cathodeSegmentationBuilder {
	return builderRegistry[segType]
}
