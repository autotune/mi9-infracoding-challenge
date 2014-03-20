// +build !appengine

package main

import (
	"net/smtp"
	"os"
	"strings"
)

var smtpHost string
var smtpPort string
var smtpUsername string
var smtpPassword string
var smtpRecipients []string
var smtpFailRecipients []string

func init() {
	smtpHost = os.Getenv("SMTP_HOST")
	smtpPort = os.Getenv("SMTP_PORT")
	smtpUsername = os.Getenv("SMTP_USERNAME")
	smtpPassword = os.Getenv("SMTP_PASSWORD")
	smtpRecipients = strings.Split(os.Getenv("SMTP_RECIPIENTS"), ",")
	smtpFailRecipients = strings.Split(os.Getenv("SMTP_FAIL_RECIPIENTS"), ",")
	if smtpPort == "" {
		smtpPort = "25"
	}
}

func SendConfirmation(candidate, email, url string) error {
	println(candidate + " (" + email + ") just passed with URL " + url)
	if smtpHost == "" {
		return nil
	}
	message := "SUBJECT: " + candidate + " just passed the coding test!\n\nSee the submission at " + url + "\n\nLet's get in touch with them at " + email
	return smtp.SendMail(smtpHost+":"+smtpPort, smtp.PlainAuth("", smtpUsername, smtpPassword, smtpHost),
		"codingchallenge-noreply@mi9.com.au", smtpRecipients, []byte(message))
}

func SendFailure(candidate, email, url string) error {
	if smtpHost == "" {
		return nil
	}
	message := "SUBJECT: " + candidate + " just attempted the coding test and failed!\n\nSee the submission at " + url
	return smtp.SendMail(smtpHost+":"+smtpPort, smtp.PlainAuth("", smtpUsername, smtpPassword, smtpHost),
		"codingchallenge-noreply@mi9.com.au", smtpFailRecipients, []byte(message))
}
