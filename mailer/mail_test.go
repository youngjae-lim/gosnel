package mailer

import (
	"errors"
	"testing"
)

func TestMail_SendSMTPMessage(t *testing.T) {
	msg := Message{
		From:        "me@here.com",
		FromName:    "Joe",
		To:          "you@there.com",
		Subject:     "Test",
		Template:    "test",
		Attachments: []string{"./testdata/mail/test.html.tmpl"},
	}

	err := mailer.SendSMTPMessage(msg)
	if err != nil {
		t.Error(err)
	}
}

func TestMail_SendUsingChan(t *testing.T) {
	msg := Message{
		From:        "me@here.com",
		FromName:    "Joe",
		To:          "you@there.com",
		Subject:     "Test",
		Template:    "test",
		Attachments: []string{"./testdata/mail/test.html.tmpl"},
	}

	mailer.Jobs <- msg
	res := <-mailer.Results
	if res.Error != nil {
		t.Error(errors.New("failed to send over channel"))
	}

	msg.To = "not_an_email_address"
	mailer.Jobs <- msg
	res = <-mailer.Results
	if res.Error == nil {
		t.Error(errors.New("no error received with invalid To address"))
	}
}

func TestMail_SendUsingAPI(t *testing.T) {
	msg := Message{
		To:          "you@there.com",
		Subject:     "Test",
		Template:    "test",
		Attachments: []string{"./testdata/mail/test.html.tmpl"},
	}

	mailer.API = "unknown"
	mailer.APIKey = "abc123"
	mailer.APIUrl = "https://www.fake.com"

	// use 'unknown' instead of SendGrid, MailGun, or SparkPost
	err := mailer.SendUsingAPI(msg, "unknown")
	if err == nil {
		t.Error(errors.New("no error received with unknown input driver"))
	}

	// reset mailer API properties so that it defaults to using SMTP
	mailer.API = ""
	mailer.APIKey = ""
	mailer.APIUrl = ""
}

func TestMail_buildHTMLMessage(t *testing.T) {
	msg := Message{
		From:        "me@here.com",
		FromName:    "Joe",
		To:          "you@there.com",
		Subject:     "Test",
		Template:    "test",
		Attachments: []string{"./testdata/mail/test.html.tmpl"},
	}

	_, err := mailer.buildHTMLMessage(msg)
	if err != nil {
		t.Error(err)
	}
}

func TestMail_buildPlainTextMessage(t *testing.T) {
	msg := Message{
		From:        "me@here.com",
		FromName:    "Joe",
		To:          "you@there.com",
		Subject:     "Test",
		Template:    "test",
		Attachments: []string{"./testdata/mail/test.html.tmpl"},
	}

	_, err := mailer.buildPlainTextMessage(msg)
	if err != nil {
		t.Error(err)
	}
}

func TestMail_Send(t *testing.T) {
	msg := Message{
		From:        "me@here.com",
		FromName:    "Joe",
		To:          "you@there.com",
		Subject:     "Test",
		Template:    "test",
		Attachments: []string{"./testdata/mail/test.html.tmpl"},
	}

	// Send will try to send an email using SMTP
	// because we reset mailer API properties to empty values
	err := mailer.Send(msg)
	if err != nil {
		t.Error(err)
	}

	// now reset mailer API properties to some wrong values
	mailer.API = "unknown"
	mailer.APIKey = "abc123"
	mailer.APIUrl = "https://www.fake.com"

	err = mailer.Send(msg)
	if err == nil {
		t.Error(errors.New("did not get an error when we should have"))
	}

	// reset mailer API properties again
	mailer.API = ""
	mailer.APIKey = ""
	mailer.APIUrl = ""
}

func TestMail_ChooseAPI(t *testing.T) {
	msg := Message{
		From:        "me@here.com",
		FromName:    "Joe",
		To:          "you@there.com",
		Subject:     "Test",
		Template:    "test",
		Attachments: []string{"./testdata/mail/test.html.tmpl"},
	}

	mailer.API = "unknown"
	err := mailer.ChooseAPI(msg)
	if err == nil {
		t.Error(errors.New("no error received with invalid API property"))
	}
}