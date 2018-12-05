package geo

import (
	"bytes"
	"fmt"
	"log"
	"testing"
)

func TestBasicWrite(t *testing.T) {
	want := `<html>
<style></style>
<body>
<svg width="1024" height="1024" viewBox="0.1 0.1 10.01 10.01">
<g class="test">
<rect x="1" y="2" width="3" height="4" /> 
<text x="10" y="20">some text</text>
<circle cx="10.000000" cy="10.000000" r="0.010000" />
</g>
<g class="points">
<polygon points="0.1,0.1 1.1,0.1 1.1,1.1 2.1,1.1 2.1,3.1 1.1,3.1 1.1,2.1 0.1,2.1 " />
<polygon points="0.1,0.1 1.1,0.1 1.1,1.1 2.1,1.1 2.1,3.1 1.1,3.1 1.1,2.1 0.1,2.1 " class"big" />
</svg>
</body>
</html>
`

	svg := NewSVGWriter(1024)

	testPolygon = Polygon{
		{0.1, 0.1},
		{1.1, 0.1},
		{1.1, 1.1},
		{2.1, 1.1},
		{2.1, 3.1},
		{1.1, 3.1},
		{1.1, 2.1},
		{0.1, 2.1},
		{0.1, 0.1}}

	svg.GroupStart("test")
	svg.Rect(1, 2, 3, 4)
	svg.Text("some text", 10, 20)
	svg.Circle(10, 10, 0.01)
	svg.GroupEnd()

	svg.GroupStart("points")
	svg.Polygon(&testPolygon)
	svg.PolygonWithClass(&testPolygon, "big")

	var buf bytes.Buffer
	svg.WriteHTML(&buf)

	got := string(buf.Bytes()[:])

	fmt.Println(bytes.Compare([]byte(want), []byte(got)))
	if got != want {
		log.Fatalf("Wanted:\n---%v+++and got:\n---%v+++", want, got)
	}

}
