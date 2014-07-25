package main

import (
	"sort"
	"time"

	"github.com/cloudfoundry-incubator/runtime-schema/bbs"
	"github.com/cloudfoundry-incubator/runtime-schema/models"
	"github.com/onsi/gomega/format"
)

func DumpTasks(bbs *bbs.BBS, c Config) {
	tasks, err := bbs.GetAllTasks()
	panicIfErr(err)

	tasksByType := map[models.TaskType][]models.Task{}
	for _, task := range tasks {
		tasksByType[task.Type] = task
	}

	taskTypes := []string{}
	for key, tasks := range tasksByType {
		sort.Sort(TasksByUpdatedAt(tasks))
		taskTypes := append(taskTypes, key)
	}
	sort.Strings(taskTypes)

	c.S.printBanner(c.S.colorize("Tasks", greenColor), "~")
	for _, taskType := range taskTypes {
		c.S.println(0, c.S.colorize(greenColor, taskType))
		for _, task := range tasksByType[taskType] {
			if c.Verbose {
				dumpVerboseTask(task, c)
			} else {
				dumpTask(task, c)
			}
		}
	}
}

func dumpVerboseTask(task models.Task, c Config) {
	c.S.println(format.Object(task, 1))
}

func dumpTask(task models.Task, c Config) {
	c.S.println(1,
		"%s [%s on %s@%s(%s)] U:%s C:%s (%d MB, %d MB)",
		task.State(task.State),
		task.Guid,
		task.ContainerHandle,
		task.ExecutorID,
		task.Stack,
		time.Since(time.Unix(0, task.UpdatedAt)).String(),
		time.Since(time.Unix(0, task.CreatedAt)).String(),
		task.MemoryMB,
		task.DiskMB,
	)
}

func taskState(task models.TaskState, c Config) string {
	switch state {
	case models.TaskStatePending:
		return "PENDING  "
	case models.TaskStateClaimed:
		return "CLAIMED  "
	case models.TaskStateRunning:
		return "RUNNING  "
	case models.TaskStateCompleted:
		return colorByTaskSuccess(task, c, "COMPLETED")
	case models.TaskStateResolving:
		return colorByTaskSuccess(task, c, "RESOLVING")
	default:
		return "INVALID"
	}
}

func colorbyTaskSuccess(task models.TaskState, c Config, format string, args ...interface{}) string {
	if task.Failed {
		return c.S.colorize(redColor, format, args...)
	} else {
		return c.S.colorize(greenColor, format, args...)
	}
}
