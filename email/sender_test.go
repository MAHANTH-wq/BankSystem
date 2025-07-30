package email

import (
	"testing"

	"github.com/mahanth/simplebank/util"
	"github.com/stretchr/testify/require"
)

func TestSendEMailWithGmail(t *testing.T) {

	if testing.Short(){
		t.Skip()
	}
	config, err := util.LoadConfig("..")
	require.NoError(t, err)
	//Added comments to explain the code
	sender := NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)

	subject := "Test Email From Bank System Project"
	content := `
	<h1>Hello World</h1>
	<p>This is a test email sent from the Bank System project.</p>
	`
	to := []string{"vallulrimahanthkumar@gmail.com"}
	attachFiles := []string{"../README.md"}

	err = sender.SendEmail(subject, content, to, nil, nil, attachFiles)
	require.NoError(t, err, "Failed to send email")
	t.Log("Email sent successfully")
}
