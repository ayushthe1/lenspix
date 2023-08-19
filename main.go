package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ayushthe1/lenspix/controllers"
	"github.com/ayushthe1/lenspix/migrations"
	"github.com/ayushthe1/lenspix/models"
	"github.com/ayushthe1/lenspix/templates"
	"github.com/ayushthe1/lenspix/views"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
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

	// Setup the database connection
	cfg := models.DefaultPostgresConfig()
	fmt.Println(cfg.String())
	db, err := models.Open(cfg)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Run migrations
	// We no longer need the "migrations" directory variable because our embedding occurs within the migrations directory in fs.go. We can instead pass an empty directory string.
	err = models.MigrateFS(db, migrations.FS, ".")
	// our fs.go file is inside the migration folder. So all the files are going to be relative to that. That's why we pass "." which means current directory w.r.t to the fs.go file.
	if err != nil {
		panic(err)
	}

	// Setup services
	// Setup our model services
	userService := models.UserService{
		DB: db,
	}
	// setup our user service
	sessionService := models.SessionService{
		DB: db,
	}

	// Setup middleware
	umw := controllers.UserMiddleware{
		SessionService: &sessionService,
	}
	csrfKey := "Wc0gT1xfFAjlRwip7l7MmEdjw7DzMXamEHLjyAUP"
	csrfMw := csrf.Protect(
		[]byte(csrfKey),
		//TODO: Fix this before deploying
		csrf.Secure(false),
	)

	// Setup controllers
	userC := controllers.Users{
		UserService:    &userService,
		SessionService: &sessionService,
	}
	userC.Templates.New = views.Must(views.ParseFS(templates.FS, "signup.gohtml", "tailwind.gohtml"))
	userC.Templates.SignIn = views.Must(views.ParseFS(templates.FS, "signin.gohtml", "tailwind.gohtml"))

	// Setup our router and routes

	r := chi.NewRouter()
	// Applying the niddleware
	// These middleware run on all the incoming request
	r.Use(csrfMw)
	r.Use(umw.SetUser)

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

	r.Get("/signup", userC.New)
	r.Get("/signin", userC.SignIn)
	r.Post("/signup", userC.Create)
	r.Post("/signin", userC.ProcessSignIn)
	r.Post("/signout", userC.ProcessSignOut)
	r.Get("/users/me", userC.CurrentUser)
	// r.Get("/users/me", userC.CurrentUser) --> before

	// Apply the router to all routes that match the prefix 'users/me'
	r.Route("/users/me", func(r chi.Router) {
		// RequireUser middleware will be used on all routes with the /users/me prefix
		r.Use(umw.RequireUser)
		r.Get("/", userC.CurrentUser)
		r.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "Hellooo")
		})
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "page not found", http.StatusNotFound)
	})

	// Wrapping csrfMw as a middleware around r and starting the server
	fmt.Println("Starting the server on port :3000 ......")
	http.ListenAndServe(":3000", r)
	// http.ListenAndServe(":3000", csrfMw(umw.SetUser(r)))
	// In this code we are wrapping the router with the middleware that sets a user, then we wrap that whole http handler with the CSRF middleware. This means CSRF protection will run first, then our user lookup, and finally our router will decide what HTTP handler to use based on the path and HTTP method. The order of these can be important.
}

// timer middleware to know the response time for our requests
func TimeMiddleware(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		h(w, r)
		fmt.Println("Request time:", time.Since(start))
	}
}

// ServeMux refers to a simple HTTP request multiplexer (or router) provided by the "net/http" package. It acts as a request router, matching incoming HTTP requests to the corresponding handler functions that should process those requests.
