package geo

import (
	"log"
	"testing"
)

func TestNewBBoxMustErrorIfCoordinatesAreInvalid(t *testing.T) {

	_, err1 := NewBBox(2, 2, 0, 3)
	_, err2 := NewBBox(2, 2, 3, 0)
	if err1 == nil || err2 == nil {
		t.Fatalf("BBox creation should fail if coordinates are invalid")
	}
}

var b, _ = NewBBox(-15, -10, 5, 20)

func TestBBoxBoundaries(t *testing.T) {
	if b.Xmin() != -15 {
		t.Errorf("expected xmin -15 and got %f", b.Xmin())
	}
	if b.Xmax() != 5 {
		t.Errorf("expected xmxa 5 and got %f", b.Xmax())
	}
	if b.Ymin() != -10 {
		t.Errorf("expected ymin -10 and got %f", b.Ymin())
	}
	if b.Ymax() != 20 {
		t.Errorf("expected ymax 20 and got %f", b.Ymax())
	}
}

func TestBBoxCenter(t *testing.T) {
	if b.Xcenter() != -5 {
		t.Errorf("expected xcenter -5 and got %f", b.Xcenter())
	}
	if b.Ycenter() != 5 {
		t.Errorf("expected ycenter 5 and got %f", b.Ycenter())
	}
}
func TestBBoxWH(t *testing.T) {
	if b.Width() != 20 {
		t.Errorf("expected width 20 and got %f", b.Width())
	}
	if b.Height() != 30 {
		t.Errorf("expected height 30 and got %f", b.Height())
	}
}

func TestBBoxString(t *testing.T) {

	expected := "bottomLeft: -15.000000,-10.000000 topRight: 5.000000,20.000000"
	if b.String() != expected {
		t.Errorf("expected string:%s but got:%s",
			expected, b.String())
	}
}
func TestBBoxIntersect(t *testing.T) {
	one, _ := NewBBox(0, 0, 4, 2)
	two, _ := NewBBox(2, -1, 5, 1)
	expected, _ := NewBBox(2, 0, 4, 1)
	inter, _ := Intersect(one, two)
	if !EqualBBox(inter, expected) {
		t.Errorf("Intersect not as expected : %s", inter.String())
		log.Println("expected", expected)
		log.Println("inter", inter)
	}
}
