package submit_task

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"time"

	"github.com/cloudfoundry-incubator/receptor"
	"github.com/pivotal-cf-experimental/veritas/common"
	"github.com/pivotal-cf-experimental/veritas/say"
)

func SubmitTask(client receptor.Client, f io.Reader) error {
	var desiredTask receptor.TaskCreateRequest

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

func interactivelyBuildDesiredTask() receptor.TaskCreateRequest {
	desiredTask := receptor.TaskCreateRequest{}
	desiredTask.TaskGuid = say.AskWithDefault("TaskGuid", fmt.Sprintf("%d", time.Now().Unix()))
	desiredTask.Domain = say.AskWithDefault("Domain", "veritas")
	desiredTask.Stack = say.AskWithDefault("Stack", "lucid64")
	desiredTask.MemoryMB = say.AskForIntegerWithDefault("MemoryMB", 256)
	desiredTask.DiskMB = say.AskForIntegerWithDefault("DiskMB", 256)
	desiredTask.CPUWeight = uint(say.AskForIntegerWithDefault("CPUWeight", 100))
	desiredTask.EnvironmentVariables = common.ModelEnvsToReceptorEnvs(common.BuildEnvs())
	desiredTask.LogGuid = desiredTask.TaskGuid
	desiredTask.LogSource = "VRT"
	desiredTask.Annotation = say.AskWithDefault("Annotation", "none")
	desiredTask.Action = common.BuildAction("Build Action", nil)

	return desiredTask
}
