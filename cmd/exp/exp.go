// package main

// import (
// 	"fmt"

// 	"github.com/ayushthe1/lenspix/models"
// )

// func main() {

// 	cfg := models.DefaultPostgresConfig()
// 	db, err := models.Open(cfg)
// 	if err != nil {
// 		panic(err)
// 	}

// 	defer db.Close()

// 	// Make sure the database is up and running
// 	err = db.Ping()
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println("Connected !")

// 	us := models.UserService{
// 		DB: db,
// 	}
// 	user, err := us.Create("bob@bob.com", "bob123")
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println(user)

	// // Create a table
	// _, err = db.Exec(`
	// CREATE TABLE IF NOT EXISTS users (
	// 	id SERIAL PRIMARY KEY,
	// 	name TEXT,
	// 	email TEXT UNIQUE NOT NULL
	// );

	// CREATE TABLE IF NOT EXISTS orders (
	// 	id SERIAL PRIMARY KEY,
	// 	user_id INT NOT NULL,
	// 	amount INT,
	// 	description TEXT
	// );`)

	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println("tables created")
	// fmt.Println("Connected")

	// // // Insert some data
	// // // using the $1 ,$2 sign and not writing the sql query ourselves protects from sql injection
	// // name := "laporta"
	// // email := "barca@email.com"
	// // row := db.QueryRow(`
	// // INSERT INTO users (name , email)
	// // VALUES ($1, $2) RETURNING id;
	// // `, name, email)

	// // var id int
	// // err = row.Scan(&id) // address of id
	// // if err != nil {
	// // 	panic(err)
	// // }
	// // fmt.Println("users created. id=", id)

	// // query a single record
	// id := 1
	// row := db.QueryRow(`
	// SELECT name, email
	// FROM users
	// WHERE id=$1;`, id)
	// var name, email string
	// err = row.Scan(&name, &email)

	// if err == sql.ErrNoRows {
	// 	fmt.Println("Error, no rows")
	// 	return
	// }
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Printf("User information: name=%s, email=%s \n", name, email)

	// // // Inserting fake orders
	// // userID := 1
	// // for i := 1; i <= 5; i++ {
	// // 	amount := i * 100
	// // 	desc := fmt.Sprintf("Fake order #%d", i)
	// // 	_, err := db.Exec(`
	// // 	INSERT INTO orders(user_id, amount, description)
	// // 	VALUES ($1,$2,$3)`, userID, amount, desc)

	// // 	if err != nil {
	// // 		panic(err)
	// // 	}
	// // }
	// // fmt.Println("Created fake orders")

	// //
	// type Order struct {
	// 	ID          int
	// 	UserID      int
	// 	Amount      int
	// 	Description string
	// }

	// var orders []Order
	// userID := 1

	// rows, err := db.Query(`
	// SELECT id, amount, description
	// FROM orders
	// WHERE user_id=$1`, userID)
	// if err != nil {
	// 	panic(err)
	// }
	// defer rows.Close()

	// // after calling next ,rows will point to the first record
	// for rows.Next() {
	// 	var order Order
	// 	order.UserID = userID
	// 	err := rows.Scan(&order.ID, &order.Amount, &order.Description)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	orders = append(orders, order)

	// }
	// err = rows.Err()
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println("Orders:", orders)
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
