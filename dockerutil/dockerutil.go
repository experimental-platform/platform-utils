package dockerutil

import (
	"github.com/docker/engine-api/types/filters"
	"fmt"
	"github.com/docker/engine-api/client"
	"github.com/docker/engine-api/types"
	"golang.org/x/net/context"
	"errors"
)

func GetContainerIP(name string) (string, error) {
	defaultHeaders := map[string]string{"User-Agent": "protonet-skvs_cli"}
	cli, err := client.NewClient("unix:///var/run/docker.sock", "v1.22", nil, defaultHeaders)
	if err != nil {
		return "", err
	}

	listOptions := types.ContainerListOptions{Filter: filters.NewArgs()}
	listOptions.Filter.Add("name", name)

	containers, err := cli.ContainerList(context.Background(), listOptions)
	if err != nil {
		return "", err
	}
	if len(containers) == 0 {
		return "", fmt.Errorf("Found no container named '%s'", name)
	}

	data, err := cli.ContainerInspect(context.Background(), containers[0].ID)
	if err != nil {
		return "", err
	}

	protonetNetworkData, ok := data.NetworkSettings.Networks["protonet"]
	if !ok {
		return "", errors.New("The SKVS container doesn't belong to the network 'protonet'.")
	}

	return protonetNetworkData.IPAddress, nil
}
