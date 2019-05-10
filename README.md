# Monexec

![Mx](docs/logo.svg)

***MON**itoring **EXE**cutables*

[![GitHub release](https://img.shields.io/github/release/reddec/monexec.svg)](https://github.com/reddec/monexec/releases)
[![license](https://img.shields.io/github/license/reddec/monexec.svg)](https://github.com/reddec/monexec)
[![](https://godoc.org/github.com/reddec/monexec/monexec?status.svg)](http://godoc.org/github.com/reddec/monexec/monexec)
[![Snap Status](https://build.snapcraft.io/badge/reddec/monexec.svg)](https://build.snapcraft.io/user/reddec/monexec)

It’s tool for controlling processes like a supervisord but with some important features:

* Easy to use - no dependencies. Just a single binary file pre-compilled for most major platforms
* Easy to hack - monexec can be used as a Golang library with clean and simple architecture
* Integrated with Consul - optionally, monexec can register all running processes as services and deregister on fail
* Optional notification to Telegram
* Supports gracefull and fast shutdown by signals
* Developed for used inside Docker containers
* Different strategies for processes
* Support additional environment files
* Support template-based email notification
* Support HTTP notification
* REST API (see swagger.yaml)
* Web UI (if REST API enabled)

[Buy me a coffe on Patreon](https://www.patreon.com/bePatron?u=15369842)

![screen1](https://reddec.net/images/project/monexec/screen1.png)
![screen2](https://reddec.net/images/project/monexec/screen2.gif)
## Installing

[![Get it from the Snap Store](https://snapcraft.io/static/images/badges/en/snap-store-black.svg)](https://snapcraft.io/monexec)

* [snapcraft: monexec](https://snapcraft.io/monexec)

* Precompilled binaries: [release page](https://github.com/reddec/monexec/releases)

* From source (required Go toolchain):

```
go get -v -u github.com/reddec/monexec/...
```

recommended way is snap

## Documentation

Usage: [https://reddec.github.io/monexec/](https://reddec.github.io/monexec/)

API: [Godoc](http://godoc.org/github.com/reddec/monexec/monexec)


## Examples

See documentation for details [https://reddec.github.io/monexec/](https://reddec.github.io/monexec/)

### Run from cmd

```bash
monexec run -l srv1 --consul -- nc -l 9000
```

### Run from config

```bash
monexec start ./myservice.yaml
```

### Notifications

Add notification to Telegram

```yaml
telegram:
  # BOT token
  token: "123456789:AAAAAAAAAAAAAAAAAAAAAA_BBBBBBBBBBBB"
  services:
      # services that will be monitored
      - "listener2"
  recipients:
      # List of telegrams chat id
      - 123456789
  template: |
    *{{.label}}*
    Service {{.label}} {{.action}}
    {{if .error}}⚠️ *Error:*  {{.error}}{{end}}
    _time: {{.time}}_
    _host: {{.hostname}}_
```

#### Email

Add email notification

```yaml
email:
  services:
    - myservice
  smtp: "smtp.gmail.com:587"
  from: "example-monitor@gmail.com"
  password: "xyzzzyyyzyyzyz"
  to:
    - "admin1@example.com"
  template: |
    Subject: {{.label}}

    Service {{.label}} {{.action}}
```

#### HTTP

Add HTTP request as notification

```yaml
http:
  services:
    - myservice
  url: "http://example.com/{{.label}}/{{.action}}"
  templateFile: "./body.txt"
```
