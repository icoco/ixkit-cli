package main

import (
	"encoding/json"
	"fmt"
	//"github.com/icoco/ixkit-cli/vendors/console"
	"os"

	//argparse "github.com/akamensky/argparse"
	argparse "github.com/icoco/ixkit-cli/vendors/akamensky/argparse"
	//bat "github.com/astaxie/bat"
	bat "github.com/icoco/ixkit-cli/vendors/astaxie/bat"
	//cli "github.com/icoco/ixkit"
	cli "github.com/icoco/ixkit-cli/core"
)

type CmdOptions struct {
	Ver  string
	Help string

	Init string
	Maps string
	Load string

	Desc string

	Space string
	User  string
}

func (c *CmdOptions) isCurl() bool {
	call := func() {
		outUsage()
	}
	return bat.ValidateCURLArg(call)
}

func (c *CmdOptions) isMurl() bool {

	if len(os.Args) == 1 {
		return false
	}
	return "" != c.Init || "" != c.Maps || "" != c.Load

}

func (c *CmdOptions) isMapping() bool {

	if len(os.Args) == 1 {
		return false
	}
	return "" != c.Maps

}

var (
	parser     *argparse.Parser
	kitVersion string = "ixkit version:0.1.0"
)

func parseCommand() *CmdOptions {
	// Create new parser object
	parser = argparse.NewParser("ixkit", usageinfo)

	// init command
	var init string = ""
	initCmd := parser.NewCommand("start", "Start new application, use -n|--name <name> specifiy the application name.")

	// load command
	var load string = ""
	loadCmd := parser.NewCommand("load", "Load full mobile platforms source code of the application.")
	//	var platformSelector *string = loadCmd.Selector("p", "platform", []string{"*", "swift", "objectc", "android"}, &argparse.Options{NullAble: true, Default4Exist: "!", Help: "description"})

	maps := parser.String("m", "map", &argparse.Options{Required: false, ValueDesc: "Service.Method name", Default4Exist: "!", Help: "Mapping http call to Method of Web Service class in application, value format: \"Service.Method\" \n eg: ixkit http://ixkit.com/api/app -m AppService.list"})

	nameFlag := parser.String("n", "name", &argparse.Options{Required: false, ValueDesc: "name", Help: "Specifiy name"})
	descFlag := parser.String("d", "desc", &argparse.Options{Required: false, Default4Exist: "", ValueDesc: "Description", Help: "Describe the item,optional"})

	help := parser.String("h", "help", &argparse.Options{Required: false, Help: "show usage"})
	ver := parser.String("v", "ver", &argparse.Options{Required: false, NullAble: true, Default4Exist: "!", Help: "Show version"})
	trace := parser.String("", "trace", &argparse.Options{Required: false, NullAble: true, Default4Exist: "!", Help: "Debug mode"})

	//space := parser.String("s", "space", &argparse.Options{Required: false, Help: "Space name"})
	//user := parser.String("u", "user", &argparse.Options{Required: false, Help: "User name"})

	// Parse input
	err := parser.Parse(os.Args)
	if err != nil {
		// In case of error print error and print usage
		// This can also be done by passing -h or --help flags
		outUsage()
	}
	setupTracer(*trace)

	if "" == *descFlag {
		d := ""
		descFlag = &d
	}

	if initCmd.Happened() {
		if "" == *nameFlag {
			outUsage()
		}

		init = *nameFlag
		cli.TraceMe.Debug("initCmd", init, *descFlag)

	}
	if loadCmd.Happened() {
		load = *nameFlag
		if "" == load {
			load = "!"

		}
		cli.TraceMe.Debug("loadCmd", load)
	}
	o := &CmdOptions{Help: *help, Ver: *ver, Init: init, Maps: *maps, Load: load, Desc: *descFlag}
	if "" != o.Ver {
		fmt.Println(kitVersion)
		os.Exit(0)
	}

	if "!" == o.Maps || "!" == o.Init || "!" == o.Help {
		outUsage()
	}
	return o
}
func outUsage() {

	fmt.Print(parser.Usage(nil))

	onEcho := func(p *PromptTyping, value string) *PromptTyping {
		fmt.Print(usageinfo2)
		return nil
	}
	tip := "🚦ixkit is standing on the bat which implmeted cURL features, press enter key to continue check the full detail usage."
	p := NewPromptTyping(tip, onEcho)
	_ = readTypingP(p)

	os.Exit(0)
}

