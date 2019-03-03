package output

import (
	"fmt"
	"html/template"
	"io"
	"log"

	"github.com/aphecetche/pigiron/mapping"
	"github.com/aphecetche/pigiron/segcontour"
)

func ToD3(w io.Writer, cseg mapping.CathodeSegmentation, bending bool, showflags segcontour.ShowFlags) {
	bendingString := "Non Bending"
	if bending {
		bendingString = "Bending"
	}

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
	{{ range .JS }}<script src="{{.}}"></script>
	{{ end }}
	<script>
	show({{ .DeID }},{{ .Bending }} )
	</script>
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
		DeID    int
		Bending bool
	}{
		Title: title,
		CSS:   "style.css",
		JS: []string{
			"https://d3js.org/d3.v5.min.js",
			"/static/d3.js",
		},
		DeID:    int(cseg.DetElemID()),
		Bending: bending,
	}
	err = t.Execute(w, data)

	if err != nil {
		log.Fatalf(err.Error())
	}
}
