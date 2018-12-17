# ixkit-cli

Usage: ixkit <Command> [-h|--help] [-m|--map "<value>"(Service.Method name)]
             [-n|--name "<value>"(name)] [-d|--desc "<value>"(Description)]
             [-v|--ver] [--trace]

             
	🐙ixkit cli is a Go implemented CLI cURL-like tool for automate build mobile
             client SDK application, inspired Httpie and bat, but ixkit do more
             steps, it can mapping the http call to native mobile application
             source code directly, iOS(Swift|Objective C), Android and
             more.

	cURL -> Request + Response => mobile app source code

	Purpose:
	🚀Base http call automate generate native App client source code for all
             mobile platforms that can invoke & consume the http web service
             easily.
	⚙️ Help you fast build mobile client SDK for your RESTful APIs!
	💡more help information of ixkit cli,please refer to
             https://github.com/icoco/ixkit


Commands:

  start  Start new application, use -n|--name <name> specifiy the application
          name.
  load   Load full mobile platforms source code of the application.

Arguments:

  -h  --help   Print help information
  -m  --map    Mapping http call to Method of Web Service class in application,
               value format: "Service.Method" 
 eg: ixkit http://ixkit.com/api/app -m AppService.list
  -n  --name   Specifiy name
  -d  --desc   Describe the item,optional
  -v  --ver    Show version
      --trace  Debug mode

🚦ixkit is standing on the bat which implmeted cURL features, press enter key to continue check the full detail usage.

🐙ixkit transfer command to bat that is a Go implemented CLI cURL-like tool for humans, then mapping the http request and response as definition to service class source code of native mobile client application. 

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


💡more help information of ixkit ,please refer to <a href="http://www.ixkit.com">ixkit</a>
usage of cURL, please refer to <a href="https://github.com/astaxie/bat">bat</a>