func main() {

	co := parseCommand()
	cli.TraceMe.Debug("cli args %s", map2string(co))
	if co.isMurl() {
		if !co.isMapping() {
			murl(co)
			return
		}

	}

	if co.isCurl() {
		if co.isMapping() {
			murl(co)
			return
		}
		TTYIndictiator().start()
		defer TTYIndictiator().stop()
		curl()
		TTYIndictiator().stop()
		return
	}

	outUsage()
}

func murl(co *CmdOptions) {
	options := make(map[string]interface{})
	//@case
	if "" != co.Init {
		// init application
		options["cmd"] = co
		initApp(co)
		return
	}
	//@case
	if "" != co.Load {
		// pull application
		TTYIndictiator().start()
		defer TTYIndictiator().stop()
		options["cmd"] = co
		var space *cli.SpaceData
		if "!" == co.Load {
			dir, _ := os.Getwd()
			space = cli.LocalSpaceData(dir)
		} else {
			space = cli.NewSpaceData("", co.Load)
		}

		if nil == space {
			return
		}

		loadApp(space, options)
		TTYIndictiator().stop()
		cli.TraceMe.Info("Load project success!")
		return
	}
	//@case
	if "" != co.Maps {
		// import HTTP call to service class
		TTYIndictiator().start()
		defer TTYIndictiator().stop()
		def := curl()
		def["map"] = append(def["map"], co.Maps)
		if "" != co.Desc {
			def["desc"] = append(def["desc"], co.Desc)
		}
		api := map2ApiDef(def)

		jsonStr := map2string(api)

		services := make([]string, 0)
		services = append(services, jsonStr)

		dir, _ := os.Getwd()
		space := cli.LocalSpaceData(dir)
		importService(space, services)

		return
	}
}

func setupTracer(traceLevel string) {

	if traceLevel == "debug" {
		cli.TraceMe.SetDebugMode(true)
	}

}

//////////////////////////////
func curl() map[string][]interface{} {
	_, def := bat.Main()
	return def
}

func getPasswordRoutine() *PromptTyping {
	onEcho := func(p *PromptTyping, value string) *PromptTyping {
		l := len(value)
		if l < 5 {
			fmt.Println("password is invalidated")
			return p
		}
		p.push("password", value)
		return nil
	}
	tip := "❓🔒please type your password:"
	p := NewPromptTyping(tip, onEcho)
	return p
}

func getUserEmailRoutine() *PromptTyping {
	onEcho := func(p *PromptTyping, value string) *PromptTyping {
		l := len(value)
		if l < 5 {
			fmt.Println("email is invalidated")
			return p
		}
		p.push("email", value)
		return getPasswordRoutine()
	}
	tip := "❓📧please type your email:"
	p := NewPromptTyping(tip, onEcho)
	return p
}

func getRgisterRoutine() *PromptTyping {
	onEcho := func(p *PromptTyping, value string) *PromptTyping {
		l := len(value)
		if l < 5 {
			fmt.Println("user name should at least 5 length")
			return p
		}
		p.push("name", value)
		return getUserEmailRoutine()
	}
	tip := "❓Type register user name"
	p := NewPromptTyping(tip, onEcho)
	return p
}

func getLoginRoutine() *PromptTyping {
	onEcho := func(p *PromptTyping, value string) *PromptTyping {
		l := len(value)
		if l < 5 {
			fmt.Println("user name should at least 5 length")
			return p
		}
		p.push("name", value)
		return getPasswordRoutine()
	}
	tip := "❓Type user name"
	p := NewPromptTyping(tip, onEcho)
	return p
}

type UserAction int

const (
	_ UserAction = iota
	Anonymous
	Register
	Login
)

///////////////////////////////
func initApp(co *CmdOptions) {

	appName := co.Init
	space := cli.NewSpaceData("", appName)
	if "" != co.Desc {
		space.Desc = co.Desc
	}
	TTYIndictiator().start()
	createApp(space, nil)
	TTYIndictiator().stop()
	if 1 > 0 {

		return
	}

	if "" == co.User {
		onEcho := func(p *PromptTyping, value string) *PromptTyping {
			//@case Anonymous
			if "N" == value || "n" == value {

				p.push("action", "Anonymous")
				return nil
			}
			//@case register
			if "" == value {
				p.push("action", "Register")
				return getRgisterRoutine()
			}
			//@case login
			p.push("action", "Login")
			return getLoginRoutine()
		}
		tip := "🛠 Create project: recommend login as registed user, \nif you had already registered user then type your user name, \nif want register new user then just press enter key will start register flow, \nor if you want anonymous then type 'N' or 'n',then press enter key."
		p := NewPromptTyping(tip, onEcho)
		u := readTypingP(p)
		fmt.Println("go last input " + u)
		values := make(map[string]string)
		values = p.getValues(values)
		fmt.Println(values)

		action := values["action"]
		if "Anonymous" == action {
			appName := co.Init
			space := cli.NewSpaceData("", appName)
			createApp(space, nil)
		}

	}

}

