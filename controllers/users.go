package controllers

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/ayushthe1/lenspix/context"
	"github.com/ayushthe1/lenspix/errors"
	"github.com/ayushthe1/lenspix/models"
)

type Users struct {
	Templates struct {
		// any type can be assigned to the New field, as long as it has all the methods defined by the Template interface
		New            Template
		SignIn         Template
		ForgotPassword Template
		CheckYourEmail Template
		ResetPassword  Template
	}
	UserService          *models.UserService
	SessionService       *models.SessionService
	PasswordResetService *models.PasswordResetService
	EmailService         *models.EmailService
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

	// NEVER TRY TO WRITE TO RESPONSEWRITER LIKE THIS FOR DEBUGGING AS WE CAN ONLY WRITE TO W ONCE. THESE LINES WILL INTRODUCE A BUG.
	// fmt.Fprint(w, "Email: ", r.FormValue("email"))
	// fmt.Fprint(w, "Password: ", r.FormValue("password"))

	// For GET requests, the server processes the data in the URL's query parameters. For POST requests, the server retrieves the encoded data from the request's body.

	var data struct {
		Email    string
		Password string
	}

	data.Email = r.FormValue("email")
	data.Password = r.FormValue("password")
	user, err := u.UserService.Create(data.Email, data.Password)
	if err != nil {
		if errors.Is(err, models.ErrEmailTaken) {
			err = errors.Public(err, "That email address is already associated with an account.")
		}
		//  execute the signup template ,and render it with passed in data
		u.Templates.New.Execute(w, r, data, err)
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
	http.Redirect(w, r, "/galleries", http.StatusFound)

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
		// http.Error(w, "Something went wrong", http.StatusInternalServerError)
		http.Redirect(w, r, "/signin", http.StatusFound)
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
	http.Redirect(w, r, "/galleries", http.StatusFound)
}

// function to take the token store in the cookie and take that to lookup the current user
// function to take a web request and print the current users information.
func (u Users) CurrentUser(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	user := context.User(ctx)
	// We don't technically need to  check if the user is nil bcoz we're assuming that RequireUser() middleware has been run.
	if user == nil {
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}
	fmt.Fprintf(w, "Cureent user: %s\n", user.Email)

	// token, err := readCookie(r, CookieSession)
	// if err != nil {
	// 	fmt.Println(err)
	// 	log.Println("Couldn't read token from cookie ..Redirecting to signIn page")
	// 	http.Redirect(w, r, "/signin", http.StatusFound)
	// 	return
	// }
	// user, err := u.SessionService.User(token)
	// if err != nil {
	// 	fmt.Println(err)
	// 	log.Println("token in cookie isn't valid ..Redirecting to signIn page")
	// 	http.Redirect(w, r, "/signin", http.StatusFound)
	// 	return
	// }
	// fmt.Fprintf(w, "Current user: %s\n", user.Email)

	// fmt.Fprintf(w, "Good Luck !")
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

// HTTP handler to render the password reset form.
// Like the sign in form, it can parse any URL query params and use those to prefill the form.
func (u Users) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}

	data.Email = r.FormValue("email")              // get the email form url query parameters
	u.Templates.ForgotPassword.Execute(w, r, data) // prefil the page with email if it exists otherwise empty string is filled in password reset form

}

// handler for processing the forgot password form
func (u Users) ProcessForgotPassword(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}

	data.Email = r.FormValue("email") // get the email from form post request

	// using the password reset service to create a new password reset token
	pwReset, err := u.PasswordResetService.Create(data.Email)
	if err != nil {
		// handle other cases
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	vals := url.Values{
		"token": {pwReset.Token},
	}
	// encode the token in url
	resetURL := "https://www.lenspix.com/reset-pw?" + vals.Encode()

	// send the mail to the user with the reset_url
	err = u.EmailService.ForgotPassword(data.Email, resetURL)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	//Render the 'check your email' page
	u.Templates.CheckYourEmail.Execute(w, r, data)
}

// handler to render the password reset form
func (u Users) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Token string
	}
	data.Token = r.FormValue("token") // parsing the token from the URL query parameters. This token will be inserted into the form as a hidden value
	u.Templates.ResetPassword.Execute(w, r, data)
}

// handler to process the password reset form
func (u Users) ProcessResetPassword(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Token    string
		Password string // user's new password they want to update
	}
	data.Token = r.FormValue("token")       // token value from form via post request
	data.Password = r.FormValue("password") // password value from form via post request

	// Attempt to consume the token
	user, err := u.PasswordResetService.Consume(data.Token)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	// Update the user's password
	err = u.UserService.UpdatePassword(user.ID, data.Password)
	if err != nil {
		http.Error(w, "Something went wromg", http.StatusInternalServerError)
		return
	}

	// Sign the user in now that they have reset their password.
	// Any errors from this point onward should redirect to the sign in page.
	// Create a new session
	session, err := u.SessionService.Create(user.ID)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}
	// sign the user in
	setCookie(w, CookieSession, session.Token)
	http.Redirect(w, r, "/users/me", http.StatusFound)

}

// Defining all the middlewares

type UserMiddleware struct {
	SessionService *models.SessionService
}

// A middleware function to look up a user if one can be found and to store it in the request context . It accepts an http handler as an argument, and returns a new http handler.
func (umw UserMiddleware) SetUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Lookup the session token via the users cookies. If we run into an error reading it, proceed with the request. The goal of this middleware isn't to limit access. It only sets the user in the context if it can.
		token, err := readCookie(r, CookieSession)
		if err != nil {
			// If there isn't a cookie or it can't get it, we will proceed with the request and assume that the user is not logged in.
			next.ServeHTTP(w, r) // Continue with the wrapped HTTP handler without setting the user.
			return
		}

		// Query for a valid session with the token
		user, err := umw.SessionService.User(token)
		if err != nil {
			// Invalid or expired token. In either case we can still proceed, we just cannot set a user.
			next.ServeHTTP(w, r) // Continue with the wrapped HTTP handler without setting the user.
			return
		}

		// Store the user associated with the session in the context
		// user has been found ,we will get the context ,set the value and then update the request with the new context set
		ctx := r.Context()
		ctx = context.WithUser(ctx, user) // updating the context
		r = r.WithContext(ctx)            // updating the request and get the new request with our updated context

		// Finally we call the handler that our middleware was applied to with the updated request.
		next.ServeHTTP(w, r) // Continue with the wrapped HTTP handler WITH a user being set.
	})
}

//  middleware that requires a user to be signed in, and otherwise redirects them to the sign in page.
// This middleware assumes that we have already run our SetUser middleware
func (umw UserMiddleware) RequireUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check to see if a user is present and if they're not present ,redirect them to the signin page
		user := context.User(r.Context())
		if user == nil {
			http.Redirect(w, r, "/signin", http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}
