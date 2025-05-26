package mail

import (
	"fmt"
	"os"
	"simple-bank/util"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSenderEmail(t *testing.T) {
	dir, _ := os.Getwd()
	fmt.Println("Running test from:", dir)

	config, err := util.LoadConfig("..")
	fmt.Printf("error load config : %s", err)
	require.NoError(t, err)

	sender := NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)

	subject := "A Test Email"
	content := `
	<h1>Hello World</h1>
	<p> This is a test message from <a href="https://cdfas.fun">Aditya Putra</a></p>
	`
	to := []string{"aditiyaputra42@gmail.com", "pratamaadityaputra777@gmail.com"}
	attachFile := []string{"../README.md"}

	err = sender.SendEmail(subject, content, to, nil, nil, attachFile)
	fmt.Printf("error : %s", err)
	require.NoError(t, err)

}
