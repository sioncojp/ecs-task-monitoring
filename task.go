package ecstaskmonitoring

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
)

// NewTask ... initialize Task
func (c *Cluster) NewTask(e *ecs.Task) (*Task, error) {
	fmt.Println("new task")
	return &Task{
		Name:            strings.TrimPrefix(aws.StringValue(e.Group), "family:"),
		IncomingWebhook: c.IncomingWebhook,
		Count:           DefaultParallelTaskCount,
		EcsDescribeTask: []*ecs.Task{e},
	}, nil
}

// NewTasks ... initialize only "RUNNING" current tasks
func (c *Cluster) NewTasks() ([]*Task, error) {
	var result []*Task

	listTask, err := c.ListTasks()
	if err != nil {
		return nil, err
	}

	if listTask == nil {
		return nil, nil
	}

	describeTasks, err := c.DescribeTasks(listTask)
	if err != nil {
		return nil, err
	}

	for _, d := range describeTasks {
		if aws.StringValue(d.LastStatus) == "RUNNING" {
			// Branch for incoming webhook for each service
			if c.Tasks != nil {
				for _, t := range c.Tasks {
					if t.Name == strings.TrimPrefix(aws.StringValue(d.Group), "family:") {
						if t.IncomingWebhook == "" {
							t.IncomingWebhook = c.IncomingWebhook
						}
						t.EcsDescribeTask = append(t.EcsDescribeTask, d)

						result = append(result, t)

					} else {
						// If nothing is set, apply incoming webhook set to cluster
						t, err := c.NewTask(d)
						if err != nil {
							return nil, err
						}
						result = append(result, t)
					}
				}
			}
		}
	}

	return result, nil
}
