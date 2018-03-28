package mapping_test

import (
	"fmt"
	"testing"

	"github.com/aphecetche/pigiron/mapping"
)

func TestCreateSegmentation(t *testing.T) {

	seg, err := mapping.NewSegmentation(100, 1)
	if err != nil {
		t.Fatalf("Could not create segmentation")
	}
	fmt.Println(seg.FindPadByFEE(76, 9))

}
