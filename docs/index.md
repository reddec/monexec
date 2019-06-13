
  
  
# MONEXEC  
  
[![GitHub release](https://img.shields.io/github/release/reddec/monexec.svg)](https://github.com/reddec/monexec/releases)  
[![license](https://img.shields.io/github/license/reddec/monexec.svg)](https://github.com/reddec/monexec)  
[![](https://godoc.org/github.com/reddec/monexec/monexec?status.svg)](http://godoc.org/github.com/reddec/monexec/monexec)  
  
It's tool for controlling processes like a **supervisord** but with some important features:  
* Easy to use - no dependencies. Just a single binary file pre-compilled for most major platforms  
* Easy to hack - monexec can be used as a Golang library with clean and simple architecture  
* Integrated with Consul - optionally, monexec can register all running processes as services and deregister on fail  
* Supports gracefull and fast shutdown by signals  
* Developed for used inside Docker containers  
* Different strategies for processes  
* Support template-based email notification  
  
[download for most major platform](https://github.com/reddec/monexec/releases)  
  
# Installing  
  
Precompilled binaries:  
[release page](https://github.com/reddec/monexec/releases)  
  
From source (required Go toolchain):  
  
```  
go get -v -u github.com/reddec/monexec/...  
```  
  
# How to integrate with Consul  
  
Consul is a service registry and service discover system. MONEXEC can automatically register application in Consul as a service.  
  
Auto(de)registration available for `run` or `start` commands.  
  
Use general flag `--consul` (or env var `MONEXEC_CONSUL=true`) for enable Consul integration. Monexec will try register and update status of service in Consul local agent.  
  
Monexec will continue work even if Consul becomes unavailable.  
  
Consul address by default located to localhost, but can be overrided by `--consul-address` or `MONEXEC_CONSUL_ADDRESS` environment variable.  
  
Additional Consul configuration is available only by [Go Consul API environment variables](https://godoc.org/github.com/hashicorp/consul/api#pkg-constants) (improvments for this are in roadmap).  
  
## Examples:  
  
**Register in local agent:**  
  
Temporary (will auto de-registrate service in a critical state or after gracefull shutdown)  
  
```bash  
monexec run -l srv1 --consul -- nc -l 9000  
```  

Permanent  
  
```bash  
monexec run -l srv1 --consul --consul-permanent -- nc -l 9000  
```  
  
**Register in remote agent:**  
  
Suppose Consul agent is running in host `registry`  
  
```bash  
monexec run --consul --consul-address "http://registry:8500" -l srv1 -- nc -l 9000
```  
  
# How to integrate with Telegram  
  
Since `0.1.1` you can receive notifications over Telegram.  
  
You have to know:  
  
* BOT token : can be obtained here http://t.me/botfather  
* Receipients ChatID's : can be obtained here http://t.me/MyTelegramID_bot  
  
Message template (based on Golang templates) also required. We recommend use this:  
  
```  
*{{.label}}*  
Service {{.label}} {{.action}}  
{{if .error}}⚠️ *Error:*  {{.error}}{{end}}_time: {{.time}}_  
_host: {{.hostname}}_  
```  
  
Available params:  
  
* `.label` - name of service  
* `.action` - servce action. Can be `spawned` or `stopped`  
* `.time` - current time in UTC format with timezone  
* `.error` - error message available only on `stopped` action  
* `.hostname` - current hostname  
  
Configuration avaiable only from .yaml files:  
  
```yaml  
telegram:  
 # BOT token token: "123456789:AAAAAAAAAAAAAAAAAAAAAA_BBBBBBBBBBBB" 
 services: # services that will be monitored 
 - "listener2" 
 recipients: # List of telegrams chat id 
 - 123456789 
 template: | *{{.label}}* Service {{.label}} {{.action}} {{if .error}}⚠️ *Error:*  {{.error}}{{end}} _time: {{.time}}_ _host: {{.hostname}}_
 ```  
  
Since `0.1.4` you also can specify `templateFile` instead of `template`  
  
# How to integrate with email  
  
Since `0.1.3` you can receive notifications over email.  
  
If you are using Google emails (tested):  
  
* Obtain application password https://myaccount.google.com/apppasswords  
* SMTP server will be: `smtp.gmail.com:587`  
  
Message template (based on Golang templates) also required. We recommend use this:  
  
```  
Content-Type: text/html  
Subject: {{.label}} {{.action}}  
  
<h2>{{.label}}</h2>  
  
<table>  
 <tr> <th>Label</th> <td>{{.label}}</td> </tr> <tr> <th>ID</th> <td>({{.id}})</td> </tr> <tr> <th>Action</th> <td>{{.action}}</td> </tr> <tr> <th>Hostname</th> <td>{{.hostname}}</td> </tr> <tr> <th>Local time</th> <td>{{.time}}</td> </tr> <tr> <th>User</th> <td>{{env "USER"}}</td> </tr></table>  
```  
  
Available params:  
  
* `.label` - name of service  
* `.action` - servce action. Can be `spawned` or `stopped`  
* `.time` - current time in UTC format with timezone  
* `.error` - error message available only on `stopped` action  
* `.hostname` - current hostname  
  
Plus all operations from http://masterminds.github.io/sprig/ (like `env` or `upper`)  
  
Configuration avaiable only from .yaml files:  
  
  
```yaml  
  
email:  
 services: 
 - myservice 
  smtp: "smtp.gmail.com:587" 
  from: "example-monitor@gmail.com" 
  password: "xyzzzyyyzyyzyz" 
  to: 
  - "admin1@example.com" 
  - "admin2@example.com" 
  template: | Subject: {{.label}}  
 Service {{.label}} {{.action}} templateFile: "./email.html"
 ```  
  
`template` will be used as fallback for `templateFile`. If template file location is not absolute, it will be calculated  
from configuration directory.  
  
# How to integrate with HTTP  
  
Since `0.1.4` you can send notifications over HTTP  
  
* Supports any kind of methods (by default - `POST` but can be changed in `http.method`)  
* **Body** - template-based text same as in `Telegram` or `Email` plugin  
* **URL** - also template-based text (yes, with same rules as in `body` ;-) )  
* **Headers** - you can also provide any headers (no, no templates here)  
* **Timeout** - limit time for request. By default - `20s`  
  
Configuration avaiable only from .yaml files:  
  
```yaml  
http:  
 services: 
 - myservice 
 url: "http://example.com/{{.label}}/{{.action}}" 
 templateFile: "./body.txt"
 ```  
  
`template` will be used as fallback for `templateFile`. If template file location is not absolute, it will be calculated from configuration directory.  
  
  
|Parameter     | Type     | Required | Default | Description |  
|--------------|----------|----------|---------|-------------|  
|`url`         |`string`  |   yes    |         | Target URL  
|`method`      |`string`  |   no     | POST    | HTTP method  
|`services`    |`list`    |   yes    |         | List of services that will trigger plugin  
|`headers`     |`map`     |   no     | {}      | Map (string -> string) of additional headers per request  
|`timeout`     |`duration`|   no     | 20s     | Request timeout  
|`template`    |`string`  |   no     | ''      | Template string  
|`templateFile`|`string`  |   no     | ''      | Path to file of template (more priority then `template`, but `template` will be used as fallback)  
  
# Usage  
  
`monexec <command> [command-flags...] [args,...]`  
  
All flags can be set by environment variables with prefix `MONEXEC_`. For example flag `--label sample` can be set as `export MONEXEC_LABEL="sample"`  
  
# How to enable REST API  
  
Since `0.1.6` you can enable simple REST API by adding `rest` plugin.  
  
Full version  
  
```yaml  
rest:  
 listen: "localhost:9900" 
 cors: false
 ```  
  
_cors option added in `0.1.9`_  
  
or minimal (default is `localhost:9900`)  
  
```yaml  
rest:  
```  
  
API documentation see in swagger.yaml file in repository  
  
**WEB UI** enable on `/ui` path  
  
![screencapture-127-0-0-1-9000-2018-06-28-20_46_16](https://user-images.githubusercontent.com/6597086/42038135-c961b11a-7b1c-11e8-9437-44de6b36510c.png)  
  
## Commands  
  
### run  
Run single executable  
  
**Usage:**  
`monexec run [flags...] <executable> [args...]`  
  
**Example:**  
`monexec run -- nc -l 9000` - will run command `nc -l 9000` and restart it forever if needed with default timeout  
  
**Flags:**  
  
* `--generate` Generate to stdout YAML configuration based on args  instead of run  
* `-r, --restart-count=-1` Restart count (negative means infinity)  
* `-d, --restart-delay=5s` Delay before restart  
* `-g, --graceful-timeout=5s` Timeout for graceful shutdown. Application first got signal `SIGTERM` and after this timeout `SIGKILL`  
* `-l, --label=LABEL` Label name for executable. Default - autogenerated  
* `-w, --workdir=WORKDIR` Workdir for executable  
* `-e, --env=ENV ...` Environment addition variables  
* `--consul` Enable consul integration  
* `--consul-address="http://localhost:8500"` Consul address  
* `--consul-permanent` Keep service in consul auto timeout  
* `--consul-ttl=3s` Keep-alive TTL for services  
* `--consul-unreg=1m` Timeout after for auto de-registration  
  
## start  
Start processes based on YAML configuration files  
  
**Usage:**  
`monexec start <config file or dir,...>`  
  
**Example:**  
  
`monexec start ./*.yaml`  
  
Configuration sources can be multiple directories and/or files. Files must contain valid YAML content and have `.yaml` or `.yml` extension.  
  
Minimal configuration file:  
  
```yaml  
command: path/to/executable  
```  
  
Full sample of configuration file:  
  
```yaml  
services:  
- label: Netcat Sample Service  
 command: nc 
 args: 
  - -l 
  - "9000" 
 stop_timeout: 5s 
 restart_delay: 5s 
 restart: -1
 consul:  
  url: http://localhost:8500 
  ttl: 3s 
  timeout: 1m0s
 ```
  
# How to generate sample config  
  
Generate configuration file based on `run` like arguments: just add `--generate`  
  
**Usage:**  
  
Same as `run`  
  
For example, during development we are using  
  
```bash  
monexec run -l srvExt1 --consul --restart-count 10 restart -- java -jar srvExt1.jar  
```  
  
We want to save this settings into configuration file. Just add `--generate`  
  
```bash  
monexec run --generate -l srvExt1  --consul --restart-count 10 restart -- java -jar srvExt1.jar  
```  
  
and get  
  
```yaml  
services:  
- label: srvExt1  
 command: restart 
 args: 
 - java 
 - -jar -
  srvExt1.jar 
  stop_timeout: 5s 
  restart_delay: 5s 
  restart: 10
  consul:  
    url: http://localhost:8500 
    ttl: 3s 
    timeout: 1m0s 
    register: 
      - srvExt1
 ```  

# How to log to file a service  
  
Since `0.1.5` you can copy content of STDERR/STDOUT  service output to specific file: option `logFile` in service section. If file path not absolute log file is putted relative to working directory.  
  
```yaml  
services:  
- label: listener4  
 command: nc 
 logFile: /var/log/listener4.log 
 args: 
 - -v 
 - -l 
 - 9001 
 stop_timeout: 5s 
 restart_delay: 5s 
 restart: -1  
```  
  
# Critical services  
  
  
When critical services stopped, all other processes have to be stopped also  
  
  
Add section `critical` to configuration:  
  
  
```yaml  
services:  
- label: srvExt1  
 command: restart 
 args: 
 - java 
 - -jar 
 - srvExt1.jar 
 stop_timeout: 5s 
 restart_delay: 5s 
 restart: 10  
- label: consul  
 command: restart 
 args: 
 - consul 
 - agen 
 - -dev 
 - -bootstrap 
 - -uiconsul:  
 url: http://localhost:8500 
 ttl: 3s 
 timeout: 1m0s
  register: 
  - srvExt1critical:  
  - consul  
```  
  
# Raw stdout  
  
For several reasons (i.e. use in a bash tools) raw stdout is required from application.  
  
Since `0.1.12` to disable all prefixes in STDOUT (in STDERR they will still persists) use flag `--raw, -R`.  
  
**Example:**  
  
```bash  
monexec -R echo 123 > sample.txt # run echo command. Use CTRL+C to interrupt when needed
```  
  
The file `sample.txt` will now contains ONLY result of echo command (i.e. `123`)

# Environment variables from file

General environment variables processing:

1. system environments
2. environments defined by `--env, -e` (CLI) or in `environment` (config)
3. (since `0.1.14`), environments files in order as they defined by `--env-file, -E` (CLI) or in `envFile` (config) 

If environment file couldn't be read, the application **can still be launched** - only warning in log presented.

Environment file format:

* pair: KEY=VALUE  
* comment: line started with #  
* empty lines or invalid (without = symbol) ignored  
* there is no way to escape new line symbol

Example of file:

```env
# some comment  
  
  
it was empty lines and this is a broken record 
MAIL=will be replaced  
MAIL=owner@reddec.net  
SUBJECT=multiple = are supported
```
