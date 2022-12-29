package docker

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"strconv"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/ozraru/ContainerStarterBot/config"
)

var apicli client.APIClient

var ErrAlreadyStarted = errors.New("container already started")

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

func ContainerStatus() (string, error) {
	data, err := apicli.ContainerInspect(context.Background(), config.Config.Container)
	if err != nil {
		return "", err
	}
	res := fmt.Sprintf("Status: %s\n", data.State.Status)
	if data.State.Dead {
		res += fmt.Sprintf("Exit code: %d\n", data.State.ExitCode)
	}
	if data.State.Running {
		res += fmt.Sprintf("Started at: %s\n", data.State.StartedAt)
	} else {
		res += fmt.Sprintf("Finished at: %s\n", data.State.FinishedAt)
	}
	if data.State.Error != "" {
		res += fmt.Sprintf("Error: %s\n", data.State.Error)
	}
	return res, nil
}

func GetLog(timestamps bool, tail int64) (io.ReadCloser, error) {
	return apicli.ContainerLogs(context.Background(), config.Config.Container, types.ContainerLogsOptions{
		ShowStdout: true,
		Timestamps: timestamps,
		Tail:       strconv.FormatInt(tail, 10),
	})
}
