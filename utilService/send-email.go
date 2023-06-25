package utilService

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
)

type EmailDataToken struct {
    Link string
	Header string
}

type EmailDataMessage struct {
    Message string
	Header string
}


// send email
func SendTokenEmail(email string, token string, header string) (string, error) {
	 // Send an email to the user with a link to reset their password
	 from := "selamu.dawit@aastustudent.edu.et"
	 to := email
	  // Create the reset URL with the token
	 resetURL := "http://localhost:3000/auth/resetPassword/" + token
	 auth := smtp.PlainAuth("", "bb0fbe593f233b", "2f47796776dd86", "sandbox.smtp.mailtrap.io")
	
	t, _ := template.ParseFiles("template.html")

	var body bytes.Buffer   
	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: This is a test subject \n%s\n\n", mimeHeaders)))
	 // Define the email data
	emailData := EmailDataToken{
        Link: resetURL,
		Header: header,
    }
	t.Execute(&body, emailData)
 
	 err := smtp.SendMail("sandbox.smtp.mailtrap.io:2525", auth, from, []string{to}, body.Bytes())
	 return "Email sent", err
}
// send email
func SendMessageEmail(email string, message string, header string) (string, error) {
	// Send an email to the user with a link to reset their password
	from := "selamu.dawit@aastustudent.edu.et"
	to := email
	 // Create the reset URL with the token
	auth := smtp.PlainAuth("", "bb0fbe593f233b", "2f47796776dd86", "sandbox.smtp.mailtrap.io")
   
   t, _ := template.ParseFiles("replyTemplate.html")

   var body bytes.Buffer   
   mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
   body.Write([]byte(fmt.Sprintf("Subject: This is a test subject \n%s\n\n", mimeHeaders)))
   var newHeader = fmt.Sprintf("%s to: %s", header, email)
	// Define the email data
   emailData := EmailDataMessage{
	   Message: message,
	   Header: newHeader,
   }
   t.Execute(&body, emailData)

	err := smtp.SendMail("sandbox.smtp.mailtrap.io:2525", auth, from, []string{to}, body.Bytes())
	return "Email sent", err
}





