package main

import (
	"fmt"

	"github.com/aphecetche/pigiron/mapping"
)

func main() {
	seg := mapping.NewSegmentation(103, true)
	fmt.Printf("seg=%v", seg)
}
