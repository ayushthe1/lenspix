package main

import (
	"os"
	"text/template"
)

type User struct {
	Name string
	Age  int
}

func main() {
	// the file path that we're passing in is relative to where we're running the file from
	t, err := template.ParseFiles("hello.gohtml")
	if err != nil {
		panic(err)
	}

	user := User{
		Name: "John Smith",
		Age:  44,
	}

	// execute is taking the template and processing it
	err = t.Execute(os.Stdout, user)
	if err != nil {
		panic(err)
	}
}
