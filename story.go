package cyoa

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
)

var tpl *template.Template

var defaultHandlerTmpl = `<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<title>Choose your own adventure</title>
</head>
<body>
	<h1>{{.Title}}</h1>
	{{range .Paragraphs}}
		<p>{{.}}</p>
	{{end}}
	<ul>
		{{range .Options}}
		<li><a href="/{{.Chapter}}">{{.Text}}</a></li>
		{{end}}
	</ul>
</body>
</html>`

func init() {
	tpl = template.Must(template.New("").Parse(defaultHandlerTmpl))

}

func StoryToJson(r io.Reader) (Story, error) {
	d := json.NewDecoder(r)

	var story Story
	if err := d.Decode(&story); err != nil {
		return nil, err
	}
	return story, nil
}

type handler struct {
	s Story
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if path == "" || path == "/" {
		path = "/intro"
	}

	path = path[1:]

	if chapter, ok := h.s[path]; ok {
		err := tpl.Execute(w, chapter)
		if err != nil {
			log.Println(err)
			http.Error(w, "Something failed", http.StatusInternalServerError)
		}
		return
	}
	http.Error(w, "Chapter not found", http.StatusNotFound)

}

func NewHandler(s Story) http.Handler {
	return handler{s}
}

type Story map[string]Chapter

type Chapter struct {
	Title      string   `json:"title,omitempty"`
	Paragraphs []string `json:"story,omitempty"`
	Options    []Option `json:"options,omitempty"`
}

type Option struct {
	Text    string `json:"text,omitempty"`
	Chapter string `json:"arc,omitempty"`
}
