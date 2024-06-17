package models

import (
	"125_isbn_new/internal/assert"
	"bytes"
	"crypto/rand"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"log/slog"
	"net"
	"net/smtp"
	"os"
)

var configFile = assert.Path + "config/config.json"

func GenerateConfirmationID() string {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}

func generateMessageID(domain string) string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}

	return fmt.Sprintf("<%s@%s>", base64.StdEncoding.EncodeToString(b), domain)

}

// fetchConfig retrieves the models.MailConfig from config/config.json
func fetchConfig() MailConfig {
	var config MailConfig
	var err error
	data, _ := os.ReadFile(configFile)

	if len(data) == 0 {
		Logger.Error(GetCurrentFuncName() + " No JSON config data found!")
		return config
	}
	err = json.Unmarshal(data, &config)
	if err != nil {
		Logger.Error(GetCurrentFuncName()+" JSON MarshalIndent error!", slog.Any("output", err))
		return config
	}
	return config
}

// SendMail sends a mail to models.TempUser to create his account
func SendMail(temp *TempUser, categorie string) {
	// Fetching mail configuration
	config := fetchConfig()

	// Recipient information
	recipientMail := []string{temp.User.Email}

	// Generating confirmation Id
	temp.ConfirmID = GenerateConfirmationID()

	// Setting the headers
	header := make(map[string]string)
	header["From"] = "Account management" + "<" + config.Email + ">"
	header["To"] = temp.User.Email
	header["Subject"] = "Email verification"
	header["Message-ID"] = generateMessageID(config.Hostname)
	header["Content-Type"] = "text/html; charset=UTF-8"

	t, err := template.ParseFiles(assert.Path + "templates/" + categorie)
	if err != nil {
		log.Fatal(err)
	}

	// Create a buffer to hold the formatted message
	var body bytes.Buffer

	// Execute the mail's template with data
	err = t.Execute(&body, struct {
		Pseudo    string
		ConfirmID string
	}{
		Pseudo:    temp.User.Name,
		ConfirmID: temp.ConfirmID,
	})
	if err != nil {
		log.Fatal(err)
	}

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body.String()

	// Setting the authentication
	auth := smtp.PlainAuth(
		"",
		config.Email,
		config.Auth,
		config.Hostname,
	)

	// Sending the mail using TLS
	err = sendMailTLS(
		fmt.Sprintf("%s:%d", config.Hostname, config.Port),
		auth,
		config.Email,
		recipientMail,
		[]byte(message),
	)

	if err != nil {
		panic(err)
	} else {
		log.Println("Send mail success!")
	}
}

// SendMail sends a mail to models.TempUser to create his account
func SendMailUpdate(user *User) {
	// Fetching mail configuration
	config := fetchConfig()

	// Recipient information
	recipientMail := []string{user.Email}

	// Setting the headers
	header := make(map[string]string)
	header["From"] = "Account management" + "<" + config.Email + ">"
	header["To"] = user.Email
	header["Subject"] = "Email de confirmation"
	header["Message-ID"] = generateMessageID(config.Hostname)
	header["Content-Type"] = "text/html; charset=UTF-8"

	t, err := template.ParseFiles(assert.Path + "templates/mailupdate.gohtml")
	if err != nil {
		log.Fatal(err)
	}

	// Create a buffer to hold the formatted message
	var body bytes.Buffer

	// Execute the mail's template with data
	err = t.Execute(&body, struct {
		Pseudo string
	}{
		Pseudo: user.Name,
	})
	if err != nil {
		log.Fatal(err)
	}

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body.String()

	// Setting the authentication
	auth := smtp.PlainAuth(
		"",
		config.Email,
		config.Auth,
		config.Hostname,
	)

	// Sending the mail using TLS
	err = sendMailTLS(
		fmt.Sprintf("%s:%d", config.Hostname, config.Port),
		auth,
		config.Email,
		recipientMail,
		[]byte(message),
	)

	if err != nil {
		panic(err)
	} else {
		log.Println("Send mail success!")
	}
}

// dial returns a smtp client.
func dial(addr string) (*smtp.Client, error) {
	conn, err := tls.Dial("tcp", addr, nil)
	if err != nil {
		log.Println("Dialing Error:", err)
		return nil, err
	}
	// Explode Host Port String
	host, _, _ := net.SplitHostPort(addr)
	return smtp.NewClient(conn, host)
}

// sendMailTLS: refer to net/smtp func SendMail().
//
// When using net.Dial to connect to the tls (SSL) port, smtp. NewClient() will be stuck and will not prompt err
// When len (to)>1, to [1] starts to prompt that it is secret delivery.
func sendMailTLS(addr string, auth smtp.Auth, from string,
	to []string, msg []byte) (err error) {

	// Create smtp client
	c, err := dial(addr)
	if err != nil {
		log.Println("Create smtp client error:", err)
		return err
	}
	defer c.Close()

	// Checking authentication
	if auth != nil {
		if ok, _ := c.Extension("AUTH"); ok {
			if err = c.Auth(auth); err != nil {
				log.Println("Error during AUTH", err)
				return err
			}
		}
	}

	// Setting recipient
	if err = c.Mail(from); err != nil {
		return err
	}

	for _, addr := range to {
		if err = c.Rcpt(addr); err != nil {
			return err
		}
	}

	// Retrieving the Writer to set the message headers and body
	w, err := c.Data()
	if err != nil {
		return err
	}

	// Writing `msg` in the Writer
	_, err = w.Write(msg)
	if err != nil {
		return err
	}

	// Closing the Writer
	err = w.Close()
	if err != nil {
		return err
	}

	return c.Quit()
}
