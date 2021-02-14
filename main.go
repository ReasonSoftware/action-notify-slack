package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/slack-go/slack"
)

// Config contains parsed user configuration
type Config struct {
	Channel         string
	AttachmentsFile string
	TimestampFile   string
	Timestamp       string
	Fields          []slack.AttachmentField
	Client          *slack.Client
}

// GetConfig returns validated config params
func GetConfig(args []string) (*Config, error) {
	conf := new(Config)

	// arguments
	sep := os.Getenv("SEPARATOR")
	if sep == "" {
		sep = "=="
	}

	fields := make([]slack.AttachmentField, 0)

	for _, arg := range args {
		if len(strings.Split(arg, "\n")) > 0 {
			for _, f := range strings.Split(arg, "\n") {
				if len(strings.Split(f, sep)) == 2 {
					field := slack.AttachmentField{
						Title: strings.TrimSpace(strings.Split(f, sep)[0]),
						Value: strings.TrimSpace(strings.Split(f, sep)[1]),
						Short: true,
					}

					fields = append(fields, field)
				}
			}
		}
	}

	// env.vars
	timestampFile := os.Getenv("TIMESTAMP_FILE")
	var timestamp string
	if timestampFile != "" {
		err := os.MkdirAll(filepath.Dir(timestampFile), os.ModePerm)
		if err != nil {
			return conf, errors.New("error creating timestamp file")
		}

		t, err := ioutil.ReadFile(os.Getenv("TIMESTAMP_FILE"))
		if err != nil && strings.Contains(err.Error(), "no such file or directory") {
			timestamp = ""
		} else if err != nil {
			return conf, errors.New("error reading timestamp file")
		} else {
			timestamp = string(t)
		}
	} else {
		timestamp = os.Getenv("TIMESTAMP")
	}

	channel := os.Getenv("CHANNEL")
	if channel == "" {
		return conf, errors.New("missing Slack channel")
	}

	attachmentsFile := os.Getenv("ATTACHMENTS_FILE")

	t := os.Getenv("TOKEN")
	if t == "" {
		return conf, errors.New("missing Slack token")
	}

	conf.Channel = channel
	conf.AttachmentsFile = attachmentsFile
	conf.TimestampFile = timestampFile
	conf.Timestamp = timestamp
	conf.Fields = fields
	conf.Client = slack.New(t)

	return conf, nil
}

func main() {
	vars := []string{
		"GITHUB_ACTOR",
		"GITHUB_REPOSITORY",
		"STATUS",
		"GITHUB_WORKFLOW",
	}

	for _, v := range vars {
		if os.Getenv(v) == "" {
			fmt.Printf("missing required env.var: '%s'", v)
			os.Exit(1)
		}
	}

	conf, err := GetConfig(os.Args[1:])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	s := Slack{
		Channel:   conf.Channel,
		Context:   ctx,
		Timestamp: conf.Timestamp,
	}

	var ts string

	if conf.AttachmentsFile == "" {
		ts, err = s.SendTemplate(conf.Client, conf.Fields)
	} else {
		ts, err = s.SendAttachmentFromFile(conf.Client, conf.AttachmentsFile, conf.Fields)
	}

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if conf.TimestampFile != "" {
		err = ioutil.WriteFile(conf.TimestampFile, []byte(ts), 0644)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	fmt.Printf("::set-output name=TIMESTAMP::%s\n", ts)
}
