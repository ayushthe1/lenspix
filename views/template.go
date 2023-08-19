package views

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"log"
	"net/http"

	"github.com/ayushthe1/lenspix/context"

	"github.com/ayushthe1/lenspix/models"
	"github.com/gorilla/csrf"
)

type Template struct {
	htmlTpl *template.Template
}

func Must(t Template, err error) Template {
	if err != nil {
		panic(err)
	}
	return t
}

func Parse(filepath string) (Template, error) {
	tpl, err := template.ParseFiles(filepath)
	if err != nil {
		return Template{}, fmt.Errorf("parsing template: %w", err)
	}

	return Template{htmlTpl: tpl}, nil
}

// Stubbed function for parsing
func ParseFS(fs fs.FS, patterns ...string) (Template, error) {
	// define the csrField function before we parse the file
	tpl := template.New(patterns[0]) // create an empty template.Template
	tpl = tpl.Funcs(                 // provide our custom placeholder function
		template.FuncMap{
			"csrfField": func() (template.HTML, error) {
				return `<!-- TODO: Implement the csrField -->`, fmt.Errorf("csrfField not implemented")
			},
			"currentUser": func() (template.HTML, error) {
				return "", fmt.Errorf("current user not implemented ")
			},
		},
	)

	tpl, err := tpl.ParseFS(fs, patterns...) // using ParseFS template provided by the standard library
	if err != nil {
		return Template{}, fmt.Errorf("parsing template: %w", err)
	}

	return Template{htmlTpl: tpl}, nil
}

// Actual function implementation with the request
func (t Template) Execute(w http.ResponseWriter, r *http.Request, data interface{}) {

	tpl, err := t.htmlTpl.Clone() // clone() to avoid race conditions
	if err != nil {
		log.Printf("Cloning template: %v", err)
		http.Error(w, "There was an error rendering the page", http.StatusInternalServerError)
		return
	}

	// pass in a new template.FuncMap with the real csrfField & currentUser implementation
	tpl = tpl.Funcs( // provide our custom function
		template.FuncMap{ // update the placeholder function
			"csrfField": func() template.HTML {
				return csrf.TemplateField(r)
			},
			"currentUser": func() *models.User {
				return context.User(r.Context())
			},
		},
	)

	w.Header().Set("Content-Type", "text/html")
	var buf bytes.Buffer

	// when we start writing to the response writer's body, if there's isn't a status code set, it automatically gets set to 200.
	err = tpl.Execute(&buf, data) // execute the template and write the data as it goes through to the response writer
	if err != nil {
		log.Printf("executing template: %v", err)
		http.Error(w, "There was an error executing the template", http.StatusInternalServerError)
		return
	}

	io.Copy(w, &buf)
}