func map2string(def interface{}) string {

	bytes, _ := json.Marshal(def)
	jsonStr := string(bytes)
	return jsonStr
}

func map2ApiDef(def map[string][]interface{}) *cli.ApiDef {
	api := cli.NewApiDef()
	api.Name = (def["map"][0]).(string)
	api.Url = (def["url"][0]).(string)
	api.Method = def["method"][0].(string)
	api.Response = def["response"][0].(string)
	api.Parameters = def["parameters"][0].(map[string]interface{})
	api.Headers = def["headers"][0].(map[string]interface{})
	api.Body = def["body"][0].(map[string]interface{})

	if nil != def["desc"] && len(def["desc"]) >= 1 {
		desc := def["desc"][0].(string)
		api.Put("desc", desc)
	}

	return api
}

func createApp(space *cli.SpaceData, options interface{}) {

	option := cli.SpaceOption(space)
	c, _ := cli.NewCLI(option)

	data := make(map[string]interface{})
	data["space"] = space
	jsonStr := map2string(data)

	c.Exec("create", jsonStr)
}

func importService(space *cli.SpaceData, services []string) {

	jsonStr := map2string(services)

	option := cli.SpaceOption(space)
	c, _ := cli.NewCLI(option)

	c.Exec("maps", jsonStr)
}

func loadApp(space *cli.SpaceData, options interface{}) {
	option := cli.SpaceOption(space)
	c, _ := cli.NewCLI(option)

	data := make(map[string]interface{})
	data["space"] = space
	jsonStr := map2string(data)

	c.Exec("load", jsonStr)
}

var usageinfo string = `
	🐙ixkit cli is a Go implemented CLI cURL-like tool for automate build mobile client SDK application, inspired Httpie and bat, but ixkit do more steps, it can mapping the http call to native mobile application source code directly, iOS(Swift|Objective C), Android and more.

	cURL -> Request + Response => mobile app source code

	Purpose:
	🚀Base http call automate generate native App client source code for all mobile platforms that can invoke & consume the http web service easily.
	⚙️ Help you fast build mobile client SDK for your RESTful APIs!
	💡more help information of ixkit cli,please refer to https://github.com/icoco/ixkit
`

var usageinfo2 string = `🐙ixkit transfer command to bat that is a Go implemented CLI cURL-like tool for humans, then mapping the http request and response as definition to service class source code of native mobile client application. 

Full Usage:

	ixkit [command] [flags] [METHOD] URL [ITEM [ITEM]] 

Command:
  start  Start new application,use -n|--name <name> specifiy the application name.
  load   Load full mobile platforms source code of the application.
flags:
  -m  --map   Mapping http call to Method of Web Service class in application,
 		eg: ixkit http://ixkit.com/api/app -m AppService.list
  -n  --name  Specifiy the application name
  -d  --desc  Describe the item,optional

curl flags:
  -a, -auth=USER[:PASS]       Pass a username:password pair as the argument
  -b, -bench=false            Sends bench requests to URL
  -b.N=1000                   Number of requests to run
  -b.C=100                    Number of requests to run concurrently
  -body=""                    Send RAW data as body
  -f, -form=false             Submitting the data as a form
  -j, -json=true              Send the data in a JSON object
  -p, -pretty=true            Print Json Pretty Format
  -i, -insecure=false         Allow connections to SSL sites without certs
  -proxy=PROXY_URL            Proxy with host and port

METHOD:
  defaults to either GET (if there is no request data) or POST (with request data).

URL:
  The only information needed to perform a request is a URL. The default scheme is http://,
  which can be omitted from the argument; example.org works just fine.

ITEM:
  Can be any of:
    Query string   key=value
    Header         key:value
    Post data      key=value
    File upload    key@/path/file

Full Example:
	ixkit ixkit.com/api/app
	
	ixkit start -n myapp -d "Fast build mobile client sdk"
	cd myapp
	ixkit ixkit.com/api/app -m AppService.list 
	ixkit POST ixkit.com/api/app name="Demo App" --map AppService.create --desc "Create Action"
	
	ixkit load

	💡more help information of ixkit cli,please refer to https://github.com/icoco/ixkit
more help information of cURL, please refer to https://github.com/astaxie/bat
`
