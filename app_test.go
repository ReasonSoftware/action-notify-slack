package main_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	app "action-notify-slack"

	"action-notify-slack/mocks"

	"github.com/pkg/errors"
	"github.com/slack-go/slack"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func init() {
	vars := []string{
		"GITHUB_ACTOR",
		"GITHUB_REPOSITORY",
		"STATUS",
		"GITHUB_WORKFLOW",
	}

	for _, v := range vars {
		if os.Getenv(v) == "" {
			os.Setenv(v, "unit-test")
		}
	}
}

func TestSendAttachmentFromFile(t *testing.T) {
	assert := assert.New(t)

	filename, err := ioutil.TempFile(os.TempDir(), "test-")
	assert.Equal(nil, err, "preparation: error creating temporary file")
	defer os.Remove(filename.Name())

	type test struct {
		Receiver       *app.Slack
		Parameter1     string
		Parameter2     []slack.AttachmentField
		Attachment     []byte
		ExpectedOutput string
		MockError      error
		ExpectedError  string
	}

	suite := map[string]test{
		"Single Attachment": {
			Receiver: &app.Slack{
				Channel:   "self",
				Context:   context.Background(),
				Timestamp: "",
			},
			Parameter1: filename.Name(),
			Parameter2: []slack.AttachmentField{},
			Attachment: []byte(`{
	"color": "#FBF000",
	"fields": [
		{
			"title": "Repository",
			"value": "<https://github.com/ReasonSoftware/project|project>",
			"short": true
		},
		{
			"title": "Workflow",
			"value": "<https://github.com/ReasonSoftware/project/actions?query=workflow%3Aproduction|production>",
			"short": true
		},
		{
			"title": "Initiator",
			"value": "<https://github.com/anton-yurchenko|anton-yurchenko>",
			"short": true
		},
		{
			"title": "Status",
			"value": "STARTED",
			"short": true
		}
	],
	"footer": "<https://github.com/ReasonSoftware/action-notify-slack|ReasonSoftware/action-notify-slack>",
	"footer_icon": "https://cdn.reasonsecurity.com/images/logo.png",
	"ts": 1588885167
}`),
			ExpectedOutput: fmt.Sprint(time.Now().Unix()),
			MockError:      nil,
			ExpectedError:  "",
		},
		"slack.PostMessageContext Error": {
			Receiver: &app.Slack{
				Channel:   "self",
				Context:   context.Background(),
				Timestamp: "",
			},
			Parameter1: filename.Name(),
			Parameter2: []slack.AttachmentField{},
			Attachment: []byte(`{
			"color": "#FBF000",
			"fields": [
				{
					"title": "Repository",
					"value": "<https://github.com/ReasonSoftware/project|project>",
					"short": true
				},
				{
					"title": "Workflow",
					"value": "<https://github.com/ReasonSoftware/project/actions?query=workflow%3Aproduction|production>",
					"short": true
				},
				{
					"title": "Initiator",
					"value": "<https://github.com/anton-yurchenko|anton-yurchenko>",
					"short": true
				},
				{
					"title": "Status",
					"value": "STARTED",
					"short": true
				}
			],
			"footer": "<https://github.com/ReasonSoftware/action-notify-slack|ReasonSoftware/action-notify-slack>",
			"footer_icon": "https://cdn.reasonsecurity.com/images/logo.png",
			"ts": 1588885167
		}`),
			ExpectedOutput: "",
			MockError:      errors.New("reason"),
			ExpectedError:  "error sending message: reason",
		},
		"Update": {
			Receiver: &app.Slack{
				Channel:   "self",
				Context:   context.Background(),
				Timestamp: fmt.Sprint(time.Now().Unix()),
			},
			Parameter1: filename.Name(),
			Parameter2: []slack.AttachmentField{},
			Attachment: []byte(`{
			"color": "#0CE823",
			"fields": [
				{
					"title": "Repository",
					"value": "<https://github.com/ReasonSoftware/project|project>",
					"short": true
				},
				{
					"title": "Workflow",
					"value": "<https://github.com/ReasonSoftware/project/actions?query=workflow%3Aproduction|production>",
					"short": true
				},
				{
					"title": "Initiator",
					"value": "<https://github.com/anton-yurchenko|anton-yurchenko>",
					"short": true
				},
				{
					"title": "Status",
					"value": "FINISHED",
					"short": true
				}
			],
			"footer": "<https://github.com/ReasonSoftware/action-notify-slack|ReasonSoftware/action-notify-slack>",
			"footer_icon": "https://cdn.reasonsecurity.com/images/logo.png",
			"ts": 1588885167
		}`),
			ExpectedOutput: fmt.Sprint(time.Now().Unix()),
			MockError:      nil,
			ExpectedError:  "",
		},
		"slack.UpdateMessageContext Error": {
			Receiver: &app.Slack{
				Channel:   "self",
				Context:   context.Background(),
				Timestamp: fmt.Sprint(time.Now().Unix()),
			},
			Parameter1: filename.Name(),
			Parameter2: []slack.AttachmentField{},
			Attachment: []byte(`{
			"color": "#0CE823",
			"fields": [
				{
					"title": "Repository",
					"value": "<https://github.com/ReasonSoftware/project|project>",
					"short": true
				},
				{
					"title": "Workflow",
					"value": "<https://github.com/ReasonSoftware/project/actions?query=workflow%3Aproduction|production>",
					"short": true
				},
				{
					"title": "Initiator",
					"value": "<https://github.com/anton-yurchenko|anton-yurchenko>",
					"short": true
				},
				{
					"title": "Status",
					"value": "FINISHED",
					"short": true
				}
			],
			"footer": "<https://github.com/ReasonSoftware/action-notify-slack|ReasonSoftware/action-notify-slack>",
			"footer_icon": "https://cdn.reasonsecurity.com/images/logo.png",
			"ts": 1588885167
		}`),
			ExpectedOutput: "",
			MockError:      errors.New("reason"),
			ExpectedError:  "error updating message: reason",
		},
		"Multiple Attachments": {
			Receiver: &app.Slack{
				Channel:   "self",
				Context:   context.Background(),
				Timestamp: "",
			},
			Parameter1: filename.Name(),
			Parameter2: []slack.AttachmentField{},
			Attachment: []byte(`[
			{
				"color": "#FBF000",
				"fields": [
					{
						"title": "Repository",
						"value": "<https://github.com/ReasonSoftware/project-1|project-1>",
						"short": true
					},
					{
						"title": "Workflow",
						"value": "<https://github.com/ReasonSoftware/project-1/actions?query=workflow%3Aproduction|production>",
						"short": true
					},
					{
						"title": "Initiator",
						"value": "<https://github.com/anton-yurchenko|anton-yurchenko>",
						"short": true
					},
					{
						"title": "Status",
						"value": "STARTED",
						"short": true
					}
				],
				"footer": "<https://github.com/ReasonSoftware/action-notify-slack|ReasonSoftware/action-notify-slack>",
				"footer_icon": "https://cdn.reasonsecurity.com/images/logo.png",
				"ts": 1588885167
			},
			{
				"color": "#FBF000",
				"fields": [
					{
						"title": "Repository",
						"value": "<https://github.com/ReasonSoftware/project-2|project-2>",
						"short": true
					},
					{
						"title": "Workflow",
						"value": "<https://github.com/ReasonSoftware/project-2/actions?query=workflow%3Aproduction|production>",
						"short": true
					},
					{
						"title": "Initiator",
						"value": "<https://github.com/anton-yurchenko|anton-yurchenko>",
						"short": true
					},
					{
						"title": "Status",
						"value": "STARTED",
						"short": true
					}
				],
				"footer": "<https://github.com/ReasonSoftware/action-notify-slack|ReasonSoftware/action-notify-slack>",
				"footer_icon": "https://cdn.reasonsecurity.com/images/logo.png",
				"ts": 1588885167
			}
		]`),
			ExpectedOutput: fmt.Sprint(time.Now().Unix()),
			MockError:      nil,
			ExpectedError:  "",
		},
		"Invalid JSON": {
			Receiver: &app.Slack{
				Channel:   "self",
				Context:   context.Background(),
				Timestamp: fmt.Sprint(time.Now().Unix()),
			},
			Parameter1: filename.Name(),
			Parameter2: []slack.AttachmentField{},
			Attachment: []byte(`{
			"color": "#0CE823",
			"ts": 1588885167,
		}`),
			ExpectedOutput: "",
			MockError:      nil,
			ExpectedError:  fmt.Sprintf("invalid JSON file '%s': invalid character '}' looking for beginning of object key string", filename.Name()),
		},
	}

	var counter int
	for name, test := range suite {
		counter++
		t.Logf("Test Case %v/%v - %s", counter, len(suite), name)

		file, err := os.OpenFile(filename.Name(), os.O_CREATE, os.ModePerm)
		assert.Equal(nil, err, "preparation: error opening a temporary file")
		defer file.Close()

		err = ioutil.WriteFile(filename.Name(), test.Attachment, 0644)
		assert.Equal(nil, err, "preparation: error writing json to temporary file")

		m := new(mocks.Client)

		m.On("UpdateMessageContext", test.Receiver.Context, test.Receiver.Channel, test.Receiver.Timestamp, mock.AnythingOfType("slack.MsgOption")).Return("", test.ExpectedOutput, "", test.MockError)
		m.On("PostMessageContext", test.Receiver.Context, test.Receiver.Channel, mock.AnythingOfType("slack.MsgOption")).Return("", test.ExpectedOutput, test.MockError)

		result, err := test.Receiver.SendAttachmentFromFile(m, test.Parameter1, test.Parameter2)

		if test.ExpectedError != "" {
			assert.EqualError(err, test.ExpectedError)
		} else {
			assert.Equal(nil, err)
		}

		assert.Equal(test.ExpectedOutput, result)
	}
}

