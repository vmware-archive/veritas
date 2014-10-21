package submit_task

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"time"

	"github.com/cloudfoundry-incubator/receptor"
	"github.com/cloudfoundry-incubator/runtime-schema/models"
	"github.com/pivotal-cf-experimental/veritas/say"
)

func SubmitTask(client receptor.Client, f io.Reader) error {
	var desiredTask receptor.CreateTaskRequest

	if f != nil {
		decoder := json.NewDecoder(f)
		err := decoder.Decode(&desiredTask)
		if err != nil {
			return err
		}
	} else {
		desiredTask = interactivelyBuildDesiredTask()

		asJson, err := json.Marshal(desiredTask)
		if err != nil {
			return err
		}

		filename := fmt.Sprintf("desired_task_%s.json", desiredTask.TaskGuid)
		err = ioutil.WriteFile(filename, asJson, 0666)
		if err != nil {
			return err
		}
		say.Println(0, say.Green("Save DesiredTask to %s", filename))
	}

	say.Println(0, say.Green("Desiring %s", desiredTask.TaskGuid))
	return client.CreateTask(desiredTask)
}

/*
type CreateTaskRequest struct {
    TaskGuid   string                  `json:"task_guid"`
    Domain     string                  `json:"domain"`
    Actions    []models.ExecutorAction `json:"actions"`
    Stack      string                  `json:"stack"`
    MemoryMB   int                     `json:"memory_mb"`
    DiskMB     int                     `json:"disk_mb"`
    CpuPercent float64                 `json:"cpu_percent"`
    Log        models.LogConfig        `json:"log"`
    Annotation string                  `json:"annotation,omitempty"`
}
*/

func interactivelyBuildDesiredTask() receptor.CreateTaskRequest {
	desiredTask := receptor.CreateTaskRequest{}
	desiredTask.TaskGuid = say.AskWithDefault("TaskGuid", fmt.Sprintf("%d", time.Now().Unix()))
	desiredTask.Domain = say.AskWithDefault("Domain", "veritas")
	desiredTask.Stack = say.AskWithDefault("Stack", "lucid64")
	desiredTask.MemoryMB = say.AskForIntegerWithDefault("MemoryMB", 256)
	desiredTask.DiskMB = say.AskForIntegerWithDefault("DiskMB", 256)
	desiredTask.CpuPercent = float64(say.AskForIntegerWithDefault("CpuPercent", 100))
	desiredTask.Log = models.LogConfig{
		Guid:       desiredTask.TaskGuid,
		SourceName: "VRT",
	}
	desiredTask.Annotation = say.AskWithDefault("Annotation", "none")
	desiredTask.Actions = interactivelyBuildActions()

	return desiredTask
}

func interactivelyBuildActions() []models.ExecutorAction {
	actions := []models.ExecutorAction{}
	for {
		choice := say.Pick("Add an action", []string{
			"Done",
			"DownloadAction",
			"RunAction",
		})

		switch choice {
		case "Done":
			return actions
		case "DownloadAction":
			actions = append(actions, interactivelyBuildDownloadAction())
		case "RunAction":
			actions = append(actions, interactivelyBuildRunAction())
		}
	}
}

func interactivelyBuildDownloadAction() models.ExecutorAction {
	return models.ExecutorAction{
		models.DownloadAction{
			From: say.Ask("Download URL"),
			To:   say.AskWithDefault("Container Destination", "."),
		},
	}
}

func interactivelyBuildRunAction() models.ExecutorAction {
	return models.ExecutorAction{
		models.RunAction{
			Path: say.AskWithValidation("Command to run", func(response string) error {
				if strings.Contains(response, " ") {
					return fmt.Errorf("You cannot specify arguments to the command, that'll come next...")
				}
				return nil
			}),
			Args: strings.Split(say.Ask("Args (split by ';')"), ";"),
			Env:  interactivelyGetEnvs("Envs (FOO=BAR;BAZ=WIBBLE)"),
		},
	}
}

func interactivelyGetEnvs(prompt string) []models.EnvironmentVariable {
	envs := say.Ask(prompt)
	splitEnvs := strings.Split(envs, ";")
	out := []models.EnvironmentVariable{}
	for _, env := range splitEnvs {
		sub := strings.Split(env, "=")
		if len(sub) == 2 {
			out = append(out, models.EnvironmentVariable{
				Name:  sub[0],
				Value: sub[1],
			})
		}
	}
	return out
}
