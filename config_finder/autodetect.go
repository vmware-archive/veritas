package config_finder

import (
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/cloudfoundry-incubator/veritas/say"
)

func Autodetect(out io.Writer) error {
	jobs, err := ioutil.ReadDir("/var/vcap/jobs")
	if err != nil {
		return err
	}

	vitalsAddrs := []string{}
	executorAddr := ""
	wardenAddr := ""
	etcdCluster := ""

	debugRe := regexp.MustCompile(`debugAddr=(\d+.\d+.\d+.\d+:\d+)`)
	etcdRe := regexp.MustCompile(`etcdCluster=\"(.+)\"`)
	executorRe := regexp.MustCompile(`listenAddr=(\d+.\d+.\d+.\d+:\d+)`)
	wardenRe := regexp.MustCompile(`wardenAddr=(\d+.\d+.\d+.\d+:\d+)`)

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

				if etcdRe.Match(data) {
					etcdCluster = string(etcdRe.FindSubmatch(data)[1])
				}

				if name == "executor" && executorRe.Match(data) {
					executorAddr = "http://" + string(executorRe.FindSubmatch(data)[1])
				}

				if name == "executor" && wardenRe.Match(data) {
					wardenAddr = string(wardenRe.FindSubmatch(data)[1])
				}
			}
		}
	}

	if len(vitalsAddrs) > 0 {
		say.Fprintln(out, 0, "export VITALS_ADDRS=%s", strings.Join(vitalsAddrs, ","))
	}
	if executorAddr != "" {
		say.Fprintln(out, 0, "export EXECUTOR_ADDR=%s", executorAddr)
	}
	if wardenAddr != "" {
		say.Fprintln(out, 0, "export WARDEN_ADDR=%s", wardenAddr)
	}
	if etcdCluster != "" {
		say.Fprintln(out, 0, "export ETCD_CLUSTER=%s", etcdCluster)
	}

	return nil
}
