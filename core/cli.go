package core

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

type CLI struct {
	space *SpaceData
	data  string
}
type CLIOption func(*CLI)

func setupSpace() {
	SetupManager()
}

func NewCLI(options ...CLIOption) (*CLI, error) {
	setupSpace()

	c := new(CLI)
	for _, o := range options {
		o(c)
	}
	return c, nil
}

func LoginOption(u string, p string) CLIOption {

	return func(c *CLI) {
		c.space = NewSpaceData("", "")
	}
}

func DataOption(data string) CLIOption {
	return func(c *CLI) {
		c.data = data
	}
}

func NewSpaceData(token string, name string) *SpaceData {
	if token == "" {
		token = GetUUID()
	}

	timestamp := time.Now().Format("20060102150405")
	if name == "" {
		name = "app-" + timestamp
	}
	space := &SpaceData{Token: token, Name: name, Timestamp: timestamp}
	return space
}

type Config struct {
	Space *SpaceData `json:"space"`
}

func LocalSpaceData(path string) *SpaceData {

	raw, err := ioutil.ReadFile(path + "/ixkit.config.json")
	if nil != err {
		TraceMe.Error("Failed LocalSpaceData->%s", err)
		return nil
	}
	config := &Config{
		Space: &SpaceData{},
	}
	if err := json.Unmarshal(raw, &config); err != nil {
		TraceMe.Error("Failed LocalSpaceData->%s", err)
	}
	TraceMe.Debug("LocalSpaceData: %s", map2string(config))
	return config.Space

}

func SpaceOption(space *SpaceData) CLIOption {
	return func(c *CLI) {
		c.space = space
	}
}

func (c *CLI) Exec(cmd string, data string) {
	if cmd == "create" {
		c.create(data)
		return
	}
	if cmd == "maps" {
		c.maps(data)
		return
	}
	if cmd == "load" {
		c.load(data)
		return
	}
}

type SpaceData struct {
	Token     string `json:"token"`
	Name      string `json:"name"`
	Timestamp string `json:"timestamp,omitempty"`
	Desc      string `json:"desc"`
}

type RequestData struct {
	Space *SpaceData `json:"space"`
	Cmd   string     `json:"cmd"`
	Data  string     `json:"data"`
}

type ResponseData struct {
	Success bool       `json:"success"`
	Data    *EmbedData `json:"data"`
}
type EmbedData struct {
	Effect int `json:"effect"`
}

//@ service
func (c *CLI) maps(def string) {

	s := KitService{}
	success := func(resp *http.Response, body []byte) {
		TraceMe.Debug("success->%s,%s", resp, string(body))
		r := &ResponseData{
			Data: &EmbedData{},
		}
		err := json.Unmarshal(body, r)
		if nil == err {
			TraceMe.Info("Mapping success! effect %d\n typing command: ixkit load, pull full platforms source code.", r.Data.Effect)
		}

	}
	failure := func(err error) {
		TraceMe.Error("failure->%s", err)
	}

	data := &RequestData{Space: c.space, Cmd: "map", Data: def}

	j, err := json.Marshal(data)
	if err != nil {
		TraceMe.Error("err->%s", err)
	}

	s.maps(string(j), success, failure)
}

func map2string(def interface{}) string {

	bytes, _ := json.Marshal(def)
	jsonStr := string(bytes)
	return jsonStr
}

func isZipResponse(resp *http.Response, body []byte) bool {
	header := resp.Header
	ctype := map2string(header["Content-Type"])
	if strings.Contains(ctype, "application/json") {
		return false
	}
	if strings.Contains(ctype, "zip") {
		return true
	}
	return false
}

//@ service
func (c *CLI) create(def string) {
	tempFolder := GetTempFolder()
	output := tempFolder + "/" + c.space.Name + ".zip"
	prjPath := "./" + c.space.Name + "/"

	s := KitService{}
	success := func(resp *http.Response, body []byte) {
		TraceMe.Debug("create done %s	", resp)
		if !isZipResponse(resp, body) {
			b, err := ioutil.ReadFile(output) // just pass the file name
			if err == nil {
				fmt.Println(string(b))
			}
			return
		}
		UnZip(output, prjPath)
		os.Remove(output)
		defer os.Remove(tempFolder)
		TraceMe.Info("Create project %s success!\n typing command: cd %s, entry the folder, start mapping http call to service source code!", c.space.Name, c.space.Name)
	}
	failure := func(err error) {
		TraceMe.Error("Failure %s", err)
		defer os.Remove(tempFolder)
	}

	data := &RequestData{Space: c.space, Cmd: "pull", Data: def}
	j, err := json.Marshal(data)
	if err != nil {
		fmt.Println("err->%s", err)
	}

	s.create(string(j), output, success, failure)
}

//@ service
func (c *CLI) load(def string) {
	tempFolder := GetTempFolder()
	output := tempFolder + "/" + c.space.Name + ".zip"
	TraceMe.Debug("load %s", output)
	prjPath := "./"

	s := KitService{}
	success := func(resp *http.Response, body []byte) {
		if !isZipResponse(resp, body) {
			b, err := ioutil.ReadFile(output) // just pass the file name
			if err == nil {
				fmt.Println(string(b))
			}
			return
		}
		UnZip(output, prjPath)
		os.Remove(output)
		defer os.Remove(tempFolder)

		//TraceMe.Info("Load project success!")
	}
	failure := func(err error) {
		TraceMe.Error("Failure %s", err)
		defer os.Remove(tempFolder)
	}
	data := &RequestData{Space: c.space, Cmd: "load", Data: def}
	j, err := json.Marshal(data)
	if err != nil {
		TraceMe.Error("error->%s", err)
	}
	s.load(string(j), output, success, failure)
}
