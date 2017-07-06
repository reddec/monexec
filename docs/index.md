# MONEXEC

[![GitHub release](https://img.shields.io/github/release/reddec/monexec.svg)]()
[![license](https://img.shields.io/github/license/reddec/monexec.svg)]()
[![](https://godoc.org/github.com/reddec/monexec/monexec?status.svg)](http://godoc.org/github.com/reddec/monexec/monexec)

It's tool for controlling processes like a **supervisord** but with some important features:
* Easy to use - no dependencies. Just a single binary file pre-compilled for most major platforms
* Easy to hack - monexec can be used as a Golang library with clean and simple architecture
* Integrated with Consul - optionally, monexec can register all running processes as services and deregister on fail
* Supports gracefull and fast shutdown by signals
* Developed for used inside Docker containers
* Different strategies for processes

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

Consul configuration is available only by [Go Consul API environment variables](https://godoc.org/github.com/hashicorp/consul/api#pkg-constants) (improvments for this are in roadmap). 
Most usefull - `CONSUL_HTTP_ADDR` that specifies URL to Consul agent.

## Examples:

**Register in local agent:**

```bash
monexec run -l srv1 --consul forever -- nc -l 9000
```

**Register in remote agent:**

Suppose Consul agent is running in host `registry`

```bash
export CONSUL_HTTP_ADDR="http://registry:8500"
monexec run -l srv1 --consul forever -- nc -l 9000
```


# Usage

`monexec [common-flags...] <command> [command-flags...] [args,...]`

All flags can be set by environment variables with prefix `MONEXEC_`. For example flag `--label sample` can be set as `export MONEXEC_LABEL="sample"`

Common flags:

* `--label -l <label name>` - mark executable with specific ID. Used as service ID in Consul and in logs. By default ID will be randomly generated.
* `--consul` - enable consul integration

## Commands

### run
Run single executable with specified strategies (modes)

**Usage:**
`monexec run [flags...] <mode> <executable> [args...]`

**Example:**
`monexec run forver -- nc -l 9000` - will run command `nc -l 9000` and restart it forever if needed with default timeout

**Modes:**

* `critical` - fail on first error in application
* `restart` - always restart application and ignore exit codes
* `forever` - same as `restart` but with unlimited restart retries
* `oneshot` - run application only once and ignore exit codes

**Flags:**
* `-r, --retries=5` - Maximum restart retries. Negative means infinity
* `-restart-timeout=5s` - Timeout before restart
* `--start-timeout=3s` - Timeout to check that process is started. Only after this timeout process will be marked alive and tried to be registered in Consul (if enabled)
* `--stop-timeout=5s` - Timeout for graceful shutdown. Application first got signal `SIGTERM` and after this timeout `SIGKILL`
* `-w, --workdir=WORKDIR` - Working directory. By default - running directory

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
label: Netcat Sample Service
command: nc
args:
- "-l"
- "9000"
retries: 5
critical: false
stop_timeout: 5s
start_timeout: 3s
restart_timeout: 5s
```


## gen
Generate configuration file based on `run` like arguments for `start` sources.

**Usage:**
Same as `run`

For example, during development we are using 

```bash
monexec run -l srvExt1 --retries 10 --start-timeout 15s restart -- java -jar srvExt1.jar
```

We want to save this settings into configuration file. Just change `run` to `gen`:

```bash
monexec gen -l srvExt1 --retries 10 --start-timeout 15s restart -- java -jar srvExt1.jar
```

and get

```yaml
label: srvExt1
command: java
args:
- -jar
- srvExt1.jar
retries: 10
stop_timeout: 5s
start_timeout: 15s
restart_timeout: 5s
```
