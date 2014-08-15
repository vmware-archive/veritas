package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/cloudfoundry-incubator/veritas/say"
)

func PushDockerAppCommand() Command {
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
	flagSet.StringVar(&command, "command", "", "start command (required)")
	flagSet.StringVar(&domain, "domain", "", "route domain (required - e.g. ketchup.cf-app.com)")

	flagSet.IntVar(&memory, "memory", 128, "memory limit (MB)")
	flagSet.IntVar(&disk, "disk", 1024, "disk limit (MB)")
	flagSet.IntVar(&instances, "instances", 1, "instances (n)")

	return Command{
		Name:        "push-docker-app",
		Description: " - Push a docker app to CC",
		FlagSet:     flagSet,
		Run: func(args []string) {
			validate(appName, "You must specify -appName")
			validate(space, "You must provide a -space")
			validate(dockerImage, "You must specify a -dockerImage")
			validate(command, "You must specify a -command")
			validate(domain, "You must specify a -domain")

			spaceGuid := getSpaceGuid(space)
			say.Println(0, "Your space guid is: %s", say.Green(spaceGuid))

			CF("curl", "/v2/apps", "-X", "POST", "-d", fmt.Sprintf(`{
			    "name":"%s",
			    "space_guid":"%s",
			    "docker_image":"%s",
			    "command":"%s",
			    "memory":%d,
			    "disk_quota":%d,
			    "instances":%d
			   }`, appName, spaceGuid, dockerImage, command, memory, disk, instances))
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
	ExitIfError("CF Failed", err)
}

func getSpaceGuid(space string) string {
	re := regexp.MustCompile(`"guid": "([a-f0-9-]+)",`)
	cf := exec.Command("cf", "space", space)
	cf.Env = append(os.Environ(), "CF_TRACE=true")
	output, err := cf.CombinedOutput()
	ExitIfError("Fetchign space guid failed", err)
	return string(re.FindSubmatch(output)[1])
}
