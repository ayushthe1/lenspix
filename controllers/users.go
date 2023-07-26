package controllers

import (
	"net/http"
)

type Users struct {
	Templates struct {
		// any type can be assigned to the New field, as long as it has all the methods defined by the Template interface
		New Template
	}
}

func (u Users) New(w http.ResponseWriter, r *http.Request) {
	u.Templates.New.Execute(w, nil)
}
