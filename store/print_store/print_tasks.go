package print_store

import (
	"time"

	"github.com/cloudfoundry-incubator/bbs/models"

	"github.com/onsi/gomega/format"
	"github.com/onsi/say"
	"github.com/pivotal-cf-experimental/veritas/veritas_models"
)

func printTasks(verbose bool, tasks veritas_models.VeritasTasks) {
	taskTypes := tasks.OrderedTaskTypes()

	say.Println(0, say.Green("Tasks"))

	for _, taskType := range taskTypes {
		say.Println(0, say.Green(taskType))
		for _, task := range tasks[taskType] {
			if verbose {
				printVerboseTask(task)
			} else {
				printTask(task)
			}
		}
	}
}

func printVerboseTask(task *models.Task) {
	say.Println(0, format.Object(task, 1))
}

func printTask(task *models.Task) {
	privileged := ""
	if task.Privileged {
		privileged = say.Red(" PRIVILEGED")
	}
	say.Println(1,
		"%s [%s on %s %s%s] U:%s C:%s (%d MB, %d MB, %d CPU)",
		taskState(task),
		task.TaskGuid,
		task.CellId,
		say.Cyan(task.RootFs),
		privileged,
		time.Since(time.Unix(0, task.UpdatedAt)).String(),
		time.Since(time.Unix(0, task.CreatedAt)).String(),
		task.MemoryMb,
		task.DiskMb,
		task.CpuWeight,
	)
}

func taskState(task *models.Task) string {
	switch task.State {
	case models.Task_Pending:
		return say.LightGray("PENDING  ")
	case models.Task_Running:
		return say.Yellow("RUNNING  ")
	case models.Task_Completed:
		return colorByTaskSuccess(task, "COMPLETED")
	case models.Task_Resolving:
		return colorByTaskSuccess(task, "RESOLVING")
	default:
		return say.Red("INVALID")
	}
}

func colorByTaskSuccess(task *models.Task, format string, args ...interface{}) string {
	if task.Failed {
		return say.Red(format, args...)
	} else {
		return say.Green(format, args...)
	}
}
