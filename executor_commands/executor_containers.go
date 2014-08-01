package executor_commands

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/cloudfoundry-incubator/executor/api"
	"github.com/cloudfoundry-incubator/executor/client"
	"github.com/cloudfoundry-incubator/veritas/say"
)

func ExecutorContainers(executorAddr string, raw bool, out io.Writer) error {
	client := client.New(http.DefaultClient, executorAddr)
	containers, err := client.ListContainers()
	if err != nil {
		return err
	}

	if raw {
		encoded, err := json.MarshalIndent(containers, "", "  ")

		if err != nil {
			return err
		}

		out.Write(encoded)
		return nil
	}

	say.Fprintln(out, 0, say.Green("Containers"))
	if len(containers) == 0 {
		say.Println(0, say.Red("None"))
	}
	for _, container := range containers {
		printContainer(out, container)
	}
	return nil
}

func printContainer(out io.Writer, container api.Container) {
	ports := []string{}
	for _, portMapping := range container.Ports {
		ports = append(ports, fmt.Sprintf("%d:%d", portMapping.HostPort, portMapping.ContainerPort))
	}
	say.Fprintln(out, 0,
		"%s@%s (%d MB, %d MB) [%s %s] %s",
		say.Green(container.Guid),
		say.Cyan(container.ContainerHandle),
		container.MemoryMB,
		container.DiskMB,
		container.State,
		time.Since(time.Unix(0, container.AllocatedAt)),
		strings.Join(ports, ","),
	)
}
