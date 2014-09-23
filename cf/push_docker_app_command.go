package cf

import (
	"encoding/json"
	"flag"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/cloudfoundry-incubator/veritas/common"
	"github.com/cloudfoundry-incubator/veritas/say"
)

type AppRequest struct {
	Name        string `json:"name,omitempty"`
	SpaceGuid   string `json:"space_guid,omitempty"`
	DockerImage string `json:"docker_image,omitempty"`
	Command     string `json:"command,omitempty"`
	Memory      int    `json:"memory"`
	DiskQuota   int    `json:"disk_quota"`
	Instances   int    `json:"instances"`
}

func PushDockerAppCommand() common.Command {
	var (
		appName     string
		space       string
		dockerImage string
		command     string
		domain      string

		memory    int
		disk      int
		instances int
	)

	flagSet := flag.NewFlagSet("push-docker-app", flag.ExitOnError)
	flagSet.StringVar(&appName, "appName", "", "app name (required)")
	flagSet.StringVar(&space, "space", "", "space (required)")
	flagSet.StringVar(&dockerImage, "dockerImage", "", "docker image (required)")
	flagSet.StringVar(&command, "command", "", "start command")
	flagSet.StringVar(&domain, "domain", "", "route domain (required - e.g. ketchup.cf-app.com)")

	flagSet.IntVar(&memory, "memory", 128, "memory limit (MB)")
	flagSet.IntVar(&disk, "disk", 1024, "disk limit (MB)")
	flagSet.IntVar(&instances, "instances", 1, "instances (n)")

	return common.Command{
		Name:        "push-docker-app",
		Description: " - Push a docker app to CC",
		FlagSet:     flagSet,
		Run: func(args []string) {
			validate(appName, "You must specify -appName")
			validate(space, "You must provide a -space")
			validate(dockerImage, "You must specify a -dockerImage")
			validate(domain, "You must specify a -domain")

			spaceGuid := getSpaceGuid(space)
			say.Println(0, "Your space guid is: %s", say.Green(spaceGuid))

			app := AppRequest{
				Name:        appName,
				SpaceGuid:   spaceGuid,
				DockerImage: dockerImage,
				Command:     command,
				Memory:      memory,
				DiskQuota:   disk,
				Instances:   instances,
			}
			encodedApp, err := json.Marshal(app)
			common.ExitIfError("Failed to build App JSON", err)

			CF("curl", "/v2/apps", "-X", "POST", "-d", string(encodedApp))
			CF("set-env", appName, "CF_DIEGO_BETA", "true")
			CF("set-env", appName, "CF_DIEGO_RUN_BETA", "true")
			CF("create-route", space, domain, "-n", appName)
			CF("map-route", appName, domain, "-n", appName)
			say.Println(0, "Your docker app is ready -- just:\n %s", say.Green("cf start %s", appName))
		},
	}
}

func validate(value string, errorMessage string) {
	if value == "" {
		say.Fprintln(os.Stderr, 0, say.Red(errorMessage))
		os.Exit(1)
	}
}

func CF(args ...string) {
	say.Println(0, say.Green("cf %s", strings.Join(args, " ")))
	cf := exec.Command("cf", args...)
	cf.Stdout = os.Stdout
	cf.Stderr = os.Stderr
	err := cf.Run()
	common.ExitIfError("CF Failed", err)
}

func getSpaceGuid(space string) string {
	re := regexp.MustCompile(`"guid": "([a-f0-9-]+)",`)
	cf := exec.Command("cf", "space", space)
	cf.Env = append(os.Environ(), "CF_TRACE=true")
	output, err := cf.CombinedOutput()
	common.ExitIfError("Fetchign space guid failed", err)
	return string(re.FindSubmatch(output)[1])
}
