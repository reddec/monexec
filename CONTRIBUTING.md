# Contributing to the project

First and most important: thanks a lot for choosing this project as a place to which you are planning to invest your time and knowledge!

This guide contains as much information as possible that should help you to join to the development.

Please feel free to enrich this document through PR requests.

## Pull requests

I am highly recommend to discuss a new feature through issues section before development to be sure that feature is not yet already done or someone already doing it.

The project trying to follow simple 'git flow': 
* `master` branch should be always available for 'testing' use (i.e. available for build, run and not to fail)
* **tagged** commit means stable version
* (recommended, but not critical) new features should come from `feature/some-name`  branches, bug-fixes from `bug/bug-name` and so on. Don't stuck on it too much: if you can't really decide which branch name required - use `feature/` prefix

Simple way to understand project is to look at plugins implementation. They are relatively simple and standardized

Useful resources:
* Good tutorial [how to create PR](https://help.github.com/en/articles/creating-a-pull-request) from GitHub.

## Environment

 - **Go**: at this time project aims to the latest golang toolchain (1.12 now), however if you will use newer version please note it in the PR.
	 - [Go SDK](https://golang.org/doc/install)
 - **Modules**: the project is using go1.11+ modules without `vendor` directory.
	 - [How to use modules](https://blog.golang.org/using-go-modules)
 - **IDE**: Me (reddec) personally using [Jetbrains Goland](https://www.jetbrains.com/go/) as a professional user, however you may use free combination of IDEA community edition + Golang plugin. Or whatever you want - just please **don't commit IDE** files to the repository

## Prepare and build

 1. Install go and setup [GOPATH](https://github.com/golang/go/wiki/GOPATH)
 2. Clone project through git
 3. In the project directory download all dependency for the command: `GO111MODULE=on go get -v ./...`
 4. Build it: `GO111MODULE=on go build ./cmd/...`
 5. Run it: `./monexec --help`

## Project structure

* `docs` - site and project documentation
* `monexec` - root configuration
* `plugins` - plugin implementation
* `pool` - core (need refactor) supervisor logic
* `sample` - examples and configuration files
* `ui` - sub-module for UI
* `swagger.yaml` - swagger definition of HTTP plugin endpoints