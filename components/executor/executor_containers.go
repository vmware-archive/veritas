package executor

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/cloudfoundry-incubator/executor"
	"github.com/cloudfoundry-incubator/executor/http/client"

	"github.com/pivotal-cf-experimental/veritas/say"
)

func ExecutorContainers(executorAddr string, raw bool, out io.Writer) error {
	client := client.New(http.DefaultClient, executorAddr)
	containers, err := client.ListContainers(nil)
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

	if len(containers) == 0 {
		say.Println(0, say.Red("No Containers"))
	}
	for _, container := range containers {
		printContainer(out, container)
	}
	return nil
}

func printContainer(out io.Writer, container executor.Container) {
	ports := []string{}
	for _, portMapping := range container.Ports {
		ports = append(ports, fmt.Sprintf("%d:%d", portMapping.HostPort, portMapping.ContainerPort))
	}
	say.Fprintln(out, 0,
		"%s (%d MB, %d MB) [%s %s] %s",
		say.Green(container.Guid),
		container.MemoryMB,
		container.DiskMB,
		container.State,
		time.Since(time.Unix(0, container.AllocatedAt)),
		strings.Join(ports, ","),
	)
}
