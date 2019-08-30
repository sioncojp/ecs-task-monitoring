package ecstaskmonitoring

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
)

// ListTasks ... Returns a list of tasks for a specified cluster.
func (c *Cluster) ListTasks() ([]*string, error) {
	var taskArns []*string
	input := &ecs.ListTasksInput{
		Cluster: aws.String(c.Name),
	}

	if err := c.Client.ecs.ListTasksPages(
		input,
		func(page *ecs.ListTasksOutput, _ bool) bool {
			for _, taskArn := range page.TaskArns {
				taskArns = append(taskArns, taskArn)
			}
			return true
		},
	); err != nil {
		return nil, err
	}
	return taskArns, nil
}

// DescribeTasks ... Returns a list of tasks for a specified cluster
func (c *Cluster) DescribeTasks(tasks []*string) ([]*ecs.Task, error) {
	input := &ecs.DescribeTasksInput{
		Cluster: aws.String(c.Name),
		Tasks:   tasks,
	}

	result, err := c.Client.ecs.DescribeTasks(input)
	if err != nil {
		return nil, fmt.Errorf("failed to DescribeTasks: %s cluster", c.Name)
	}

	return result.Tasks, nil
}

// NewClient ... Creates a new instance of the ECS client with a session.
func (c *Cluster) NewClient() (*Client, error) {
	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Profile:           c.AwsProfile,
	})

	if err != nil {
		return nil, err
	}

	return &Client{
		ecs: ecs.New(sess, aws.NewConfig().WithRegion(c.AwsRegion)),
	}, nil
}
