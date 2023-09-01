package docker

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"strconv"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/ozraru/ContainerStarterBot/config"
)

var apicli client.APIClient

var ErrAlreadyStarted = errors.New("container already started")
var ErrAlreadyStopped = errors.New("container already stopped")

func init() {
	var err error
	apicli, err = client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		log.Panic("Failed to make docker client: ", err)
	}
}

func StartContainer() error {
	res, err := apicli.ContainerInspect(context.Background(), config.Config.Container)
	if err != nil {
		return err
	}
	if res.State.Running {
		return ErrAlreadyStarted
	}
	return apicli.ContainerStart(context.Background(), config.Config.Container, types.ContainerStartOptions{})
}

func StopContainer() error {
	res, err := apicli.ContainerInspect(context.Background(), config.Config.Container)
	if err != nil {
		return err
	}
	if !res.State.Running {
		return ErrAlreadyStopped
	}
	return apicli.ContainerStop(context.Background(), config.Config.Container, container.StopOptions{
		Timeout: &config.Config.Timeout,
	})
}

func ContainerStatus() (string, error) {
	data, err := apicli.ContainerInspect(context.Background(), config.Config.Container)
	if err != nil {
		return "", err
	}
	res := fmt.Sprintf("Status: `%s`\n", data.State.Status)
	if data.State.Running {
		res += fmt.Sprintf("Started at: %s\n", convertTime(data.State.StartedAt))
	} else {
		res += fmt.Sprintf("Exit code: `%d`\n", data.State.ExitCode)
		res += fmt.Sprintf("Finished at: %s\n", convertTime(data.State.FinishedAt))
	}
	if data.State.Error != "" {
		res += fmt.Sprintf("Error: `%s`\n", data.State.Error)
	}
	return res, nil
}

func convertTime(before string) string {
	t, err := time.Parse(time.RFC3339Nano, before)
	if err != nil {
		log.Print("Failed to parse time: ", err)
		return before
	}
	return fmt.Sprintf("<t:%d:T>", t.Unix())
}

func GetLog(timestamps bool, tail int64) (io.ReadCloser, error) {
	return apicli.ContainerLogs(context.Background(), config.Config.Container, types.ContainerLogsOptions{
		ShowStdout: true,
		Timestamps: timestamps,
		Tail:       strconv.FormatInt(tail, 10),
	})
}
