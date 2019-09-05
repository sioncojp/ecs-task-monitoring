package ecstaskmonitoring

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"time"

	"github.com/BurntSushi/toml"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/ecs/ecsiface"
	"go.uber.org/zap"
)

const (
	// AppName ... This tool name.
	AppName = "ecs-task-monitoring"

	// ColorRED ... Red color code
	ColorRED = "#F08080"
)

var (
	log Logger

	// MonitorInterval ... Monitoring interval
	MonitorInterval time.Duration

	// DefaultTaskCount ... Default number of tasks that can move in parallel
	DefaultParallelTaskCount int

	CurrentTaskThresholdFailureCount = 0
)

// Logger ... Store logging
type Logger struct {
	sugar *zap.SugaredLogger
}

// Config ... Store from xxxx.toml
type Config struct {
	Clusters []*Cluster `toml:"cluster"`
}

// Cluster ... Store ecs cluster data
type Cluster struct {
	Name            string  `toml:"name"`
	FailureCount    int     `toml:"failure_count"`
	TaskThreshold   int     `toml:"task_threshold"`
	AwsProfile      string  `toml:"aws_profile"`
	AwsRegion       string  `toml:"aws_region"`
	Tasks           []*Task `toml:"task"`
	IncomingWebhook string  `toml:"incoming_webhook"`
	Client          *Client
	Slack
}

// Task ... Store ecs task data
type Task struct {
	Name            string `toml:"name"`
	IncomingWebhook string `toml:"incoming_webhook"`
	Count           int    `toml:"count"`
	EcsDescribeTask []*ecs.Task
	Slack
}

// Client ... Store ECS client with a session
type Client struct {
	ecs ecsiface.ECSAPI
}

// SlackMessage ... Store slack message data
type SlackMessage struct {
	Attachments []*Attachment `json:"attachments"`
}

// Attachment ... Slack Attachment Data
type Attachment struct {
	Color  string `json:"color,omitempty"`
	Title  string `json:"title,omitempty"`
	Text   string `json:"text,omitempty"`
	Footer string `json:"footer,omitempty"`
}

// LoadToml ... Read the toml file in the directory
func LoadToml(dir string) (*Config, error) {
	// 末尾が / で終わってなければ追加
	if string(dir[len(dir)-1]) != "/" {
		dir = dir + "/"
	}

	// load config. ディレクトリ配下の設定ファイルを結合して読み込む
	files, _ := ioutil.ReadDir(dir)
	openFiles := make([]io.Reader, len(files)*2)

	// ファイル間の結合の際、改行を加える
	for i := 0; i < len(files); i++ {
		num := int(2 * float64(i))
		if i == 0 {
			num = 0
		}
		openFiles[num], _ = os.Open(fmt.Sprintf("%s%s", dir, files[i].Name()))
		openFiles[num+1] = strings.NewReader("\n")
	}

	reader := io.MultiReader(openFiles...)

	var config Config
	if _, err := toml.DecodeReader(reader, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

// Validation ... Validation that the value set in toml is correct
func (c *Config) Validation() error {
	invalidMessage := "invalid parameter"
	for _, v := range c.Clusters {
		if v.FailureCount == 0 {
			return fmt.Errorf("%s: failure_count is 0", invalidMessage)
		}

		if v.TaskThreshold == 0 {
			return fmt.Errorf("%s: task_threshold is 0", invalidMessage)
		}

		if v.TaskThreshold == 0 {
			return fmt.Errorf("%s: aws region is 0", invalidMessage)
		}

		if v.Name == "" {
			return fmt.Errorf("%s: name is empty", invalidMessage)
		}

		if v.AwsRegion == "" {
			return fmt.Errorf("%s: aws_region is empty", invalidMessage)
		}

		if v.IncomingWebhook == "" {
			return fmt.Errorf("%s: incoming_webhook is empty", invalidMessage)
		}

	}
	return nil
}
