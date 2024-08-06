package rugs

import (
	"context"
	"fmt"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	//dcontainer "github.com/docker/docker/api/types/container"
)

type rugsContainer struct {
	testcontainers.Container
	URIBase string
}

func startRugsContainer(ctx context.Context) (*rugsContainer, error) {
	// the RUGS -- Rust UnrealGameSync server -- automatically runs a schema apply/update on start,
	// which -- for ephemeral containers -- means to create the sqlite database
	//
	// starup therefore takes ~2s, but will respond to a healthcheck (http://host:3000/health) once
	// available
	req := testcontainers.ContainerRequest{
		Image:        "ghcr.io/jorgenpt/rugs:latest",
		ExposedPorts: []string{"3000/tcp"},
		WaitingFor:   wait.ForHTTP("/health"),
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, fmt.Errorf("ContainerRequest: %v", err)
	}

	ip, err := container.Host(ctx)
	if err != nil {
		return nil, fmt.Errorf("Container.Host(%v): %v", ip, err)
	}

	mappedPort, err := container.MappedPort(ctx, "3000")
	if err != nil {
		return nil, fmt.Errorf("Container.MappedPort(3000): %v", err)
	}

	uri := fmt.Sprintf("http://%s:%s", ip, mappedPort.Port())

	return &rugsContainer{Container: container, URIBase: uri}, nil
}
