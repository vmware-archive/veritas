package submit_lrp

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"github.com/tedsuo/rata"

	"github.com/cloudfoundry-incubator/runtime-schema/routes"

	"github.com/cloudfoundry-incubator/runtime-schema/bbs"
	"github.com/cloudfoundry-incubator/runtime-schema/models"
	"github.com/cloudfoundry/gunk/timeprovider"
	"github.com/cloudfoundry/gunk/workpool"
	"github.com/cloudfoundry/storeadapter/etcdstoreadapter"
	"github.com/pivotal-cf-experimental/veritas/common"
	"github.com/pivotal-cf-experimental/veritas/say"
	"github.com/pivotal-golang/lager"
)

func SubmitLRP(cluster []string, f io.Reader) error {
	adapter := etcdstoreadapter.NewETCDStoreAdapter(cluster, workpool.NewWorkPool(10))
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
	desiredLRP.CPUWeight = uint(say.AskForIntegerWithDefault("CPUWeight", 100))
	desiredLRP.EnvironmentVariables = common.BuildEnvs()
	desiredLRP.Routes = []string{say.AskWithDefault("Route", desiredLRP.ProcessGuid+".10.244.0.34.xip.io")}
	ports := say.AskWithDefault("Ports to open (comma separated)", "8080")
	desiredLRP.Ports = []uint32{}
	for _, portString := range strings.Split(ports, ",") {
		port, err := strconv.Atoi(portString)
		if err != nil {
			say.Println(0, say.Red("Ignoring invalid port %s", portString))
			continue
		}
		desiredLRP.Ports = append(desiredLRP.Ports, uint32(port))
	}
	desiredLRP.LogGuid = desiredLRP.ProcessGuid
	desiredLRP.LogSource = "VRT"

	requestGenerator := rata.NewRequestGenerator("http://file_server.service.dc1.consul:8080", routes.FileServerRoutes)
	circusDownloadRequest, _ := requestGenerator.CreateRequest(routes.FS_STATIC, nil, nil)

	setup := common.BuildAction("Build Setup Action", []common.PreFabAction{
		common.PreFabAction{
			Name: "Download Spy",
			ActionBuilder: func() models.Action {
				return &models.DownloadAction{
					From: circusDownloadRequest.URL.String(),
					To:   "/tmp/circus",
				}
			},
		},
	})

	if setup != nil {
		desiredLRP.Setup = setup
	}

	desiredLRP.Action = common.BuildAction("Build Action", nil)

	monitor := common.BuildAction("Build Monitor Action", []common.PreFabAction{
		common.PreFabAction{
			Name: "Run Spy with Port Check on 8080",
			ActionBuilder: func() models.Action {
				return &models.RunAction{
					Path: "/tmp/circus/spy",
					Args: []string{"-addr=:" + say.AskWithDefault("Port", "8080")},
				}
			},
		},
	})

	if monitor != nil {
		desiredLRP.Monitor = monitor
	}

	return desiredLRP
}
