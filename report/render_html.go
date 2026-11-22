package report

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"html/template"
)

//go:embed template.html
var htmlTemplate string

var parsedHTMLTmpl = template.Must(template.New("report").Parse(htmlTemplate))

// RenderHTML renders a self-contained offline HTML dashboard for r.
func RenderHTML(r *Report) ([]byte, error) {
	jsonBytes, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	err = parsedHTMLTmpl.Execute(&buf, struct {
		RunID    string
		Mode     string
		JSONData template.JS
	}{
		RunID:    r.RunID,
		Mode:     r.Mode,
		JSONData: template.JS(jsonBytes),
	})
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
