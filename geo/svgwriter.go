package geo

import (
	"fmt"
	"io"
	"math"
)

// SVGWriter is a utility object to create SVG output
// for polygon, contours etc...
type SVGWriter struct {
	width                        int
	styleBuffer                  string
	elements                     []element
	xleft, xright, ybottom, ytop float64
}

func NewSVGWriter(width int) *SVGWriter {
	w := width
	if w <= 0 {
		w = 1024
	}
	return &SVGWriter{width: w}
}

func (w *SVGWriter) assertViewBox() {
	if w.xleft < w.xright && w.ybottom < w.ytop {
		fmt.Printf("view box already correct")
		return
	}
	// loop over all elements to find the boundaries
	xmin := math.MaxFloat64
	xmax := -math.MaxFloat64
	ymin := math.MaxFloat64
	ymax := -math.MaxFloat64
	for _, e := range w.elements {
		b := e.bbox()
		if b.Width() > tiny && b.Height() > tiny {
			xmin = math.Min(xmin, b.Xmin())
			xmax = math.Max(xmax, b.Xmax())
			ymin = math.Min(ymin, b.Ymin())
			ymax = math.Max(ymax, b.Ymax())
		}
	}
	w.xleft = xmin
	w.xright = xmax
	w.ybottom = ymin
	w.ytop = ymax
}

func (w *SVGWriter) Translate(x0, y0 float64) {
	for _, e := range w.elements {
		e.translate(x0, y0)
	}
	w.xleft += x0
	w.xright += x0
	w.ytop += y0
	w.ybottom += y0
}

func (w *SVGWriter) MoveToOrigin() {
	w.assertViewBox()
	w.Translate(-w.xleft, -w.ybottom)
}

func (w *SVGWriter) ViewBox(xleft, xright, ybottom, ytop float64) {
	w.xleft = xleft
	w.xright = xright
	w.ytop = ytop
	w.ybottom = ybottom
}

func (w *SVGWriter) ViewBoxHeight() float64 {
	return w.ytop - w.ybottom
}

func (w *SVGWriter) ViewBoxWidth() float64 {
	return w.xright - w.xleft
}

func (w *SVGWriter) Width() int {
	if w.width > 0 {
		return w.width
	}
	return 1024
}

// Height returns the height of this SVG object
func (w *SVGWriter) Height() int {
	return int(float64(w.Width()) * w.ViewBoxHeight() / w.ViewBoxWidth())
}

// GroupStart starts a group tag with a given classname.
func (w *SVGWriter) GroupStart(classname string) {
	w.appendElement(grouptag(classname))
}

// GroupEnd ends a group tag
func (w *SVGWriter) GroupEnd() {
	w.appendElement(groupendtag{})
}

// Rect adds a rectangle object
func (w *SVGWriter) Rect(x, y, width, height float64) {
	w.appendElement(&rectag{x, y, width, height})
}

// Text adds a text object
func (w *SVGWriter) Text(text string, x, y float64) {
	w.appendElement(&textag{x, y, text})
}

// Polygon adds a polygon object
func (w *SVGWriter) Polygon(p *Polygon) {
	w.PolygonWithClass(p, "")
}

// PolygonWithClass adds a polygon object with a given CSS class
func (w *SVGWriter) PolygonWithClass(p *Polygon, class string) {
	vertices := p.getVertices()
	x := make([]float64, len(vertices))
	y := make([]float64, len(vertices))
	for i, v := range vertices {
		x[i] = v.X
		y[i] = v.Y
	}
	w.appendElement(&poltag{x, y, class})
}

func (w *SVGWriter) Circle(x, y, radius float64) {
	w.appendElement(&cirtag{x, y, radius})
}

// Contour adds one polygon object per sub-contour
func (w *SVGWriter) Contour(c *Contour) {
	for _, p := range *c {
		w.Polygon(&p)
	}
}

// Style add some style
func (w *SVGWriter) Style(style string) {
	w.styleBuffer += style
}

// WriteStyle outputs <style> options
func (w *SVGWriter) WriteStyle(out io.Writer) {
	fmt.Fprintf(out, "<style>%s</style>\n", w.styleBuffer)
}

