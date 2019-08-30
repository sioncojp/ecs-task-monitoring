package ecstaskmonitoring

import (
	"strconv"
	"sync"
)

// Monitor ... Monitor the cluster
func (c *Cluster) Monitor(exitErrCh chan error) {
	tasks, err := c.NewTasks()
	if err != nil {
		exitErrCh <- err
		return
	}

	if tasks == nil {
		return
	}

	c.Tasks = tasks

	// concurrency
	var wg sync.WaitGroup
	wg.Add(2)
	go MonitorTaskThreshold(*c, exitErrCh, &wg)
	go MonitorTaskParallel(*c, exitErrCh, &wg)
	wg.Wait()

	return
}

// MonitorTaskThreshold ... Monitor the number of tasks in the cluster
func MonitorTaskThreshold(c Cluster, exitErrCh chan error, wg *sync.WaitGroup) {
	defer wg.Done()

	if len(c.Tasks) > c.TaskThreshold {
		a := c.NewSlackAttachmentMessage(strconv.Itoa(len(c.Tasks)), "")
		if err := a.PostSlackMessage(c.IncomingWebhook); err != nil {
			exitErrCh <- err
			return
		}
	}
}

// MonitorTaskParallel ... Monitor if task is running in parallel
func MonitorTaskParallel(c Cluster, exitErrCh chan error, wg *sync.WaitGroup) {
	defer wg.Done()

	for _, v := range c.Tasks {
		if v.Count == 0 {
			return
		}

		if len(c.Tasks) > v.Count {
			a := v.NewSlackAttachmentMessage(strconv.Itoa(len(c.Tasks)), c.AwsProfile)
			if err := a.PostSlackMessage(v.IncomingWebhook); err != nil {
				exitErrCh <- err
				return
			}
		}
	}
}
