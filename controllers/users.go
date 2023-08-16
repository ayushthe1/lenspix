package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ayushthe1/lenspix/models"
)

type Users struct {
	Templates struct {
		// any type can be assigned to the New field, as long as it has all the methods defined by the Template interface
		New    Template
		SignIn Template
	}
	UserService    *models.UserService
	SessionService *models.SessionService
}

// handler function for signup route
func (u Users) New(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
		// CSRFField template.HTML
	}
	// New will be used in get request ,so FormValue will return the "email" value from query string and not from body parameters
	data.Email = r.FormValue("email")

	// // give us the HTML for a hidden <input> tag that has the CSRF token for the incoming request.
	// data.CSRFField = csrf.TemplateField(r)

	// if the emailid is present in url, the signup page will have the email field filled with the emailid
	//csrf token will be added to the signup page form

	u.Templates.New.Execute(w, r, data) // data struct will be available inside of our template.

	// template.Execute method is used to fill in the placeholders within a template with actual values and generate the final output. Here we are taking the email_id present as query parameter to fill the email field in the signup page (template) which will be rendered (parsed)
}

func (u Users) Create(w http.ResponseWriter, r *http.Request) {

	// For GET requests, the server processes the data in the URL's query parameters. For POST requests, the server retrieves the encoded data from the request's body.

	fmt.Fprint(w, "Email: ", r.FormValue("email"))
	fmt.Fprint(w, "Password: ", r.FormValue("password"))

	email := r.FormValue("email")
	password := r.FormValue("password")
	user, err := u.UserService.Create(email, password)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	log.Printf("User created: %+v", user)

	// create a session after the user is created
	session, err := u.SessionService.Create(user.ID)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}

	// set a cookie
	setCookie(w, CookieSession, session.Token)
	http.Redirect(w, r, "/users/me", http.StatusFound)

}

func (u Users) SignIn(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")

	// if the emailid is present as url query parameter, the signin page will have the email field filled with the emailid

	u.Templates.SignIn.Execute(w, r, data)
}

// Handler for processing the sign in form
func (u Users) ProcessSignIn(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email    string
		Password string
	}
	data.Email = r.FormValue("email")
	data.Password = r.FormValue("password")
	user, err := u.UserService.Authenticate(data.Email, data.Password)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	log.Printf("User authenticated: %+v", user)

	// create a session
	session, err := u.SessionService.Create(user.ID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	// cookie := http.Cookie{
	// 	Name:     "MeraLoginSession",
	// 	Value:    session.Token,
	// 	Path:     "/",  // which paths on the server have access to the cookie
	// 	HttpOnly: true, // cookies should be only accessible via http browser request and not javascript request(securing cookies from XSS)
	// }
	// http.SetCookie(w, &cookie)

	setCookie(w, CookieSession, session.Token)
	http.Redirect(w, r, "/users/me", http.StatusFound)
}

// function to take the token store in the cookie and take that to lookup the current user
// function to take a web request and print the current users information.
func (u Users) CurrentUser(w http.ResponseWriter, r *http.Request) {
	token, err := readCookie(r, CookieSession)
	if err != nil {
		fmt.Println(err)
		log.Println("Couldn't read token from cookie ..Redirecting to signIn page")
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}
	user, err := u.SessionService.User(token)
	if err != nil {
		fmt.Println(err)
		log.Println("token in cookie isn't valid ..Redirecting to signIn page")
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}
	fmt.Fprintf(w, "Current user: %s\n", user.Email)

	fmt.Fprintf(w, "Good Luck !")
}

func (u Users) ProcessSignOut(w http.ResponseWriter, r *http.Request) {
	token, err := readCookie(r, CookieSession)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}

	// delete the user's session
	err = u.SessionService.Delete(token)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "SOmething went wrong", http.StatusInternalServerError)
		return
	}
	// delete the user's cookie
	deleteCookie(w, CookieSession)
	http.Redirect(w, r, "/signin", http.StatusFound)
}
