package ecstaskmonitoring

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"fmt"

	"go.uber.org/zap"
)

func init() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	s := logger.Sugar()
	log = Logger{sugar: s}
}

// Run ... run ecs-update-notify
func Run(interval int64, parallel int, dir string) int {
	MonitorInterval = time.Duration(interval)
	DefaultParallelTaskCount = parallel

	ctx, cancel := context.WithCancel(context.Background())

	log.sugar.Infof("start %s...\n", AppName)

	// load config
	config, err := LoadToml(dir)
	if err != nil {
		log.sugar.Fatalf("failed to load toml directory: %s", dir)
	}

	// validation config
	if err := config.Validation(); err != nil {
		log.sugar.Fatalf("config validation: %v", err)
	}

	// start cron-monitoring
	exitErrCh := make(chan error)
	go config.Start(ctx, exitErrCh)

	// receive syscall
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	go signalHandler(sig, cancel)

	for {
		select {
		case err := <-exitErrCh:
			if err != nil {
				log.sugar.Fatal(err)
			}
			return 0
		}
	}
}

// Start ... start monitor ecs cluster
// like : https://play.golang.org/p/mUDHFuYlEAY
func (c *Config) Start(ctx context.Context, exitErrCh chan error) {
	cluster := make(chan *Cluster)

	for _, v := range c.Clusters {
		var err error
		if v.Client, err = v.NewClient(); err != nil {
			exitErrCh <- fmt.Errorf("initialize ecs client error: %v", err)
			return
		}

		// infinity loop
		go func(v *Cluster) {
			for {
				v.Monitor(exitErrCh)
				// Execute again at intervals.
				time.Sleep(MonitorInterval * time.Second)
			}
		}(v)
	}

	// receive context
	select {
	case <-ctx.Done():
		log.sugar.Info("received done, exiting in 500 milliseconds")
		time.Sleep(500 * time.Millisecond)
		close(cluster)
		exitErrCh <- nil
		return
	}
}

// signalHandler ... Receive signal handler and do context.cancel
func signalHandler(sig chan os.Signal, cancel context.CancelFunc) {
	for {
		select {
		case s := <-sig:
			switch s {
			case syscall.SIGINT:
				log.sugar.Info("received SIGINT signal")
				log.sugar.Info("shutdown...")
				cancel()
			case syscall.SIGTERM:
				log.sugar.Info("received SIGTERM signal")
				log.sugar.Info("shutdown...")
				cancel()
			}
		}
	}
}
