# Monexec
***MON**itoring **EXE**cutables*

[![GitHub release](https://img.shields.io/github/release/reddec/monexec.svg)](https://github.com/reddec/monexec/releases)
[![license](https://img.shields.io/github/license/reddec/monexec.svg)](https://github.com/reddec/monexec)
[![](https://godoc.org/github.com/reddec/monexec/monexec?status.svg)](http://godoc.org/github.com/reddec/monexec/monexec)

It’s tool for controlling processes like a supervisord but with some important features:

* Easy to use - no dependencies. Just a single binary file pre-compilled for most major platforms
* Easy to hack - monexec can be used as a Golang library with clean and simple architecture
* Integrated with Consul - optionally, monexec can register all running processes as services and deregister on fail
* Optional notification to Telegram
* Supports gracefull and fast shutdown by signals
* Developed for used inside Docker containers
* Different strategies for processes

## Installing

Precompilled binaries: [release page](https://github.com/reddec/monexec/releases)

From source (required Go toolchain):

```
go get -v -u github.com/reddec/monexec/...
```

## Documentation

Usage: [https://reddec.github.io/monexec/](https://reddec.github.io/monexec/)

API: [Godoc](http://godoc.org/github.com/reddec/monexec/monexec)
