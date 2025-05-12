package mail

import (
	"testing"

	"github.com/guncv/Simple-Bank/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSendEmailWithGmail(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping email sending for testing")
	}

	config, err := util.LoadConfig("..")
	require.NoError(t, err)

	sender := NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)

	subject := "A test email"
	content := `
		<h1>Hello world</h1>
		<p>This is a test email</p>
	`
	to := []string{"6431309721@student.chula.ac.th"}

	attachFiles := []string{"../README.md"}
	err = sender.SendEmail(subject, content, to, nil, nil, attachFiles)
	assert.NoError(t, err)
}
