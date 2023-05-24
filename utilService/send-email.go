package utilService

import (
	"fmt"
	"net/mail"
	"net/smtp"
)

// send email
func SendEmail(email string, token string) (string, error) {
	 // Send an email to the user with a link to reset their password
	 resetURL := "http://localhost:8080/login"
	 from := mail.Address{Name: "Agenagn", Address: "selamu.dawit@aastustudent.edu.et"}
	 to := mail.Address{Name: "Selamu Dawit", Address: email}
	 subject := "Password reset request"
	 body := fmt.Sprintf("To reset your password, please follow this link: %s?token=%s", resetURL, token)
 
	 msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s", from.String(), to.String(), subject, body)
	 auth := smtp.PlainAuth("", "selamudev@gmail.com", "sele2inttroduction", "smtp.gmail.com")
	 err1:= smtp.SendMail("smtp.gmail.com:587", auth, from.Address, []string{to.Address}, []byte(msg))
	 if err1 != nil {
		return "", err1
	 }
 
	 return "Email sent", nil
}



