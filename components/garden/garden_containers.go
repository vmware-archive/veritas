package garden

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/cloudfoundry-incubator/garden/api"
	"github.com/cloudfoundry-incubator/garden/client"
	"github.com/cloudfoundry-incubator/garden/client/connection"

	"github.com/pivotal-cf-experimental/veritas/say"
)

type ContainerInfo struct {
	Handle string
	Info   api.ContainerInfo
}

func GardenContainers(gardenAddr string, gardenNetwork string, raw bool, out io.Writer) error {
	client := client.New(connection.New(gardenNetwork, gardenAddr))
	containers, err := client.Containers(nil)
	if err != nil {
		return err
	}

	containerInfos := []ContainerInfo{}
	for _, container := range containers {
		info, err := container.Info()
		if err != nil {
			say.Println(1, say.Red("Failed to fetch container: %s\n", container.Handle()))
			return err
		}
		containerInfos = append(containerInfos, ContainerInfo{
			container.Handle(),
			info,
		})
	}

	if raw {
		encoded, err := json.MarshalIndent(containerInfos, "", "  ")

		if err != nil {
			return err
		}

		out.Write(encoded)
		return nil
	}

	if len(containerInfos) == 0 {
		say.Println(0, say.Red("No Containers"))
	}
	for _, containerInfo := range containerInfos {
		printContainer(out, containerInfo)
	}
	return nil
}

func printContainer(out io.Writer, containerInfo ContainerInfo) {
	info := containerInfo.Info
	say.Fprintln(out, 0,
		"%s - %s @ %s",
		say.Green(containerInfo.Handle),
		info.State,
		info.ContainerPath,
	)

	say.Fprintln(out, 1,
		"Memory: %.3f MB",
		float64(info.MemoryStat.TotalRss+info.MemoryStat.TotalCache-info.MemoryStat.TotalInactiveFile)/1024.0/1024.0,
	)

	say.Fprintln(out, 1,
		"Disk: %.3f MB %d Inodes",
		float64(info.DiskStat.BytesUsed)/1024.0/1024.0,
		info.DiskStat.InodesUsed,
	)

	ports := []string{}
	for _, portMapping := range info.MappedPorts {
		ports = append(ports, fmt.Sprintf("%d:%d", portMapping.HostPort, portMapping.ContainerPort))
	}

	say.Fprintln(out, 1,
		"%s=>%s: %s",
		say.Green(info.HostIP),
		say.Green(containerInfo.Handle),
		strings.Join(ports, ","),
	)

	if len(info.Events) > 0 {
		say.Fprintln(out, 1,
			"Events: %s",
			strings.Join(info.Events, ","),
		)
	}

	if len(info.ProcessIDs) > 0 {
		say.Fprintln(out, 1,
			"Running: %d processes",
			len(info.ProcessIDs),
		)
	}

	if len(info.Properties) > 0 {
		say.Fprintln(out, 1,
			"Properties:",
		)
		for key, value := range info.Properties {
			say.Fprintln(out, 2,
				"%s: %s",
				key, value,
			)
		}
	}
}
