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
			log.sugar.Info("detect %s cluster in %s task threshold: %d", c.Name, c.AwsProfile, len(c.Tasks))

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
		if v.Times == 0 {
			continue
		}

		if len(v.EcsDescribeTasks) > v.Times {
			IncreaseTimesParallels(v.Name)

			log.sugar.Infof("detect parallel: %s", v.Name)
			for _, v := range v.EcsDescribeTasks {
				log.sugar.Infof("createdAt: %d, stoppedAt: %d, taskArn: %s", *v.CreatedAt, *v.StoppedAt, *v.TaskArn)
			}
		}
	}

	// Notify every 60 minutes
	if time.Now().Sub(ParallelNotifyTime).Minutes() >= float64(ParallelNotifyInterval) {
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

// IncreaseTimesParallels ... Increase in number of times Parallels variable.
func IncreaseTimesParallels(taskName string) {
	for _, v := range Parallels {
		if v.Name == taskName {
			v.Times++
			return
		}
	}
	Parallels = append(Parallels, &Parallel{taskName, 1})
}

// ParallelsToParallelNotify ... Parallels struct to ParallelNotify
func ParallelsToParallelNotify(c Cluster) *ParallelNotify {
	result := &ParallelNotify{
		ClusterName: c.Name,
		AwsProfile:  c.AwsProfile,
		AwsRegion:   c.AwsRegion,
	}
	for _, v := range Parallels {
		result.Message += fmt.Sprintf("%s: %d times\n", v.Name, v.Times)
	}

	return result
}
