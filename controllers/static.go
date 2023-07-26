package controllers

import (
	"html/template"
	"net/http"
)

// closure is a nested function that helps us access the outer function's variables even after the outer function is closed
func StaticHandler(tpl Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, nil)
	}
}

func FAQ(tpl Template) http.HandlerFunc {
	questions := []struct {
		Question string
		Answer   template.HTML
	}{
		{Question: "Is there a free version?",
			Answer: "Yes! We have a free trial offer",
		},
		{
			Question: "What is your name ?",
			Answer:   "My name is Ayush",
		},
	}

	return func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, questions)
	}
}
