package geo

import (
	"fmt"
	"io"
)

// SVGWriter is a utility object to create SVG output
// for polygon, contours etc...
type SVGWriter struct {
	Width int
	BBox
	buffer      string
	styleBuffer string
}

// Height returns the height of this SVG object
func (w *SVGWriter) Height() int {
	return int(float64(w.Width) * w.BBox.Height() / w.BBox.Width())
}

// GroupStart starts a group tag with a given classname.
func (w *SVGWriter) GroupStart(classname string) {
	w.buffer += fmt.Sprintf("<g class=\"%s\">\n", classname)
}

// GroupEnd ends a group tag
func (w *SVGWriter) GroupEnd() {
	w.buffer += "</g>\n"
}

// Rect adds a rectangle object
func (w *SVGWriter) Rect(x, y, width, height float64) {
	w.buffer += fmt.Sprintf("<ref x=\"%v\" y=\"%v\" width=\"%v\" height=\"%v\" /> ", x, y, width, height)
}

// Text adds a text object
func (w *SVGWriter) Text(text string, x, y float64) {
	w.buffer += fmt.Sprintf("<text x=\"%v\" y=\"%f\">%s</text>", x, y, text)
}

// Polygon adds a polygon object
func (w *SVGWriter) Polygon(p *Polygon) {
	w.buffer += fmt.Sprintf("<polygon points=\"")
	for _, v := range p.getVertices() {
		w.buffer += fmt.Sprintf("%v,%v ", v.X, v.Y)
	}
	w.buffer += "\"/>\n"
}

func (w *SVGWriter) Circle(x, y, radius float64) {
	w.buffer += fmt.Sprintf("<circle cx=\"%f\" cy=\"%f\" r=\"%f\" />", x, y, radius)
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
	fmt.Fprintf(
		out,
		"<svg width=\"%v\" height=\"%v\" viewBox=\"%v %v %v %v\">\n",
		w.Width, w.Height(), w.BBox.Xmin(), w.BBox.Ymin(), w.BBox.Xmax(), w.BBox.Ymax(),
	)
	fmt.Fprintf(out, "%s\n</svg>\n", w.buffer)
}

// WriteHTML output
func (w *SVGWriter) WriteHTML(out io.Writer) {
	fmt.Fprintf(out, "<html>\n")
	w.WriteStyle(out)
	fmt.Fprintf(out, "<body>\n")
	w.WriteSVG(out)
	fmt.Fprintf(out, "</body>\n</html>\n")
}
