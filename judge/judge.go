package judge

import (
	"bytes"
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"os/exec"
	"time"
)

func EvaluateUserCode(userTarBuffer *bytes.Buffer, evaluatorImage string) (string, error) {
	ctx := context.Background()

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return "", fmt.Errorf("unable to create Docker client: %w", err)
	}

	cfg := &container.Config{
		Image: evaluatorImage,
		Cmd:   []string{"bash", "./doevaluate.sh"},
	}
	hostConfig := &container.HostConfig{
		NetworkMode: container.NetworkMode("none"), // 使用 none 网络模式，禁止容器访问外部网络
	}
	resp, err := cli.ContainerCreate(ctx, cfg, hostConfig, nil, nil, "")
	if err != nil {
		return "", fmt.Errorf("unable to create container: %w", err)
	}

	defer func() {
		// Ensure the container is removed when we are done.
		if err := cli.ContainerRemove(ctx, resp.ID, types.ContainerRemoveOptions{Force: true}); err != nil {
			fmt.Println("warning: failed to remove container:", err)
		}
	}()

	reader := bytes.NewReader(userTarBuffer.Bytes())
	err = cli.CopyToContainer(ctx, resp.ID, "/", reader, types.CopyToContainerOptions{
		AllowOverwriteDirWithFile: true,
		CopyUIDGID:                false,
	})
	if err != nil {
		println("copy file error")
		return "", err
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return "", fmt.Errorf("unable to start container: %w", err)
	}

	// Set a time limit for the code execution.
	done := make(chan error, 1)
	go func() {
		statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
		select {
		case err := <-errCh:
			if err != nil {
				done <- fmt.Errorf("error waiting for container: %w", err)
			}
		case <-statusCh:
		}
		close(done)
	}()

	// If the execution takes more than 2 seconds, consider it as timeout.
	select {
	case <-time.After(10 * time.Second):
		cmd := exec.Command("docker", "logs", resp.ID)
		outLog, err := cmd.Output()
		if err != nil {
			println("get docker log fail:", err.Error())
		}
		println(outLog)
		return "timeout", fmt.Errorf("execution timed out")
	case err := <-done:
		if err != nil {
			return "", fmt.Errorf("container run error: %w", err)
		}
	}

	cmd := exec.Command("docker", "logs", resp.ID)
	outLog, err := cmd.Output()
	if err != nil {
		println("get docker log fail:", err.Error())
		return "", err
	}
	return string(outLog), nil
}
