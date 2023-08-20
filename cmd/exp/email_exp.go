// Host:
// sandbox.smtp.mailtrap.io
// Port:
// 25 or 465 or 587 or 2525
// Username:
// 5833c25125663f
// Password:
// d12c48afa4d564

package main

import (
	"fmt"

	"github.com/ayushthe1/lenspix/models"
)

func main() {
	smptpConfig := models.SMTPConfig{
		Host:     "sandbox.smtp.mailtrap.io",
		Port:     25,
		Username: "5833c25125663f",
		Password: "d12c48afa4d564",
	}

	// Create an emailService
	emailService := models.NewEmailService(smptpConfig)

	// Define the test message
	email := models.Email{
		From:    "ayushsharmaa101@gmail.com",
		To:      "keshavkumarr07@gmail.com",
		Subject: "Hello, This is a test email",
		HTML:    "<h1>Beilive in Yourself</h1>",
	}

	// send the email
	err := emailService.Send(email)
	if err != nil {
		fmt.Printf("Error sending the email : %w", err)
	}

	fmt.Println("Email sent successfully !")

}

// Create a emailService

// Send the email