// WriteSVG output
func (w *SVGWriter) WriteSVG(out io.Writer) {

	w.assertViewBox()

	fmt.Fprintf(
		out,
		"<svg width=\"%v\" height=\"%v\" viewBox=\"%v %v %v %v\">\n",
		w.Width(), w.Height(), w.xleft, w.ybottom, w.xright, w.ytop,
	)
	for _, e := range w.elements {
		fmt.Fprintln(out, e)
	}
	fmt.Fprintf(out, "</svg>\n")
}

// WriteHTML output
func (w *SVGWriter) WriteHTML(out io.Writer) {
	fmt.Fprintf(out, "<html>\n")
	w.WriteStyle(out)
	fmt.Fprintf(out, "<body>\n")
	w.WriteSVG(out)
	fmt.Fprintf(out, "</body>\n</html>\n")
}

type element interface {
	fmt.Stringer
	translate(x0, y0 float64)
	bbox() BBox
}

var (
	tiny       float64 = 1E-9
	tinyBox, _         = NewBBox(-tiny/2.0, -tiny/2.0, tiny/2.0, tiny/2.0)
)

type grouptag string
type groupendtag struct{}
type rectag struct {
	x, y, width, height float64
}
type textag struct {
	x, y float64
	text string
}
type poltag struct {
	x     []float64
	y     []float64
	class string
}
type cirtag struct {
	x, y, radius float64
}

func (g grouptag) String() string {
	return fmt.Sprintf("<g class=\"%s\">", string(g))
}

func (g grouptag) translate(x, y float64) {
}

func (g grouptag) bbox() BBox {
	return tinyBox
}

func (g groupendtag) String() string {
	return "</g>"
}

func (g groupendtag) bbox() BBox {
	return tinyBox
}

func (g groupendtag) translate(x, y float64) {
}

func (r rectag) String() string {
	return fmt.Sprintf("<rect x=\"%v\" y=\"%v\" width=\"%v\" height=\"%v\" /> ", r.x, r.y, r.width, r.height)
}

func (r *rectag) translate(x0, y0 float64) {
	r.x += x0
	r.y += y0
}

func (r *rectag) bbox() BBox {
	b, err := NewBBox(r.x-r.width/2.0, r.y-r.height/2.0, r.x+r.width/2.0, r.y+r.height/2.0)
	if err != nil {
		panic(err)
	}
	return b
}

func (t textag) String() string {
	return fmt.Sprintf("<text x=\"%v\" y=\"%v\">%s</text>", t.x, t.y, t.text)
}

func (t *textag) translate(x0, y0 float64) {
	t.x += x0
	t.y += y0
}

func (t *textag) bbox() BBox {
	//FIXME: is there a way to actually get the size here ??
	return tinyBox
}

func (p poltag) String() string {
	s := fmt.Sprintf("<polygon points=\"")
	for i, _ := range p.x {
		s += fmt.Sprintf("%v,%v ", p.x[i], p.y[i])
	}
	s += "\""
	if len(p.class) > 0 {
		s += " class\"" + p.class + "\""
	}
	s += " />"
	return s
}

func (p *poltag) translate(x0, y0 float64) {
	for i, _ := range p.x {
		p.x[i] += x0
		p.y[i] += y0
	}
}

func (p *poltag) bbox() BBox {
	xmin := math.MaxFloat64
	xmax := -math.MaxFloat64
	ymin := math.MaxFloat64
	ymax := -math.MaxFloat64
	for i, _ := range p.x {
		xmin = math.Min(xmin, p.x[i])
		xmax = math.Max(xmax, p.x[i])
		ymin = math.Min(ymin, p.y[i])
		ymax = math.Max(ymax, p.y[i])
	}
	b, err := NewBBox(xmin, ymin, xmax, ymax)
	if err != nil {
		panic(err)
	}
	return b
}

func (c cirtag) String() string {
	return fmt.Sprintf("<circle cx=\"%f\" cy=\"%f\" r=\"%f\" />", c.x, c.y, c.radius)
}

func (c *cirtag) translate(x0, y0 float64) {
	c.x += x0
	c.y += y0
}

func (c *cirtag) bbox() BBox {
	b, err := NewBBox(c.x-c.radius, c.y-c.radius, c.x+c.radius, c.y+c.radius)
	if err != nil {
		panic(err)
	}
	return b
}

func (w *SVGWriter) appendElement(e element) {
	w.elements = append(w.elements, e)
}
