package config_finder

import (
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/onsi/say"
)

func Autodetect(out io.Writer) error {
	jobs, err := ioutil.ReadDir("/var/vcap/jobs")
	if err != nil {
		return err
	}

	vitalsAddrs := []string{}
	gardenAddr := ""
	gardenNetwork := ""
	bbsEndpoint := ""

	debugRe := regexp.MustCompile(`debugAddr=(\d+.\d+.\d+.\d+:\d+)`)
	gardenTCPAddrRe := regexp.MustCompile(`gardenAddr=(\d+.\d+.\d+.\d+:\d+)`)
	gardenUnixAddrRe := regexp.MustCompile(`gardenAddr=([/\-\w+\.\d]+)`)
	bbsEndpointRe := regexp.MustCompile(`bbsAddress=([:/\-\w+\.\d]+)`)

	for _, job := range jobs {
		jobDir := filepath.Join("/var/vcap/jobs", job.Name(), "bin")
		ctls, err := ioutil.ReadDir(jobDir)
		if err != nil {
			return err
		}

		for _, ctl := range ctls {
			if ctl.IsDir() {
				continue
			}
			if strings.HasSuffix(ctl.Name(), "_ctl") {
				name := strings.TrimSuffix(ctl.Name(), "_ctl")
				path := filepath.Join(jobDir, ctl.Name())
				data, err := ioutil.ReadFile(path)
				if err != nil {
					return err
				}

				if debugRe.Match(data) {
					addr := string(debugRe.FindSubmatch(data)[1])
					vitalsAddrs = append(vitalsAddrs, fmt.Sprintf("%s:%s", name, addr))
				}

				if name == "rep" {
					if gardenTCPAddrRe.Match(data) {
						gardenAddr = string(gardenTCPAddrRe.FindSubmatch(data)[1])
						gardenNetwork = "tcp"
					} else if gardenUnixAddrRe.Match(data) {
						gardenAddr = string(gardenUnixAddrRe.FindSubmatch(data)[1])
						gardenNetwork = "unix"
					}

					if bbsEndpointRe.Match(data) {
						bbsEndpoint = string(bbsEndpointRe.FindSubmatch(data)[1])
					}
				}
			}
		}
	}

	if len(vitalsAddrs) > 0 {
		say.Fprintln(out, 0, "export VITALS_ADDRS=%s", strings.Join(vitalsAddrs, ","))
	}
	if gardenAddr != "" {
		say.Fprintln(out, 0, "export GARDEN_ADDR=%s", gardenAddr)
		say.Fprintln(out, 0, "export GARDEN_NETWORK=%s", gardenNetwork)
	}
	if bbsEndpoint != "" {
		say.Fprintln(out, 0, "export BBS_ENDPOINT=%s", bbsEndpoint)
	}

	return nil
}
