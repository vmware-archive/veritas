package submit_lrp

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"time"

	RepRoutes "github.com/cloudfoundry-incubator/rep/routes"
	"github.com/cloudfoundry-incubator/runtime-schema/bbs"
	"github.com/cloudfoundry-incubator/runtime-schema/models"
	SchemaRouter "github.com/cloudfoundry-incubator/runtime-schema/router"
	"github.com/cloudfoundry-incubator/veritas/say"
	"github.com/cloudfoundry/gunk/timeprovider"
	"github.com/cloudfoundry/gunk/urljoiner"
	"github.com/cloudfoundry/storeadapter/etcdstoreadapter"
	"github.com/cloudfoundry/storeadapter/workerpool"
	"github.com/pivotal-golang/lager"
	"github.com/tedsuo/rata"
)

func SubmitLRP(cluster []string, f io.Reader) error {
	adapter := etcdstoreadapter.NewETCDStoreAdapter(cluster, workerpool.NewWorkerPool(10))
	err := adapter.Connect()
	if err != nil {
		return err
	}

	store := bbs.NewVeritasBBS(adapter, timeprovider.NewTimeProvider(), lager.NewLogger("veritas"))

	var desiredLRP models.DesiredLRP

	if f != nil {
		decoder := json.NewDecoder(f)
		err := decoder.Decode(&desiredLRP)
		if err != nil {
			return err
		}
	} else {
		desiredLRP = interactivelyBuildDesiredLRP()

		asJson, err := json.Marshal(desiredLRP)
		if err != nil {
			return err
		}

		filename := fmt.Sprintf("desired_lrp_%s.json", desiredLRP.ProcessGuid)
		err = ioutil.WriteFile(filename, asJson, 0666)
		if err != nil {
			return err
		}
		say.Println(0, say.Green("Save DesiredLPR to %s", filename))
	}

	say.Println(0, say.Green("Desiring %s", desiredLRP.ProcessGuid))
	return store.DesireLRP(desiredLRP)
}

func interactivelyBuildDesiredLRP() models.DesiredLRP {
	desiredLRP := models.DesiredLRP{}
	desiredLRP.ProcessGuid = say.AskWithDefault("ProcessGuid", fmt.Sprintf("%d", time.Now().Unix()))
	desiredLRP.Domain = say.AskWithDefault("Domain", "veritas")
	desiredLRP.Instances = say.AskForIntegerWithDefault("Instances", 1)
	desiredLRP.Stack = say.AskWithDefault("Stack", "lucid64")
	desiredLRP.MemoryMB = say.AskForIntegerWithDefault("MemoryMB", 256)
	desiredLRP.DiskMB = say.AskForIntegerWithDefault("DiskMB", 256)
	desiredLRP.Routes = []string{say.AskWithDefault("Route", desiredLRP.ProcessGuid+".10.244.0.34.xip.io")}
	desiredLRP.Ports = []models.PortMapping{
		{ContainerPort: 8080},
	}
	desiredLRP.Log = models.LogConfig{
		Guid:       desiredLRP.ProcessGuid,
		SourceName: "VRT",
	}
	desiredLRP.Actions = interactivelyBuildActions(desiredLRP.ProcessGuid)

	return desiredLRP
}

func interactivelyBuildActions(processGuid string) []models.ExecutorAction {
	actions := []models.ExecutorAction{}
	for {
		choice := say.Pick("Add an action", []string{
			"Done",
			"DownloadAction",
			"Health-Monitored RunAction",
		})

		switch choice {
		case "Done":
			return actions
		case "DownloadAction":
			actions = append(actions, interactivelyBuildDownloadAction())
		case "Health-Monitored RunAction":
			staticRoute, _ := SchemaRouter.NewFileServerRoutes().RouteForHandler(SchemaRouter.FS_STATIC)
			circusURL := urljoiner.Join("PLACEHOLDER_FILESERVER_URL", staticRoute.Path, "linux-circus/linux-circus.tgz")
			actions = append(actions, models.ExecutorAction{
				models.DownloadAction{
					From:    circusURL,
					To:      "/tmp/circus",
					Extract: true,
				},
			})
			actions = append(actions, interactivelyBuildHealthMonitoredRunAction(processGuid))
		}
	}
}

func interactivelyBuildDownloadAction() models.ExecutorAction {
	return models.ExecutorAction{
		models.DownloadAction{
			From:    say.Ask("Download URL"),
			To:      say.AskWithDefault("Container Destination", "."),
			Extract: say.AskForBoolWithDefault("Extract", true),
		},
	}
}

func interactivelyBuildHealthMonitoredRunAction(processGuid string) models.ExecutorAction {
	repRequests := rata.NewRequestGenerator(
		"http://127.0.0.1:20515",
		RepRoutes.Routes,
	)

	healthyHook, _ := repRequests.CreateRequest(
		RepRoutes.LRPRunning,
		rata.Params{
			"process_guid":  processGuid,
			"index":         "PLACEHOLDER_INSTANCE_INDEX",
			"instance_guid": "PLACEHOLDER_INSTANCE_GUID",
		},
		nil,
	)

	return models.Parallel(
		models.ExecutorAction{
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
		},
		models.ExecutorAction{
			models.MonitorAction{
				Action: models.ExecutorAction{
					models.RunAction{
						Path: "/tmp/circus/spy",
						Args: []string{"-addr=:8080"},
					},
				},
				HealthyThreshold:   1,
				UnhealthyThreshold: 1,
				HealthyHook: models.HealthRequest{
					Method: healthyHook.Method,
					URL:    healthyHook.URL.String(),
				},
			},
		},
	)
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
