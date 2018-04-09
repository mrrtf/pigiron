package geo

import (
	"testing"
)

func TestRemoveDuplicates(t *testing.T) {
	input := []float64{12, 14.1, 1.4, 12, 14.1, 4, 42}
	expected := []float64{12, 14.1, 1.4, 4, 42}
	output := removeDuplicates(input)
	if !EqualFloat64Slice(output, expected) {
		t.Errorf("expected %v and got %v", expected, output)
	}

}
