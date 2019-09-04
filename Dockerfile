# build stage
FROM golang:1.12 AS build-env

ENV GO111MODULE auto
ENV CGO_ENABLED=0

ADD . /src
WORKDIR /src
RUN make build

# final stage
FROM alpine

WORKDIR /app
RUN mkdir toml
COPY --from=build-env /src/bin/ecs-task-monitoring /app/
ENTRYPOINT ./ecs-task-monitoring -d toml/ -i $Interval -p $DefaultParallelCount
