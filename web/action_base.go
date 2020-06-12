package web

import (
	"encoding/json"
	"encoding/xml"
	"html/template"
	"net/http"
	"path"
)

// ActionBase provides basic methods of service action.
type ActionBase struct {
	Logger         *Logger
	ResponseWriter http.ResponseWriter
	Request        *http.Request
}

func (p *ActionBase) setLogger(logger *Logger) {
	p.Logger = logger
}

func (p *ActionBase) setResponseWriter(w http.ResponseWriter) {
	p.ResponseWriter = w
}

func (p *ActionBase) setRequest(r *http.Request) {
	p.Request = r
}

// FormValue returns first value of key from request.
func (p *ActionBase) FormValue(key string) string {
	return p.Request.FormValue(key)
}

// SetContentType sets content type of response.
func (p *ActionBase) SetContentType(contentType string) {
	p.ResponseWriter.Header().Set("Content-Type", contentType)
}

// WriteError writes error to response.
func (p *ActionBase) WriteError(err error) {
	writeError(err, http.StatusInternalServerError, p.ResponseWriter, p.Logger)
}

func (p *ActionBase) Write(data []byte) {
	_, err := p.ResponseWriter.Write(data)
	if err != nil {
		p.WriteError(err)
	}
}

// WriteHTML writes html text.
func (p *ActionBase) WriteHTML(html string) {
	p.SetContentType("text/html")
	p.Write([]byte(html))
}

// WriteJSON writes the JSON encoding of v.
func (p *ActionBase) WriteJSON(v interface{}) {
	data, err := json.Marshal(v)
	if err != nil {
		p.WriteError(err)
		return
	}
	p.SetContentType("application/json")
	p.Write(data)
}

// WriteXML writes the XML encoding of v.
func (p *ActionBase) WriteXML(v interface{}) {
	data, err := xml.Marshal(v)
	if err != nil {
		p.WriteError(err)
		return
	}
	p.SetContentType("text/xml")
	p.Write(data)
}

// Forward forwards data to template.
func (p *ActionBase) Forward(data interface{}, funcMap template.FuncMap, filenames ...string) {
	if len(filenames) == 0 {
		return
	}
	tmpl := template.New(path.Base(filenames[0]))
	tmpl.Funcs(funcMap)
	_, err := tmpl.ParseFiles(filenames...)
	if err != nil {
		p.WriteError(err)
		return
	}
	err = tmpl.Execute(p.ResponseWriter, data)
	if err != nil {
		p.WriteError(err)
	}
}
