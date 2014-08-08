package executor_commands

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/cloudfoundry-incubator/executor/api"
	"github.com/cloudfoundry-incubator/executor/client"
	"github.com/cloudfoundry-incubator/veritas/say"
)

func ExecutorResources(executorAddr string, raw bool, out io.Writer) error {
	client := client.New(http.DefaultClient, executorAddr)
	remaining, err := client.RemainingResources()
	if err != nil {
		return err
	}
	total, err := client.TotalResources()
	if err != nil {
		return err
	}

	if raw {
		encoded, err := json.MarshalIndent(struct {
			RemainingResources api.ExecutorResources `json:"remaining_resources"`
			TotalResources     api.ExecutorResources `json:"total_resources"`
		}{remaining, total}, "", "  ")

		if err != nil {
			return err
		}

		out.Write(encoded)
		return nil
	}

	say.Fprintln(out, 0, say.Green("Resource Usage"))
	printResource(out, "Memory (MB)", total.MemoryMB, remaining.MemoryMB)
	printResource(out, "Disk (MB)", total.DiskMB, remaining.DiskMB)
	printResource(out, "Containers", total.Containers, remaining.Containers)
	return nil
}

func printResource(out io.Writer, label string, total int, remaining int) {
	used := total - remaining
	usedString := fmt.Sprintf("%d", used)
	usedPercentage := fmt.Sprintf("%.1f%%", float64(used)/float64(total)*100.0)

	if float64(used)/float64(total) > 0.8 {
		usedString = say.Red(usedString)
		usedPercentage = say.Red(usedPercentage)
	} else {
		usedString = say.Green(usedString)
		usedPercentage = say.Green(usedPercentage)
	}

	say.Fprintln(out, 1, "%s: %s/%s (%s)", label, usedString, say.Green("%d", total), usedPercentage)
}
