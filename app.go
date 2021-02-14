package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/slack-go/slack"
)

// Client represents Slack client
type Client interface {
	PostMessageContext(context.Context, string, ...slack.MsgOption) (string, string, error)
	UpdateMessageContext(ctx context.Context, channelID, timestamp string, options ...slack.MsgOption) (string, string, string, error)
}

// Slack represents app config
type Slack struct {
	Channel   string
	Context   context.Context
	Timestamp string
}

// GetTemplate returns default Slack Message Template
func GetTemplate(status string, failed bool, additions []slack.AttachmentField) slack.Attachment {
	var color string
	var failureColor string = "#fd0000"

	if strings.ToLower(status) == "running" || strings.ToLower(status) == "started" || strings.ToLower(status) == "building" || strings.ToLower(status) == "initializing" {
		color = "#fbf000"
	} else if strings.ToLower(status) == "deploying" || strings.ToLower(status) == "uploading" || strings.ToLower(status) == "publishing" || strings.ToLower(status) == "creating" {
		color = "#fda100"
	} else if strings.ToLower(status) == "finished" || strings.ToLower(status) == "succeeded" || strings.ToLower(status) == "passed" || strings.ToLower(status) == "built" || strings.ToLower(status) == "released" {
		color = "#0ce823"
	} else if strings.ToLower(status) == "failed" || strings.ToLower(status) == "aborted" || strings.ToLower(status) == "canceled" || strings.ToLower(status) == "terminated" {
		color = failureColor
	} else {
		color = "#777777"
	}

	s := status
	if failed {
		color = failureColor
		s = "failed"
	}

	fields := []slack.AttachmentField{
		{
			Title: "Repository",
			Value: fmt.Sprintf("<https://github.com/%s|%s>", os.Getenv("GITHUB_REPOSITORY"), strings.Split(os.Getenv("GITHUB_REPOSITORY"), "/")[1]),
			Short: true,
		},
		{
			Title: "Workflow",
			Value: fmt.Sprintf("<https://github.com/%s/actions?query=workflow", os.Getenv("GITHUB_REPOSITORY")) + "%3A" + fmt.Sprintf("%s|%s>", os.Getenv("GITHUB_WORKFLOW"), os.Getenv("GITHUB_WORKFLOW")),
			Short: true,
		},
		{
			Title: "Initiator",
			Value: fmt.Sprintf("<https://github.com/%s|%s>", os.Getenv("GITHUB_ACTOR"), os.Getenv("GITHUB_ACTOR")),
			Short: true,
		},
		{
			Title: "Status",
			Value: fmt.Sprintf("<https://github.com/%s/actions/runs/%s|%s>", os.Getenv("GITHUB_REPOSITORY"), os.Getenv("GITHUB_RUN_ID"), strings.ToUpper(s)),
			Short: true,
		},
	}

	if len(additions) > 0 {
		fields = append(fields, additions...)
	}

	msg := slack.Attachment{
		Color:      color,
		Fields:     fields,
		Footer:     "<https://github.com/ReasonSoftware/action-notify-slack|ReasonSoftware/action-notify-slack>",
		FooterIcon: "https://cdn.reasonsecurity.com/images/logo.png",
		Ts:         json.Number(strconv.FormatInt(time.Now().Unix(), 10)),
	}

	return msg
}

// SendTemplate sends a template message
func (s *Slack) SendTemplate(cli Client, fields []slack.AttachmentField) (string, error) {
	failure := false

	if os.Getenv("FAIL") != "" {
		var err error
		failure, err = strconv.ParseBool(os.Getenv("FAIL"))
		if err != nil {
			return "", errors.Wrap(err, "error parsing env.var 'FAIL'")
		}
	}

	t := GetTemplate(os.Getenv("STATUS"), failure, fields)
	return s.send(cli, slack.MsgOptionAttachments(t))
}

// SendAttachmentFromFile sends an attachment provided via json file
func (s *Slack) SendAttachmentFromFile(cli Client, filename string, fields []slack.AttachmentField) (string, error) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("error reading file '%s'", filename))
	}

	var slice []slack.Attachment
	if err := json.Unmarshal(file, &slice); err != nil {
		var single slack.Attachment

		if err := json.Unmarshal(file, &single); err != nil {
			return "", errors.Wrap(err, fmt.Sprintf("invalid JSON file '%s'", filename))
		}

		single.Fields = append(single.Fields, fields...)
		return s.send(cli, slack.MsgOptionAttachments(single))
	}

	for _, a := range slice {
		a.Fields = append(a.Fields, fields...)
	}

	return s.send(cli, slack.MsgOptionAttachments(slice...))
}

func (s *Slack) send(cli Client, options ...slack.MsgOption) (string, error) {
	if s.Timestamp != "" {
		_, ts, _, err := cli.UpdateMessageContext(s.Context, s.Channel, s.Timestamp, options...)
		if err != nil {
			return "", errors.Wrap(err, "error updating message")
		}

		return ts, nil
	}

	_, ts, err := cli.PostMessageContext(s.Context, s.Channel, options...)
	if err != nil {
		return "", errors.Wrap(err, "error sending message")
	}

	return ts, nil
}
