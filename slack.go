package ecstaskmonitoring

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"io/ioutil"
)

// https://golang.org/doc/faq#guarantee_satisfies_interface
var _ Slack = (*Cluster)(nil)
var _ Slack = (*ParallelNotify)(nil)

// Slack ... Slack operation.
type Slack interface {
	// NewSlackAttachmentMessage ... Generate a message to send to slack
	NewSlackAttachmentMessage(message string) *Attachment
}

// PostSlackMessage ... Verify the revision number and notify the message
func (a *Attachment) PostSlackMessage(incomingWebhook string) {
	s := SlackMessage{
		Attachments: []*Attachment{a},
	}

	msg, _ := json.Marshal(s)

	resp, err := http.PostForm(
		incomingWebhook,
		url.Values{"payload": {string(msg)}},
	)
	if err != nil {
		log.sugar.Warnf("cannot post slack: %v: url: %s", err, incomingWebhook)
		return
	}

	defer resp.Body.Close()
	if _, err := ioutil.ReadAll(resp.Body); err != nil {
		log.sugar.Warnf("cannot post slack: %v", err)
		return
	}
}

// NewSlackAttachmentMessage ... Initialize attachment data of slack for cluster messages
func (c *Cluster) NewSlackAttachmentMessage(message string) *Attachment {
	return &Attachment{
		Color:  ColorRED,
		Title:  "ECS cluster task count threshold has been exceeded",
		Text:   fmt.Sprintf("current: %s > threshold: %d", message, c.TaskThreshold),
		Footer: fmt.Sprintf("%s: %s cluster: %s", c.AwsProfile, c.Name, c.AwsRegion),
	}
}

// NewSlackAttachmentMessage ... Initialize attachment data of slack for Parallel messages
func (p *ParallelNotify) NewSlackAttachmentMessage(_ string) *Attachment {
	return &Attachment{
		Color:  ColorRED,
		Title:  "ECS Task Parallel",
		Text:   fmt.Sprintf("```\n%s```", p.Message),
		Footer: fmt.Sprintf("%s: %s cluster: %s", p.AwsProfile, p.ClusterName, p.AwsRegion),
	}
}
