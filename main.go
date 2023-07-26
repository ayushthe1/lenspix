package main

import (
	"fmt"
	"net/http"

	"github.com/ayushthe1/lenspix/controllers"
	"github.com/ayushthe1/lenspix/templates"
	"github.com/ayushthe1/lenspix/views"
	"github.com/go-chi/chi/v5"
)

// func executeTemplate(w http.ResponseWriter, filepath string) {

// 	htmlTpl, err := views.Parse(filepath)
// 	if err != nil {
// 		log.Println("An error occured while parsing the template")
// 		http.Error(w, "Error parsing the template", http.StatusInternalServerError)
// 		return
// 	}

// 	htmlTpl.Execute(w, nil)
// }

// func homeHandler(w http.ResponseWriter, r *http.Request) {
// 	tplPath := filepath.Join("templates", "home.gohtml")
// 	executeTemplate(w, tplPath)
// }

// type Router struct{}
// func (router Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	switch r.URL.Path {
// 	case "/":
// 		homeHandler(w, r)
// 	case "/contact":
// 		contactHandler(w, r)
// 	default:
// 		http.Error(w, "Page Not Found", http.StatusNotFound)
// 		// w.WriteHeader(http.StatusNotFound)
// 		// fmt.Fprintln(w, "Page not Found")
// 	}
// }

func main() {

	r := chi.NewRouter()

	// parse the template
	tpl := views.Must(views.ParseFS(templates.FS, "home.gohtml", "tailwind.gohtml"))

	r.Get("/", controllers.StaticHandler(tpl))
	// r.Get("/", http.HandlerFunc(homeHandler))

	tpl = views.Must(views.ParseFS(templates.FS, "contact.gohtml", "tailwind.gohtml"))
	r.Get("/contact", controllers.StaticHandler(tpl))

	tpl = views.Must(views.ParseFS(templates.FS, "faq.gohtml", "tailwind.gohtml"))
	r.Get("/faq", controllers.FAQ(tpl))

	// tpl = views.Must(views.ParseFS(templates.FS, "signup.gohtml", "tailwind.gohtml"))
	// r.Get("/signup", controllers.StaticHandler(tpl))

	userC := controllers.Users{}
	userC.Templates.New = views.Must(views.ParseFS(templates.FS, "signup.gohtml", "tailwind.gohtml"))
	r.Get("/signup", userC.New)
	r.Post("/users", userC.Create)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "page not found", http.StatusNotFound)
	})

	fmt.Println("Starting the server on port :3000 ......")
	http.ListenAndServe(":3000", r)
}

// ServeMux refers to a simple HTTP request multiplexer (or router) provided by the "net/http" package. It acts as a request router, matching incoming HTTP requests to the corresponding handler functions that should process those requests.
