package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/ayushthe1/lenspix/models"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load() // load the env variables from .env file and make them available to code through os.Getenv()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	host := os.Getenv("SMTP_HOST")
	portStr := os.Getenv("SMTP_PORT")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		panic(err)
	}

	username := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")

	smptpConfig := models.SMTPConfig{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
	}

	// Create an emailService
	emailService := models.NewEmailService(smptpConfig)

	// send the email
	err = emailService.ForgotPassword("mrbeast@gmail.com", "https://verocios.com")
	if err != nil {
		fmt.Printf("Error sending the email : %s", err.Error())
		panic(err)
	}

	fmt.Println("Email sent successfully !")

}
