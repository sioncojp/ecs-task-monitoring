# ecs-task-monitoring

1. Monitor the number of tasks in the cluster and notify with slack IncomingWebhook.
![](docs/ecs-task-monitoring01.png)

2. Monitor tasks running in parallel. To notify every 60 minutes.
![](docs/ecs-task-monitoring02.png)

# Motivation

If we run ecs task with cloudwatch events, we cannot monitor the number of tasks running in the cluster with datadog.

If task is run as cron in cloudwatch events, datadog cannot monitor whether it is moving in parallel.

# Development

```shell
$ make help
$ make build
$ ./bin/ecs-task-monitoring -d examples/ -i 3 -p 1
```

# License
The MIT License

Copyright Shohei Koyama / sioncojp

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.