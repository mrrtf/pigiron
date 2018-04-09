package geo

import (
	"fmt"
	"testing"
)

func newTestNode(dummyCardinality int) *node {
	left := newNode(0, 4, 2)
	right := newNode(4, 8, 6)

	node := newNode(0, 8, 4)
	node.setLeft(left).setRight(right)
	left.setCardinality(dummyCardinality)
	right.setCardinality(dummyCardinality)
	return node
}

var (
	dummyCardinality = 3
	dummyNode        = newNode(0, 8, 4)
)

func TestNeedAtLeastTwoValuesToBuildASegmentTree(t *testing.T) {
	var values = []float64{42}
	_, err := createSegmentTree(values)
	if err == nil {
		t.Error("should not be able to create a segment tree from only one value")
	}
}

func TestJustCreatedNodeIsNotPotent(t *testing.T) {
	if dummyNode.potent != false {
		t.Error("just created node should not be potent")
	}
}
func TestJustCreatedNodeHasCardinalityEqualsToZero(t *testing.T) {
	if dummyNode.cardinality != 0 {
		t.Error("just created node should have cardinality equals to zero")
	}
}

func TestPromoteNode(t *testing.T) {
	testNode := newTestNode(dummyCardinality)
	testNode.promote()
	if testNode.cardinality != 1 {
		t.Error("promoted testNode should have cardinality of 1")
	}
	if testNode.leftChild.cardinality != dummyCardinality-1 {
		t.Errorf("promoted left node should have cardinality of %d but got %d",
			dummyCardinality-1, testNode.leftChild.cardinality)
	}
	if testNode.rightChild.cardinality != dummyCardinality-1 {
		t.Errorf("promoted right node should have cardinality of %d but got %d",
			dummyCardinality-1, testNode.rightChild.cardinality)
	}
}

func TestDemoteNode(t *testing.T) {
	testNode := newTestNode(dummyCardinality)
	testNode.promote()
	testNode.demote()
	if testNode.cardinality != 0 {
		t.Error("demoted testNode should have cardinality of 0")
	}
	if testNode.leftChild.cardinality != dummyCardinality {
		t.Errorf("demoted left node should have cardinality of %d but got %d",
			dummyCardinality, testNode.leftChild.cardinality)
	}
	if testNode.rightChild.cardinality != dummyCardinality {
		t.Errorf("demoted right node should have cardinality of %d but got %d",
			dummyCardinality, testNode.rightChild.cardinality)
	}
	if !testNode.potent {
		t.Error("demoted node should be potent")
	}
}

func TestMidPointOfANodeIsNotHalfPoint(t *testing.T) {
	ypos := []float64{-2.0, -1.5, -1, 0}
	node, err := createSegmentTree(ypos)
	if err != nil {
		t.Fatal("could not create segment tree")
	}
	right := node.rightChild
	if (right.interval != interval{-1.5, 0.0}) {
		t.Errorf("right node interval expected to be -1.5,0 but is %f,%f",
			right.interval.b, right.interval.e)
	}
	mp := right.midpoint
	if mp == 1.5/2 {
		t.Error("wrong")
	}
	if mp != -1 {
		t.Error("wrong")
	}
}

func TestNodeInsertAndDelete(t *testing.T) {
	ypos := []float64{0.0, 0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8}

	tree, err := createSegmentTree(ypos)
	if err != nil {
		t.Fatal("could not create segmentree")
	}

	tree.insertInterval(interval{0.1, 0.5})
	tree.insertInterval(interval{0.5, 0.8})
	tree.deleteInterval(interval{0.6, 0.7})
	expected := `[0,0.8] potent
      [0,0.4] potent
            [0,0.2] potent
                  [0,0.1]
                  [0.1,0.2] C=1
            [0.2,0.4] C=1
                  [0.2,0.3]
                  [0.3,0.4]
      [0.4,0.8] potent
            [0.4,0.6] C=1
                  [0.4,0.5]
                  [0.5,0.6]
            [0.6,0.8] potent
                  [0.6,0.7]
                  [0.7,0.8] C=1
`

	actual := tree.String()
	if actual != expected {
		t.Error("output string not as expected")
		fmt.Printf("expected:%s.", expected)
		fmt.Printf("actual:%s.", actual)
		fmt.Println(len(actual), len(expected))
	}
}
