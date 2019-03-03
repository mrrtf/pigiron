package output

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"log"

	"github.com/aphecetche/pigiron/geo"
	"github.com/aphecetche/pigiron/mapping"
	"github.com/aphecetche/pigiron/segcontour"

	// must include the specific implementation package of the mapping
	_ "github.com/aphecetche/pigiron/mapping/impl4"
)

func ToSVG(w io.Writer, cseg mapping.CathodeSegmentation, bending bool, showflags segcontour.ShowFlags) {
	svg := geo.NewSVGWriter(1024)
	segcontour.SVGSegmentation(cseg, svg, showflags)

	svg.MoveToOrigin()

	bendingString := "Non Bending"
	if bending {
		bendingString = "Bending"
	}

	var buf bytes.Buffer
	svg.WriteSVG(&buf)
	svgString := string(buf.Bytes()[:])

	title := fmt.Sprintf("Cathode Segmentation DE %d %s", cseg.DetElemID(), bendingString)

	const tpl = `
<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<title>{{.Title}}</title>
		<link rel="stylesheet" type="text/css" href="/static/{{.CSS}}" />
	</head>
	<body>
	<h1>{{.Title}}</h1>
	{{.Content}}
	{{ range .JS }}<script src="{{.}}"></script>
	{{ end }}
	</body>
</html>
`

	t, err := template.New("webpage").Parse(tpl)
	if err != nil {
		log.Fatalf(err.Error())
	}
	data := struct {
		Title   string
		CSS     string
		JS      []string
		Content template.HTML
	}{
		Title: title,
		CSS:   "style.css",
		JS: []string{
			"https://d3js.org/d3.v5.min.js",
			"/static/svg.js",
		},
		Content: template.HTML(svgString),
	}
	err = t.Execute(w, data)

	if err != nil {
		log.Fatalf(err.Error())
	}
}
