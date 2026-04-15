package helpers

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
)

var dockerClient *client.Client

func init() {
	var err error
	dockerClient, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		Warning("Failed to initialize Docker client: %v", err)
	}
}

// SpawnWorkspaceContainer spawns a new Docker container and returns the container ID and the internal target URL for the VNC websocket.
// Example imageName: "kasmweb/ubuntu-bionic-desktop:1.10.0"
func SpawnWorkspaceContainer(imageName string, prefix string) (string, string, error) {
	if dockerClient == nil {
		return "", "", fmt.Errorf("docker client is not initialized")
	}

	ctx := context.Background()

	// Ensure image is pulled
	_, _, err := dockerClient.ImageInspectWithRaw(ctx, imageName)
	if err != nil {
		Warning("Image %s not found locally or inspect failed: %v. Attempting to pull...", imageName, err)
		reader, err := dockerClient.ImagePull(ctx, imageName, image.PullOptions{})
		if err != nil {
			return "", "", fmt.Errorf("failed to pull image %s: %w", imageName, err)
		}
		defer reader.Close()
		
		if _, err := io.Copy(io.Discard, reader); err != nil {
			return "", "", fmt.Errorf("error waiting for image pull of %s: %w", imageName, err)
		}
		Success("Image %s pulled successfully.", imageName)
	}

	containerName := fmt.Sprintf("%s-%d", prefix, time.Now().UnixNano())

	// We don't expose ports to the host natively because the Go app will route directly to the container's IP on the bridge network.
	// 6901 is the default Websockify VNC port for kasm web images.
	resp, err := dockerClient.ContainerCreate(ctx, &container.Config{
		Image: imageName,
		Env:   []string{"VNC_PW=password"}, // Default password, can be randomized
	}, &container.HostConfig{
		AutoRemove: true, // Automatically cleanup when stopped
	}, &network.NetworkingConfig{}, nil, containerName)
	if err != nil {
		return "", "", fmt.Errorf("failed to create container: %w", err)
	}

	if err := dockerClient.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return "", "", fmt.Errorf("failed to start container: %w", err)
	}

	// Inspect container to grab its internal Bridge IP
	inspect, err := dockerClient.ContainerInspect(ctx, resp.ID)
	if err != nil {
		// Attempt to cleanup if inspect fails
		_ = StopWorkspaceContainer(resp.ID)
		return "", "", fmt.Errorf("failed to inspect container: %w", err)
	}

	ipAddress := inspect.NetworkSettings.IPAddress
	if ipAddress == "" && len(inspect.NetworkSettings.Networks) > 0 {
		for _, net := range inspect.NetworkSettings.Networks {
			ipAddress = net.IPAddress
			if ipAddress != "" {
				break
			}
		}
	}

	if ipAddress == "" {
		_ = StopWorkspaceContainer(resp.ID)
		return "", "", fmt.Errorf("container started but no IP address was assigned")
	}

	// URL inside the Docker bridge network where websockify is listening.
	// You might need to change the port based on the image you actually build.
	targetURL := fmt.Sprintf("https://%s:6901", ipAddress)

	return resp.ID, targetURL, nil
}

// StopWorkspaceContainer sends a stop signal to the container ID. Since AutoRemove: true, it also deletes it.
func StopWorkspaceContainer(containerID string) error {
	if dockerClient == nil {
		return fmt.Errorf("docker client is not initialized")
	}
	ctx := context.Background()
	// No timeout provided, waits gracefully
	return dockerClient.ContainerStop(ctx, containerID, container.StopOptions{})
}
