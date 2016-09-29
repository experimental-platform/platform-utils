package dockerutil

import (
	"fmt"

	"github.com/docker/engine-api/client"
	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/filters"
	"golang.org/x/net/context"
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
		protonetNetworkData, ok = data.NetworkSettings.Networks["docker"]
		if !ok {
			return "", fmt.Errorf("The container '%s' doesn't exist on the networks 'protonet' and 'docker'.", name)
		}
	}

	return protonetNetworkData.IPAddress, nil
}
