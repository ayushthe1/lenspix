package main

import (
	"fmt"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	for i, arg := range os.Args {
		fmt.Println(i, arg)
	}

	if len(os.Args) < 2 {
		fmt.Println("No cli arguments provided.")
		return
	}

	switch os.Args[1] {
	case "hash":
		// hash the password
		hash(os.Args[2])
	case "compare":
		if len(os.Args) < 4 {
			fmt.Println("Not provided three cli arguments that are needed with compare.")
			return
		}
		compare(os.Args[2], os.Args[3])
	default:
		fmt.Println("PLease provide valid inputs")
	}

}

func hash(password string) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Printf("error hashing: %v \n", password)
		return
	}
	fmt.Println(string(hashedBytes))
}

func compare(password, hash string) {

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	if err != nil {
		fmt.Println("Password is not equal to hash value")
		return
	}

	fmt.Println("Password matches Hash value !")
}
