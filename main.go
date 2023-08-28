package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/ayushthe1/lenspix/controllers"
	"github.com/ayushthe1/lenspix/migrations"
	"github.com/ayushthe1/lenspix/models"
	"github.com/ayushthe1/lenspix/templates"
	"github.com/ayushthe1/lenspix/views"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
	"github.com/joho/godotenv"
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

type config struct {
	PSQL models.PostgresConfig
	SMTP models.SMTPConfig
	CSRF struct {
		Key    string
		Secure bool
	}
	Server struct {
		Address string
	}
}

func loadEnvConfig() (config, error) {
	var cfg config
	err := godotenv.Load()
	if err != nil {
		return cfg, err
	}
	//TODO: Read PSQL from an env variable
	cfg.PSQL = models.DefaultPostgresConfig()
	//TODO: Read SMTP from an env variable
	cfg.SMTP.Host = os.Getenv("SMTP_HOST")
	portStr := os.Getenv("SMTP_PORT")

	cfg.SMTP.Port, err = strconv.Atoi(portStr)
	if err != nil {
		return cfg, err
	}

	cfg.SMTP.Username = os.Getenv("SMTP_USERNAME")
	cfg.SMTP.Password = os.Getenv("SMTP_PASSWORD")

	//TODO: Read the CSRF values from an env variable
	cfg.CSRF.Key = "Wc0gT1xfFAjlRwip7l7MmEdjw7DzMXamEHLjyAUP"
	cfg.CSRF.Secure = false

	//TODO: Read the server values from an Env variable
	cfg.Server.Address = ":3000"

	return cfg, nil
}

func main() {

	cfg, err := loadEnvConfig()
	if err != nil {
		panic(err)
	}

	// Setup the database connection
	db, err := models.Open(cfg.PSQL)
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
	userService := &models.UserService{
		DB: db,
	}
	// setup our user service
	sessionService := &models.SessionService{
		DB: db,
	}
	// setup password reset service
	pwResetService := &models.PasswordResetService{
		DB: db,
	}
	// setup email service
	emailService := models.NewEmailService(cfg.SMTP)
	// setup gallery service
	galleryService := &models.GalleryService{
		DB: db,
	}

	// Setup middleware
	umw := controllers.UserMiddleware{
		SessionService: sessionService,
	}

	csrfMw := csrf.Protect(
		[]byte(cfg.CSRF.Key),
		csrf.Secure(cfg.CSRF.Secure),
		// By default, the CSRF library sets the path attribute based on the current URL.
		// To make the CSRF cookie work, on every page regardless of what the current path is ,csrf should always use the path "/"
		csrf.Path("/"),
	)

	// Setup controllers
	userC := controllers.Users{
		UserService:          userService,
		SessionService:       sessionService,
		PasswordResetService: pwResetService,
		EmailService:         emailService,
	}
	userC.Templates.New = views.Must(views.ParseFS(templates.FS, "signup.gohtml", "tailwind.gohtml"))
	userC.Templates.SignIn = views.Must(views.ParseFS(templates.FS, "signin.gohtml", "tailwind.gohtml"))
	userC.Templates.ForgotPassword = views.Must(views.ParseFS(
		templates.FS,
		"forgot-pw.gohtml", "tailwind.gohtml",
	))
	userC.Templates.CheckYourEmail = views.Must(views.ParseFS(
		templates.FS,
		"check-your-email.gohtml", "tailwind.gohtml",
	))
	userC.Templates.ResetPassword = views.Must(views.ParseFS(
		templates.FS,
		"reset-pw.gohtml", "tailwind.gohtml",
	))

	galleriesC := controllers.Galleries{
		GalleryService: galleryService,
	}
	galleriesC.Templates.New = views.Must(views.ParseFS(
		templates.FS,
		"galleries/new.gohtml", "tailwind.gohtml",
	))
	galleriesC.Templates.Edit = views.Must(views.ParseFS(
		templates.FS,
		"galleries/edit.gohtml", "tailwind.gohtml",
	))
	galleriesC.Templates.Index = views.Must(views.ParseFS(
		templates.FS,
		"galleries/index.gohtml", "tailwind.gohtml",
	))

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
	r.Get("/forgot-pw", userC.ForgotPassword)
	r.Post("/forgot-pw", userC.ProcessForgotPassword)
	r.Get("/users/me", userC.CurrentUser)
	r.Get("/reset-pw", userC.ResetPassword)
	r.Post("/reset-pw", userC.ProcessResetPassword)
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

	r.Route("/galleries", func(r chi.Router) {
		r.Get("/{id}", galleriesC.Show) // This route will be outside the Group because we don’t want it to require a user
		r.Group(func(r chi.Router) {
			// This middleware will apply on these group of routes
			r.Use(umw.RequireUser) // middleware to ensure only signed in user can access this page
			r.Get("/new", galleriesC.New)
			r.Get("/{id}/edit", galleriesC.Edit)
			r.Post("/{id}", galleriesC.Update)
			r.Post("/", galleriesC.Create)
			r.Get("/", galleriesC.Index)
		})

	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "page not found", http.StatusNotFound)
	})

	// Wrapping csrfMw as a middleware around r and starting the server
	fmt.Printf("Starting the server on port :%s ......", cfg.Server.Address)
	err = http.ListenAndServe(cfg.Server.Address, r)
	if err != nil {
		panic(err)
	}
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
