package smtp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"golang.org/x/oauth2/clientcredentials"
	"golang.org/x/oauth2/microsoft"
)

// CONFIGURATION
const (
	TenantID     = "908g90dfgkj-5172-4bfe-a987-gfndgioret8" // From Azure Overview
	ClientID     = "3a5229e4-dfg-fgddfg-fdg497a531267777"   // From Azure Overview
	ClientSecret = "YCs8Q~ikgjfdiodfjgjgkl-gdfdffgdfg4554"  // From Certificates & Secrets
	SenderEmail  = "vexora-noreply@codershubinc.tech"       // Must match your verified domain
)

// Email Structure for Microsoft Graph
type EmailMessage struct {
	Message struct {
		Subject string `json:"subject"`
		Body    struct {
			ContentType string `json:"contentType"`
			Content     string `json:"content"`
		} `json:"body"`
		ToRecipients []Recipient `json:"toRecipients"`
	} `json:"message"`
	SaveToSentItems bool `json:"saveToSentItems"`
}

type Recipient struct {
	EmailAddress struct {
		Address string `json:"address"`
	} `json:"emailAddress"`
}

func Mail() {
	// 1. Configure the OAuth2 Client Credentials flow
	conf := &clientcredentials.Config{
		ClientID:     ClientID,
		ClientSecret: ClientSecret,
		TokenURL:     microsoft.AzureADEndpoint(TenantID).TokenURL,
		Scopes:       []string{"https://graph.microsoft.com/.default"},
	}

	// 2. Get an authenticated HTTP client
	client := conf.Client(context.Background())
	// 3. Prepare the email payload
	emailReq := EmailMessage{
		SaveToSentItems: false, // Set to true if you want a copy in the "Sent" folder
	}

	// Define the HTML content using backticks
	htmlBody := `
	<!DOCTYPE html>
	<html>
	<head>
		<meta charset="utf-8">
		<style>
			body { font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; background-color: #f4f4f4; margin: 0; padding: 0; }
			.container { max-width: 600px; margin: 20px auto; background-color: #ffffff; border-radius: 8px; overflow: hidden; box-shadow: 0 4px 6px rgba(0,0,0,0.1); }
			.header { background-color: #0f172a; padding: 30px; text-align: center; }
			.header h1 { color: #ffffff; margin: 0; font-size: 24px; letter-spacing: 1px; }
			.content { padding: 30px; color: #334155; line-height: 1.6; }
			.button { display: inline-block; padding: 12px 24px; background-color: #2563eb; color: #ffffff; text-decoration: none; border-radius: 6px; font-weight: bold; margin-top: 20px; }
			.footer { background-color: #f1f5f9; padding: 20px; text-align: center; font-size: 12px; color: #64748b; }
		</style>
	</head>
	<body>
		<div class="container">
			<div class="header">
				<h1>VEXORA STUDIO</h1>
			</div>
			<div class="content">
				<h2>Welcome, Developer! 🚀</h2>
				<p>We are thrilled to have you on board. Vexora Studio is ready to help you build amazing things.</p>
				<p>This is a test email sent directly from your <strong>Go backend</strong> using the Microsoft Graph API.</p>
				<center>
					<a href="https://codershubinc.tech" class="button">Visit Dashboard</a>
				</center>
			</div>
			<div class="footer">
				<p>&copy; 2025 CodersHub Inc. All rights reserved.</p>
				<p>Sent via Microsoft Graph API</p>
			</div>
		</div>
	</body>
	</html>
	`

	emailReq.Message.Subject = "🚀 Welcome to Vexora Studio!"
	emailReq.Message.Body.ContentType = "HTML"
	emailReq.Message.Body.Content = htmlBody

	emailReq.Message.ToRecipients = []Recipient{
		{EmailAddress: struct {
			Address string `json:"address"`
		}{Address: "your-email@example.com"}},
	}
	// Add Recipient
	emailReq.Message.ToRecipients = []Recipient{
		{EmailAddress: struct {
			Address string `json:"address"`
		}{Address: "ingleswapnil2004@gmail.com"}},
	}

	jsonData, _ := json.Marshal(emailReq)

	// 4. Make the POST request to Graph API
	// Note: We use the /users/{id}/sendMail endpoint
	url := fmt.Sprintf("https://graph.microsoft.com/v1.0/users/%s/sendMail", SenderEmail)

	resp, err := client.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// 5. Check Response
	if resp.StatusCode == http.StatusAccepted {
		fmt.Println("Success! Email queued for sending (202 Accepted).")
	} else {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("Error sending email: %s\nResponse: %s\n", resp.Status, string(body))
	}
}
