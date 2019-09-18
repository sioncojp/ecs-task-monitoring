package main

import (
	"flag"
	"os"

	"github.com/sioncojp/ecs-task-monitoring"
)

func main() {
	interval := flag.Int64("i", 15, "check interval (second)")
	parallel := flag.Int("p", 1, "default number of tasks that can move in parallel")
	parallelNotifyInterval := flag.Int("n", 15, "interval to notify the summary of parallel tasks")
	dir := flag.String("d", "", "directory where toml file is stored")
	flag.Parse()
	os.Exit(ecstaskmonitoring.Run(*interval, *parallel, *parallelNotifyInterval, *dir))
}
