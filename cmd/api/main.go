package main

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v4/stdlib"
)

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	SSLMode  string
}

//  fmt.Printf is used to print formatted strings to the standard output, while fmt.Sprintf is used to format strings and store the result in a new string variable. Both functions use format strings and arguments, but they differ in their output destinations.
func (cfg PostgresConfig) String() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database, cfg.SSLMode)
}

func main() {

	cfg := PostgresConfig{
		Host:     "localhost",
		Port:     "5432",
		User:     "baloo",
		Password: "junglebook",
		Database: "lenspix",
		SSLMode:  "disable",
	}

	db, err := sql.Open("pgx", cfg.String())
	if err != nil {
		panic(err)
	}

	defer db.Close()

	// Make sure the database is up and running
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected")
}

// func main() {
// 	// the file path that we're passing in is relative to where we're running the file from
// 	t, err := template.ParseFiles("hello.gohtml")
// 	if err != nil {
// 		panic(err)
// 	}

// 	user := User{
// 		Name: "John Smith",
// 		Age:  44,
// 	}

// 	// execute is taking the template and processing it
// 	err = t.Execute(os.Stdout, user)
// 	if err != nil {
// 		panic(err)
// 	}
// }
