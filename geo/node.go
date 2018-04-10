package geo

import (
	"errors"
	"fmt"
	"log"
	"sort"
	"strings"
)

// in order to reduce the amount of code (and because the usage of this struct
// is limited to this package), we are not going through the hassle of defining
// an interface just to avoid the usage of bare node{...} (where ... could result in
// an invalid node). We just assume we're good citizen within this package
// and just use newNode function whenever we need a (properly constructed) node
type node struct {
	left        *node
	right       *node
	interval    interval
	midpoint    float64 // midpoint (not necessarily exactly half)
	cardinality int
	potent      bool
}

func newNode(b, e, m float64) *node {
	iv, err := newInterval(b, e)
	if err != nil {
		log.Fatal("could not create interval in newNode")
		return nil
	}
	return &node{
		left:        nil,
		right:       nil,
		interval:    iv,
		midpoint:    m,
		cardinality: 0,
		potent:      false,
	}
}

func (n *node) isActive() bool {
	return n.cardinality > 0 || n.potent
}

func (n *node) setLeft(left *node) *node {
	n.left = left
	return n
}

func (n *node) setRight(right *node) *node {
	n.right = right
	return n
}

func (n *node) setCardinality(c int) {
	n.cardinality = c
}

func buildNode(values []float64, b, e int) *node {
	mid := ((b + e) / 2)
	node := newNode(values[b], values[e], values[mid])
	if e-b == 1 {
		return node
	}
	node.setLeft(buildNode(values, b, mid)).setRight(buildNode(values, mid, e))
	return node
}

func createSegmentTree(values []float64) (*node, error) {
	if len(values) < 2 {
		return nil, errors.New("must get at least two values")
	}
	// make a copy of the slice to leave it unsorted for the caller
	s := append([]float64{}, values...)
	sort.Float64s(s)
	return buildNode(values, 0, len(values)-1), nil
}

func (n *node) promote() {
	n.left.cardinality--
	n.right.cardinality--
	n.cardinality++
}

func (n *node) demote() {
	n.left.cardinality++
	n.right.cardinality++
	n.cardinality--
	n.potent = true
}

func (n *node) isLeaf() bool {
	return n.left == nil && n.right == nil
}

func (n *node) PaddedString(pad int) string {

	padding := ""
	if pad > 0 {
		padding = strings.Repeat(" ", pad)
	}

	s := padding
	s += fmt.Sprint(&n.interval)
	if n.cardinality != 0 {
		s += fmt.Sprintf(" C=%d", n.cardinality)
	}
	if n.potent {
		s += " potent"
	}
	s += "\n"
	if !n.isLeaf() {
		s += n.left.PaddedString(pad + 6)
		s += n.right.PaddedString(pad + 6)
	}
	return s
}

func (n *node) String() string {
	return n.PaddedString(0)
}

func (n *node) goLeft(ival interval) bool {
	return IsStrictlyBelowFloat(ival.begin(), n.midpoint)
}
func (n *node) goRight(ival interval) bool {
	return IsStrictlyBelowFloat(n.midpoint, ival.end())
}

/// Contribution of an edge (b,e) to the final contour
func (n *node) contribution(ival interval, edgeStack *[]interval) {
	if n.cardinality != 0 {
		return
	}
	if n.interval.isFullyContainedIn(ival) && n.potent == false {
		if len(*edgeStack) == 0 {
			*edgeStack = append(*edgeStack, n.interval)
		} else {
			back := &((*edgeStack)[len(*edgeStack)-1])
			if !back.extend(n.interval) {
				// add a new segment if it can not be merged with current one
				*edgeStack = append(*edgeStack, n.interval)
			}
		}
	} else {
		if n.goLeft(ival) {
			n.left.contribution(ival, edgeStack)
		}
		if n.goRight(ival) {
			n.right.contribution(ival, edgeStack)
		}
	}
}

func (n *node) insertInterval(ival interval) {
	if n.interval.isFullyContainedIn(ival) {
		n.cardinality++
	} else {
		if n.goLeft(ival) {
			n.left.insertInterval(ival)
		}
		if n.goRight(ival) {
			n.right.insertInterval(ival)
		}
	}
	n.update()
}

func (n *node) deleteInterval(ival interval) {
	if n.interval.isFullyContainedIn(ival) {
		n.cardinality--
	} else {
		if n.cardinality > 0 {
			n.demote()
		}
		if n.goLeft(ival) {
			n.left.deleteInterval(ival)
		}
		if n.goRight(ival) {
			n.right.deleteInterval(ival)
		}
	}
	n.update()
}

func (n *node) update() {
	if n.left == nil {
		n.potent = false
	} else {
		if n.left.cardinality > 0 && n.right.cardinality > 0 {
			n.promote()
		}
		n.potent = !(!n.left.isActive() && !n.right.isActive())
	}
}
