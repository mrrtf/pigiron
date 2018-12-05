package main

import (
	"github.com/aphecetche/pigiron/mapping"
	// must include the specific implementation package of the mapping
	_ "github.com/aphecetche/pigiron/mapping/impl4"
)

func main() {
	_ = mapping.NewCathodeSegmentation(103, true)
}
