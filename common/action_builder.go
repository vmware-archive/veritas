package common

import (
	"fmt"
	"strings"

	"github.com/cloudfoundry-incubator/receptor"
	"github.com/cloudfoundry-incubator/runtime-schema/models"
	"github.com/pivotal-cf-experimental/veritas/say"
)

type PreFabAction struct {
	Name          string
	ActionBuilder func() models.Action
}

func ModelEnvsToReceptorEnvs(in []models.EnvironmentVariable) []receptor.EnvironmentVariable {
	out := []receptor.EnvironmentVariable{}

	for _, env := range in {
		out = append(out, receptor.EnvironmentVariable{
			Name:  env.Name,
			Value: env.Value,
		})
	}

	return out
}

func BuildEnvs() []models.EnvironmentVariable {
	envs := say.Ask("Environment Variables (FOO=BAR;BAZ=WIBBLE)")
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

func BuildAction(description string, preFabActions []PreFabAction) models.Action {
	action, _ := buildActionWithDelete(description, preFabActions, false)
	return action
}

func buildActionWithDelete(description string, preFabActions []PreFabAction, allowDelete bool) (models.Action, bool) {
	choices := []string{
		"Done",
		"SerialAction",
		"ParallelAction",
		"DownloadAction",
		"RunAction",
		"UploadAction",
	}

	for _, preFabAction := range preFabActions {
		choices = append(choices, preFabAction.Name)
	}

	if allowDelete {
		choices = append(choices, "Delete Previous")
	}

	choice := say.Pick(description, choices)

	switch choice {
	case "Done":
		return nil, false
	case "SerialAction":
		return buildSerialAction(preFabActions), false
	case "ParallelAction":
		return buildParallelAction(preFabActions), false
	case "DownloadAction":
		return buildDownloadAction(), false
	case "RunAction":
		return buildRunAction(), false
	case "UploadAction":
		return buildUploadAction(), false
	case "Delete Previous":
		return nil, true
	default:
		for _, preFabAction := range preFabActions {
			if preFabAction.Name == choice {
				return preFabAction.ActionBuilder(), false
			}
		}
	}

	return nil, false
}

func buildActionCollection(description string, preFabActions []PreFabAction) []models.Action {
	actions := []models.Action{}
	for {
		action, deletePrev := buildActionWithDelete(description, preFabActions, len(actions) > 0)
		if deletePrev {
			actions = actions[0 : len(actions)-1]
			continue
		}
		if action == nil {
			break
		}
		actions = append(actions, action)
	}
	return actions
}
func buildSerialAction(preFabActions []PreFabAction) models.Action {
	actions := buildActionCollection("Build Series Action", preFabActions)
	if len(actions) == 0 {
		return nil
	}
	return &models.SerialAction{
		Actions: actions,
	}
}

func buildParallelAction(preFabActions []PreFabAction) models.Action {
	actions := buildActionCollection("Build Parallel Action", preFabActions)
	if len(actions) == 0 {
		return nil
	}
	return &models.ParallelAction{
		Actions: actions,
	}
}

func buildDownloadAction() models.Action {
	return &models.DownloadAction{
		From:     say.Ask("Download URL"),
		To:       say.AskWithDefault("Container Destination", "."),
		CacheKey: say.Ask("CacheKey"),
	}
}

func buildUploadAction() models.Action {
	return &models.UploadAction{
		From: say.Ask("Container Source"),
		To:   say.Ask("Upload URL"),
	}
}

func buildRunAction() models.Action {
	return &models.RunAction{
		Path: say.AskWithValidation("Command to run", func(response string) error {
			if strings.Contains(response, " ") {
				return fmt.Errorf("You cannot specify arguments to the command, that'll come next...")
			}
			return nil
		}),
		Args: strings.Split(say.Ask("Args (split by ';')"), ";"),
		Env:  BuildEnvs(),
	}
}
