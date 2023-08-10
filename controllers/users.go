package controllers

import (
	"fmt"
	"net/http"

	"github.com/ayushthe1/lenspix/models"
)

type Users struct {
	Templates struct {
		// any type can be assigned to the New field, as long as it has all the methods defined by the Template interface
		New    Template
		SignIn Template
	}
	UserService *models.UserService
}

// handler function for signup route
func (u Users) New(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	// New will be used in get request ,so FormValue will return the "email" value from query string and not from body parameters
	data.Email = r.FormValue("email")

	// if the emailid is present in url, the signup page will have the email field filled with the emailid
	u.Templates.New.Execute(w, data)
	// template.Execute method is used to fill in the placeholders within a template with actual values and generate the final output. Here we are taking the email_id present as query parameter to fill the email field in the signup page (template) which will be rendered (parsed)
}

func (u Users) Create(w http.ResponseWriter, r *http.Request) {

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

	fmt.Fprintf(w, "User created: %+v", user)
}

func (u Users) SignIn(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	// if the emailid is present in url, the signin page will have the email field filled with the emailid
	u.Templates.SignIn.Execute(w, data)
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

	cookie := http.Cookie{
		Name:  "email",
		Value: user.Email,
		Path:  "/", // which paths on the server have access to the cookie
	}
	http.SetCookie(w, &cookie)

	fmt.Fprintf(w, "User authenticated: %+v", user)
}

// function to take a web request and print the current users information.
func (u Users) CurrentUser(w http.ResponseWriter, r *http.Request) {
	email, err := r.Cookie("email")
	if err != nil {
		fmt.Fprint(w, "The email cookie couldn't be read")
		return
	}
	fmt.Fprintf(w, "Email cookie: %s\n", email.Value)
	fmt.Fprintf(w, "Headers: %+v\n", r.Header)
}
