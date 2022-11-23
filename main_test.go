package main_test

import (
	"fmt"
	"os"
	"testing"

	app "action-notify-slack"

	"github.com/slack-go/slack"
	"github.com/stretchr/testify/assert"
)

func TestGetConfig(t *testing.T) {
	assert := assert.New(t)

	type test struct {
		Channel         string
		AttachmentsFile string
		Token           string
		TimestampFile   bool
		Timestamp       string
		Arguments       []string
		ExpectedFields  []slack.AttachmentField
		ExpectedError   string
	}

	suite := map[string]test{
		"Template": {
			Channel:         "self",
			AttachmentsFile: "",
			Token:           "secret-text",
			TimestampFile:   false,
			Timestamp:       "",
			Arguments:       []string{},
			ExpectedFields:  []slack.AttachmentField{},
			ExpectedError:   "",
		},
		"Timestamp File": {
			Channel:         "self",
			AttachmentsFile: "",
			Token:           "secret-text",
			TimestampFile:   true,
			Timestamp:       "1589146397.007200",
			Arguments:       []string{},
			ExpectedFields:  []slack.AttachmentField{},
			ExpectedError:   "",
		},
		"Template with Fields": {
			Channel:         "self",
			AttachmentsFile: "",
			Token:           "secret-text",
			TimestampFile:   false,
			Timestamp:       "",
			Arguments: []string{`field1==value1
		field 2==value 2`},
			ExpectedFields: []slack.AttachmentField{
				{
					Title: "field1",
					Value: "value1",
					Short: true,
				},
				{
					Title: "field 2",
					Value: "value 2",
					Short: true,
				},
			},
			ExpectedError: "",
		},
		"Custom Attachments File": {
			Channel:         "self",
			AttachmentsFile: "attachments.json",
			Token:           "secret-text",
			TimestampFile:   false,
			Timestamp:       "",
			Arguments:       []string{},
			ExpectedFields:  []slack.AttachmentField{},
		},
		"Custom Attachments File with Fields": {
			Channel:         "self",
			AttachmentsFile: "attachments.json",
			Token:           "secret-text",
			TimestampFile:   false,
			Timestamp:       "",
			Arguments: []string{`field1==value1
		field 2==value 2`},
			ExpectedFields: []slack.AttachmentField{
				{
					Title: "field1",
					Value: "value1",
					Short: true,
				},
				{
					Title: "field 2",
					Value: "value 2",
					Short: true,
				},
			},
			ExpectedError: "",
		},
		"Missing Channel": {
			Channel:         "",
			AttachmentsFile: "attachments.json",
			Token:           "secret-text",
			TimestampFile:   false,
			Timestamp:       "",
			Arguments:       []string{},
			ExpectedFields:  []slack.AttachmentField{},
			ExpectedError:   "missing Slack channel",
		},
		"Missing Token": {
			Channel:         "self",
			AttachmentsFile: "",
			Token:           "",
			TimestampFile:   false,
			Timestamp:       "",
			Arguments:       []string{},
			ExpectedFields:  []slack.AttachmentField{},
			ExpectedError:   "missing Slack token",
		},
		"Arguments": {
			Channel:         "self",
			AttachmentsFile: "",
			Token:           "secret-text",
			TimestampFile:   false,
			Timestamp:       "",
			Arguments: []string{`field1==value1
		field 2==value 2`},
			ExpectedFields: []slack.AttachmentField{
				{
					Title: "field1",
					Value: "value1",
					Short: true,
				},
				{
					Title: "field 2",
					Value: "value 2",
					Short: true,
				},
			},
			ExpectedError: "",
		},
	}

	var counter int
	for name, test := range suite {
		counter++
		t.Logf("Test Case %v/%v - %s", counter, len(suite), name)

		err := os.Setenv("CHANNEL", test.Channel)
		assert.Equal(nil, err, "preparation: error setting env.var 'CHANNEL'")
		defer os.Unsetenv("CHANNEL")

		err = os.Setenv("ATTACHMENTS_FILE", test.AttachmentsFile)
		assert.Equal(nil, err, "preparation: error setting env.var 'ATTACHMENTS_FILE'")
		defer os.Unsetenv("ATTACHMENTS_FILE")

		err = os.Setenv("TOKEN", test.Token)
		assert.Equal(nil, err, "preparation: error setting env.var 'TOKEN'")
		defer os.Unsetenv("TOKEN")

		var file string
		if test.TimestampFile {
			dir, err := os.MkdirTemp(".", "unittests")
			if err != nil {
				assert.Equal(nil, err, "preparation: error creating temporary directory")
			}
			defer os.RemoveAll(dir)

			file = fmt.Sprintf("%s/file", dir)

			err = os.WriteFile(file, []byte(test.Timestamp), 0644)
			if err != nil {
				assert.Equal(nil, err, "preparation: error writing to timestamp file")
			}

			err = os.Setenv("TIMESTAMP_FILE", file)
			assert.Equal(nil, err, "preparation: error setting env.var 'TIMESTAMP_FILE'")
			defer os.Unsetenv("TIMESTAMP_FILE")
		}

		conf, err := app.GetConfig(test.Arguments)

		if test.ExpectedError != "" {
			assert.EqualError(err, test.ExpectedError)
		} else {
			assert.Equal(nil, err)

			c := app.Config{
				Channel:         test.Channel,
				AttachmentsFile: test.AttachmentsFile,
				TimestampFile:   file,
				Timestamp:       test.Timestamp,
				Fields:          test.ExpectedFields,
				Client:          slack.New(test.Token),
			}

			assert.Equal(c, *conf)
		}

		os.Unsetenv("TIMESTAMP_FILE")
	}
}