func TestSendTemplate(t *testing.T) {
	assert := assert.New(t)

	type test struct {
		Receiver       *app.Slack
		Parameter1     []slack.AttachmentField
		ExpectedOutput string
		MockError      error
		ExpectedError  string
	}

	suite := map[string]test{
		"Template": {
			Receiver: &app.Slack{
				Channel:   "self",
				Context:   context.Background(),
				Timestamp: "",
			},
			Parameter1:     []slack.AttachmentField{},
			ExpectedOutput: fmt.Sprint(time.Now().Unix()),
			MockError:      nil,
			ExpectedError:  "",
		},
		"Update": {
			Receiver: &app.Slack{
				Channel:   "self",
				Context:   context.Background(),
				Timestamp: fmt.Sprint(time.Now().Unix()),
			},
			Parameter1:     []slack.AttachmentField{},
			ExpectedOutput: fmt.Sprint(time.Now().Unix()),
			MockError:      nil,
			ExpectedError:  "",
		},
		"Template with Fields": {
			Receiver: &app.Slack{
				Channel:   "self",
				Context:   context.Background(),
				Timestamp: "",
			},
			Parameter1: []slack.AttachmentField{
				{
					Title: "key-1",
					Value: "value-1",
					Short: true,
				},
				{
					Title: "key-2",
					Value: "value-2",
					Short: true,
				},
			},
			ExpectedOutput: fmt.Sprint(time.Now().Unix()),
			MockError:      nil,
			ExpectedError:  "",
		},
		"Update with Fields": {
			Receiver: &app.Slack{
				Channel:   "self",
				Context:   context.Background(),
				Timestamp: fmt.Sprint(time.Now().Unix()),
			},
			Parameter1: []slack.AttachmentField{
				{
					Title: "key-1",
					Value: "value-1",
					Short: true,
				},
				{
					Title: "key-2",
					Value: "value-2",
					Short: true,
				},
			},
			ExpectedOutput: fmt.Sprint(time.Now().Unix()),
			MockError:      nil,
			ExpectedError:  "",
		},
		"slack.PostMessageContext Error": {
			Receiver: &app.Slack{
				Channel:   "self",
				Context:   context.Background(),
				Timestamp: "",
			},
			Parameter1:     []slack.AttachmentField{},
			ExpectedOutput: "",
			MockError:      errors.New("reason"),
			ExpectedError:  "error sending message: reason",
		},
		"slack.UpdateMessageContext Error": {
			Receiver: &app.Slack{
				Channel:   "self",
				Context:   context.Background(),
				Timestamp: fmt.Sprint(time.Now().Unix()),
			},
			Parameter1:     []slack.AttachmentField{},
			ExpectedOutput: "",
			MockError:      errors.New("reason"),
			ExpectedError:  "error updating message: reason",
		},
	}

	var counter int
	for name, test := range suite {
		counter++
		t.Logf("Test Case %v/%v - %s", counter, len(suite), name)

		m := new(mocks.Client)

		m.On("UpdateMessageContext", test.Receiver.Context, test.Receiver.Channel, test.Receiver.Timestamp, mock.AnythingOfType("slack.MsgOption")).Return("", test.ExpectedOutput, "", test.MockError)
		m.On("PostMessageContext", test.Receiver.Context, test.Receiver.Channel, mock.AnythingOfType("slack.MsgOption")).Return("", test.ExpectedOutput, test.MockError)

		result, err := test.Receiver.SendTemplate(m, test.Parameter1)

		if test.ExpectedError != "" {
			assert.EqualError(err, test.ExpectedError)
		} else {
			assert.Equal(nil, err)
		}
		assert.Equal(test.ExpectedOutput, result)
	}
}

func TestMain(m *testing.M) {
	os.Setenv("GITHUB_ACTOR", "username")
	os.Setenv("GITHUB_REPOSITORY", "ore/proj")
	os.Setenv("STATUS", "running")
	os.Setenv("GITHUB_WORKFLOW", "testing")

	code := m.Run()

	os.Unsetenv("GITHUB_ACTOR")
	os.Unsetenv("GITHUB_REPOSITORY")
	os.Unsetenv("STATUS")
	os.Unsetenv("GITHUB_WORKFLOW")

	os.Exit(code)
}
