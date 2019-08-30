package ecstaskmonitoring

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"strings"

	"github.com/aws/aws-sdk-go/aws"
)

// https://golang.org/doc/faq#guarantee_satisfies_interface
var _ Slack = (*Cluster)(nil)
var _ Slack = (*Task)(nil)

// Slack ... Slack operation.
type Slack interface {
	// NewSlackAttachmentMessage ... Generate a message to send to slack
	NewSlackAttachmentMessage(message, awsProfile string) *Attachment
}

// PostSlackMessage ... Verify the revision number and notify the message
func (a *Attachment) PostSlackMessage(incomingWebhook string) error {
	s := SlackMessage{
		Attachments: []*Attachment{a},
	}

	msg, _ := json.Marshal(s)

	resp, _ := http.PostForm(
		incomingWebhook,
		url.Values{"payload": {string(msg)}},
	)

	retry := 3
	for retry > 0 {
		if _, err := ioutil.ReadAll(resp.Body); err != nil {
			log.sugar.Warnf("failed to post message to slack: monitor: %v", err)
			retry--
		}
	}
	defer resp.Body.Close()

	if retry == 0 {
		return errors.New("failed to post message to slack 3 retry")
	}

	return nil
}

// NewSlackAttachmentMessage ... Initialize attachment data of slack for failure messages for cluster
func (c *Cluster) NewSlackAttachmentMessage(message, _ string) *Attachment {
	return &Attachment{
		Color:  ColorRED,
		Title:  "ECS cluster task count threshold has been exceeded",
		Text:   fmt.Sprintf("current: %s > threshold: %d", message, c.TaskThreshold),
		Footer: fmt.Sprintf("%s: %s cluster: %s", c.AwsProfile, c.Name, c.AwsRegion),
	}
}

// NewSlackAttachmentMessage ... Initialize attachment data of slack for failure messages for task
func (t *Task) NewSlackAttachmentMessage(message, awsProfile string) *Attachment {
	// clusterArn is e.g. "arn:aws:ecs:ap-northeast-1:123456789:cluster/cron"
	cluster := strings.SplitAfter(aws.StringValue(t.EcsDescribeTask[0].ClusterArn), "/")[1]
	region := strings.Split(aws.StringValue(t.EcsDescribeTask[0].ClusterArn), ":")[3]

	return &Attachment{
		Color: ColorRED,
		Title: "ECS task parallel",
		Text: fmt.Sprintf("%s\n"+
			"current: %s > threshold: %d", t.Name, message, t.Count),
		Footer: fmt.Sprintf("%s: %s cluster: %s", awsProfile, cluster, region),
	}
}
