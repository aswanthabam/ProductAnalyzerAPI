package mailer

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"productanalyzer/api/config"
	api_error "productanalyzer/api/errors"
)

type HTMLEmailRequest struct {
	From    string   `json:"from"`
	To      string   `json:"to"`
	Subject string   `json:"subject"`
	HTML    string   `json:"html"`
	BCC     []string `json:"bcc"`
}

func SendHTMLEmail(email, subject, mailContent string) *api_error.APIError {
	url := "https://api.resend.com/emails"

	emailReq := HTMLEmailRequest{
		From:    "Product Analyzer <" + config.Config.MAIL_ID + ">",
		To:      email,
		Subject: subject,
		HTML:    mailContent,
		BCC:     []string{},
	}

	jsonData, err := json.Marshal(emailReq)
	if err != nil {
		return api_error.UnexpectedError(err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return api_error.UnexpectedError(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+config.Config.RESEND_API_KEY)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return api_error.UnexpectedError(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return api_error.UnexpectedError(err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Println(string(body))
		return api_error.UnexpectedError(nil)
	}
	return nil
}
