package chug

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/cloudfoundry-incubator/veritas/say"
	"github.com/pivotal-golang/lager"
	"github.com/pivotal-golang/lager/chug"
)

type JSFriendlyChugEntry struct {
	IsLager bool   `json:"is_lager"`
	Raw     string `json:"raw"`
	Log     struct {
		Timestamp int64      `json:"timestamp"`
		LogLevel  string     `json:"level"`
		Source    string     `json:"source"`
		Message   string     `json:"message"`
		Session   string     `json:"session"`
		Error     string     `json:"error,omitempty"`
		Trace     string     `json:"trace,omitempty"`
		Data      lager.Data `json:"data"`
	} `json:"log,omitempty"`
}

func NewJSFriendlyChugEntry(entry chug.Entry) JSFriendlyChugEntry {
	jsEntry := JSFriendlyChugEntry{}
	jsEntry.IsLager = entry.IsLager
	jsEntry.Raw = string(entry.Raw)
	if entry.IsLager {
		jsEntry.Log.Timestamp = entry.Log.Timestamp.UnixNano()
		jsEntry.Log.LogLevel = lagerLogLevel(entry.Log.LogLevel)
		jsEntry.Log.Source = entry.Log.Source
		jsEntry.Log.Message = entry.Log.Message
		jsEntry.Log.Session = entry.Log.Session
		if entry.Log.Error != nil {
			jsEntry.Log.Error = entry.Log.Error.Error()
		}
		jsEntry.Log.Trace = entry.Log.Trace
		jsEntry.Log.Data = entry.Log.Data
	}
	return jsEntry
}

func lagerLogLevel(level lager.LogLevel) string {
	switch level {
	case lager.DEBUG:
		return "DEBUG"
	case lager.INFO:
		return "INFO"
	case lager.ERROR:
		return "ERROR"
	case lager.FATAL:
		return "FATAL"
	}
	return ""
}

func ServeLogs(addr string, dev bool, src io.Reader) error {
	out := make(chan chug.Entry)
	go chug.Chug(src, out)
	entries := []JSFriendlyChugEntry{}
	for entry := range out {
		if isEmptyInigoLog(entry) {
			continue
		}
		jsEntry := NewJSFriendlyChugEntry(entry)
		entries = append(entries, jsEntry)
	}

	http.HandleFunc("/assets/", AssetServer(dev))

	http.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(entries)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/assets/index.html", http.StatusTemporaryRedirect)
	})

	listener, err := net.Listen("tcp", addr)
	say.Println(0, say.Green("Serving up on http://127.0.0.1:%d", listener.Addr().(*net.TCPAddr).Port))
	if err != nil {
		return err
	}

	return http.Serve(listener, nil)
}

func AssetServer(dev bool) http.HandlerFunc {
	if dev {
		regenerateAssetsTrigger := make(chan bool)
		go func() {
			for {
				<-regenerateAssetsTrigger
				regenerateAssets()
			}
		}()
		return func(w http.ResponseWriter, r *http.Request) {
			fname := path.Base(r.URL.Path)
			f, err := os.Open(filepath.Join("./chug/assets", fname))
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			io.Copy(w, f)
			regenerateAssetsTrigger <- true
		}
	} else {
		return func(w http.ResponseWriter, r *http.Request) {
			fname := path.Base(r.URL.Path)
			data, hasAsset := assets[fname]

			if hasAsset {
				stream := base64.NewDecoder(base64.StdEncoding, strings.NewReader(data))
				io.Copy(w, stream)
			} else {
				w.WriteHeader(http.StatusNotFound)
			}
		}
	}
}

func regenerateAssets() {
	dir, err := ioutil.ReadDir("./chug/assets")
	if err != nil {
		panic(err)
	}

	out, err := os.Create("./chug/assets.go")
	if err != nil {
		panic(err)
	}

	out.WriteString(`package chug

var assets = map[string]string{
`)

	for _, info := range dir {
		f, err := os.Open(filepath.Join("./chug/assets", info.Name()))
		if err != nil {
			panic(err)
		}
		out.WriteString(fmt.Sprintf(`"%s": "`, info.Name()))
		encoder := base64.NewEncoder(base64.StdEncoding, out)
		io.Copy(encoder, f)
		encoder.Close()
		out.WriteString("\",\n")
	}
	out.WriteString("}")
	out.Close()

	err = exec.Command("go", "fmt", "./chug/assets.go").Run()
	if err != nil {
		panic(err)
	}
}
