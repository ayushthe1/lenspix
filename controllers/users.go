package controllers

import (
	"fmt"
	"net/http"
)

type Users struct {
	Templates struct {
		// any type can be assigned to the New field, as long as it has all the methods defined by the Template interface
		New Template
	}
}

func (u Users) New(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	// New will be used in get request ,so FormValue will return the "email" value from query string
	data.Email = r.FormValue("email")
	u.Templates.New.Execute(w, data)
}

func (u Users) Create(w http.ResponseWriter, r *http.Request) {

	fmt.Fprint(w, "Email: ", r.FormValue("email"))
	fmt.Fprint(w, "Password: ", r.FormValue("password"))
}
