package models

import (
	"sort"

	"github.com/cloudfoundry-incubator/runtime-schema/models"
)

type VeritasTasks map[string][]models.Task

func (t *VeritasTasks) OrderedTaskTypes() []string {
	taskTypes := []string{}
	for key := range t {
		taskTypes := append(taskTypes, key)
	}
	sort.Strings(taskTypes)
	return taskTypes
}

func (t *VeritasTasks) SortedTasksForType(taskType string) []models.Task {
	tasks := t[taskType]
	sort.Sort(TasksByUpdatedAt(tasks))
	return tasks
}

type TasksByUpdatedAt []models.Task

func (a TasksByUpdatedAt) Len() int           { return len(a) }
func (a TasksByUpdatedAt) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a TasksByUpdatedAt) Less(i, j int) bool { return a[i].UpdatedAt.Before(a[j].UpdatedAt) }
