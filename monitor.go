package ecstaskmonitoring

import (
	"fmt"
	"strconv"
	"sync"
	"time"
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
	go MonitorTaskThreshold(*c, &wg)
	go MonitorTaskParallel(*c, &wg)
	wg.Wait()

	return
}

// MonitorTaskThreshold ... Monitor the number of tasks in the cluster
func MonitorTaskThreshold(c Cluster, wg *sync.WaitGroup) {
	defer wg.Done()

	CurrentTaskThresholdFailureCount++

	if len(c.Tasks) > c.TaskThreshold {
		if (CurrentTaskThresholdFailureCount - c.FailureCount) == 0 {
			a := c.NewSlackAttachmentMessage(strconv.Itoa(len(c.Tasks)))
			a.PostSlackMessage(c.IncomingWebhook)
			CurrentTaskThresholdFailureCount = 0
		}
	}
	return
}

// MonitorTaskParallel ... Monitor if task is running in parallel
func MonitorTaskParallel(c Cluster, wg *sync.WaitGroup) {
	defer wg.Done()

	for _, v := range c.Tasks {
		// if setting count 0, ignore check
		if v.Count == 0 {
			continue
		}

		if len(v.EcsDescribeTasks) > v.Count {
			CountUpParallels(v.Name)
		}
	}

	// Notify every 60 minutes
	if time.Now().Sub(ParallelNotifyTime).Minutes() > float64(ParallelNotifyTimeInterval) {
		notify := ParallelsToParallelNotify(c)
		if notify.Message != "" {
			a := notify.NewSlackAttachmentMessage("")
			a.PostSlackMessage(c.IncomingWebhook)
		}

		// after notify, clear time and Parallels data
		ParallelNotifyTime = time.Now()
		Parallels = []*Parallel{}
	}

}

// CountUpParallels ... count up Parallels variable.
func CountUpParallels(taskName string) {
	for _, v := range Parallels {
		if v.Name == taskName {
			v.Count++
			return
		}
	}
	Parallels = append(Parallels, &Parallel{taskName, 1})
}

// ParallelsToParallelNotify ... Parallels struct To ParallelNotify
func ParallelsToParallelNotify(c Cluster) *ParallelNotify {
	result := &ParallelNotify{
		ClusterName: c.Name,
		AwsProfile:  c.AwsProfile,
		AwsRegion:   c.AwsRegion,
	}
	for _, v := range Parallels {
		result.Message += fmt.Sprintf("%s: %d count\n", v.Name, v.Count)
	}

	return result
}
