package vitals

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/pivotal-cf-experimental/veritas/say"
)

func Vitals(vitalsAddrs map[string]string, out io.Writer) error {
	http.DefaultClient.Timeout = time.Second

	components := []string{}

	for component := range vitalsAddrs {
		components = append(components, component)
	}

	sort.Strings(components)

	say.Println(0, "Vitals on %s", time.Now())
	for _, component := range components {
		dumpVitals(component, vitalsAddrs[component], out)
	}

	return nil
}

func dumpVitals(component string, addr string, out io.Writer) {
	response, err := http.Get("http://" + addr + "/debug/pprof/")
	if err != nil {
		say.Println(0, say.Red("%s: %s"), component, err.Error())
		return
	}
	if response.StatusCode != http.StatusOK {
		say.Println(0, say.Red("%s: %d"), component, response.StatusCode)
		return
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		say.Println(0, say.Red("%s: %d"), component, err.Error())
		return
	}

	s := string(body)
	report := []string{}
	rows := strings.Split(s, "<tr>")[1:]
	for _, row := range rows {
		columns := strings.Split(row, "<td>")
		value, _ := strconv.Atoi(strings.Split(columns[0], ">")[1])
		name := strings.Split(columns[1], ">")[1]
		name = name[:len(name)-3]
		if value > 1000 {
			report = append(report, say.Red("%20s", fmt.Sprintf("%d %s", value, name)))
		} else {
			report = append(report, fmt.Sprintf("%20s", fmt.Sprintf("%d %s", value, name)))
		}
	}

	say.Println(0, "%s: %s %s", say.Green("%25s", component), string(strings.Join(report, " ")), addr)
}
