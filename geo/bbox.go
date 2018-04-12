package geo

import (
	"errors"
	"fmt"
	"math"
)

// BBox describes a simple bounding box.
type BBox interface {
	Xcenter() float64
	Ycenter() float64
	Width() float64
	Height() float64
	Xmin() float64
	Xmax() float64
	Ymin() float64
	Ymax() float64
	Contains(x, y float64) bool
	fmt.Stringer
}

type bbox struct {
	xmin, xmax, ymin, ymax float64
}

var (
	errRightLowerThanLeft = errors.New("xmax coordinate is less than xmin one")
	errTopLowerThanBottom = errors.New("ymax coordinate is less than ymin one")
)

// NewBBox creates a bounding box that is guaranteed
// to be valid if error is nil
func NewBBox(xmin, ymin, xmax, ymax float64) (BBox, error) {

	if xmin >= xmax {
		return nil, errRightLowerThanLeft
	}
	if ymin >= ymax {
		return nil, errTopLowerThanBottom
	}
	return &bbox{xmin, xmax, ymin, ymax}, nil
}

// NewBBoxUnchecked return a bbox that might contain
// garbage (mostly used in test to ease declaration of multiple
// bounding boxes)
func NewBBoxUnchecked(xmin, ymin, xmax, ymax float64) BBox {
	b, _ := NewBBox(xmin, ymin, xmax, ymax)
	return b
}

func (b bbox) Xcenter() float64 {
	return (b.xmin + b.xmax) / 2
}

func (b bbox) Ycenter() float64 {
	return (b.ymin + b.ymax) / 2
}

func (b bbox) Width() float64 {
	return (b.xmax - b.xmin)
}

func (b bbox) Height() float64 {
	return (b.ymax - b.ymin)
}

func (b bbox) Xmin() float64 {
	return b.xmin
}
func (b bbox) Xmax() float64 {
	return b.xmax
}
func (b bbox) Ymin() float64 {
	return b.ymin
}
func (b bbox) Ymax() float64 {
	return b.ymax
}

// Intersect returns the common part of boxes a and b.
func Intersect(a, b BBox) (BBox, error) {
	return NewBBox(math.Max(a.Xmin(), b.Xmin()),
		math.Max(a.Ymin(), b.Ymin()),
		math.Min(a.Xmax(), b.Xmax()),
		math.Min(a.Ymax(), b.Ymax()))
}

// Contains returns true if (x,y) is inside the box
func (b bbox) Contains(x, y float64) bool {
	return IsInRangeFloat64(x, b.xmin, b.xmax) && IsInRangeFloat64(y, b.ymin, b.ymax)
}

func (b bbox) String() string {
	return fmt.Sprintf("bottomLeft: %7.2f,%7.2f topRight: %7.2f,%7.2f",
		b.Xmin(), b.Ymin(), b.Xmax(), b.Ymax())
}

// EqualBBox checks if two boxes are equal.
// For the precision of the comparison see EqualFloat function.
func EqualBBox(a, b BBox) bool {
	return EqualFloat(a.Xmin(), b.Xmin()) &&
		EqualFloat(a.Xmax(), b.Xmax()) &&
		EqualFloat(a.Ymin(), b.Ymin()) &&
		EqualFloat(a.Ymax(), b.Ymax())
}
