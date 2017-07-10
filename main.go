package main

import (
	"github.com/reddec/monexec/monexec"
	"gopkg.in/yaml.v2"
	"os"
	"io/ioutil"
	"strings"
	"github.com/hashicorp/consul/api"
	"log"
	"time"
	"github.com/Pallinder/go-randomdata"
	"gopkg.in/alecthomas/kingpin.v2"
	"context"
	"os/signal"
	"syscall"
	"sync"
)

func loadConfigs(locations ...string) ([]*monexec.Executable, error) {
	var ans []*monexec.Executable
	var files []os.FileInfo
	for _, location := range locations {
		if stat, err := os.Stat(location); err != nil {
			return nil, err
		} else if stat.IsDir() {
			fs, err := ioutil.ReadDir(location)
			if err != nil {
				return nil, err
			}
			files = fs
		} else {
			files = []os.FileInfo{stat}
		}
		for _, info := range files {
			if strings.HasSuffix(info.Name(), ".yml") || strings.HasSuffix(info.Name(), ".yaml") {
				data, err := ioutil.ReadFile(info.Name())
				if err != nil {
					return nil, err
				}
				var exe monexec.Executable
				err = yaml.Unmarshal(data, &exe)
				if err != nil {
					return nil, err
				}
				ans = append(ans, &exe)
			}
		}
	}
	return ans, nil
}

func buildExec(mode, executable string, args []string, restart int, label string) *monexec.Executable {
	mon := monexec.Monitor{}
	var exe *monexec.Executable
	switch mode {
	case "oneshot":
		exe = mon.Oneshot(executable, args...)
	case "critical":
		exe = mon.Critical(executable, args...)
	case "forever":
		exe = mon.Forever(executable, args...)
	case "restart":
		exe = mon.Restart(restart, executable, args...)
	default:
		panic("Unknown mode " + mode)
	}
	exe.Label = label
	return exe
}

func dummyConsumer(events <-chan monexec.Event) {
	for range events {

	}
}

func consulEventsConsumer(events <-chan monexec.Event) {

	config := api.DefaultConfig()
	client, err := api.NewClient(config)
	if err != nil {
		log.Println("Unexpected error while connecting to Consul:", err)
		dummyConsumer(events)
		return
	}

	state := map[*monexec.Executable]chan monexec.Event{}

	interval := 5 * time.Second

	wg := sync.WaitGroup{}

	for event := range events {
		ch, ok := state[event.Executable]
		if !ok {
			ch = make(chan monexec.Event, 1)
			state[event.Executable] = ch
			wg.Add(1)
			go func(event monexec.Event, ch chan monexec.Event) {
				defer wg.Done()
				runService(ch, interval, client)
			}(event, ch)
		}
		ch <- event
	}
	for _, ch := range state {
		close(ch)
	}
	wg.Wait()

}

func main() {
	app := kingpin.New("monexec", "Light supervisor for monitoring processes")
	app.Version("1.0.0").Author("Baryshnikov Alexander <dev@baryshnikov.net>")

	cmdGen := app.Command("gen", "Generate basic configuration file for executable")

	label := app.Flag("label", "Label for service").Short('l').Default(randomdata.Noun() + "-" + randomdata.Adjective()).String()
	retries := app.Flag("retries", "Restart count").Short('r').Default("5").Int()
	restartInterval := app.Flag("restart-timeout", "Timeout before restart").Default("5s").Duration()
	startTimeout := app.Flag("start-timeout", "Timeout to check that process is started").Default("3s").Duration()
	stopTimeout := app.Flag("stop-timeout", "Timeout for graceful shutdown").Default("5s").Duration()
	workdir := app.Flag("workdir", "Working directory").Short('w').String()

	genMode := cmdGen.Arg("mode", "Mode types").Required().Enum("oneshot", "critical", "forever", "restart")
	genExecutable := cmdGen.Arg("executable", "Applications to start").Required().String()
	genArgs := cmdGen.Arg("arg", "Arguments").Strings()

	start := app.Command("start", "Start supervisor")
	configLocations := start.Arg("config", "Config file or directory with .yaml/.yml files").Strings()

	run := app.Command("run", "Run single executable")
	runMode := run.Arg("mode", "Mode types").Required().Enum("oneshot", "critical", "forever", "restart")
	runExecutable := run.Arg("executable", "Applications to start").Required().String()
	runArgs := run.Arg("arg", "Arguments").Strings()

	consulEnabled := app.Flag("consul", "Enable consul auto-registration (used ENV for config)").Bool()

	app.DefaultEnvars()

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case "gen":
		exe := buildExec(*genMode, *genExecutable, *genArgs, *retries, *label)
		exe.RestartTimeout = *restartInterval
		exe.StartTimeout = *startTimeout
		exe.StopTimeout = *stopTimeout
		exe.WorkDir = *workdir
		data, err := yaml.Marshal(exe)
		if err != nil {
			panic(err)
		}
		os.Stdout.Write(data)
	case "start":
		execs, err := loadConfigs(*configLocations...)
		if err != nil {
			panic(err)
		}

		mon := monexec.Monitor{Executables: execs}
		ctx, stp := context.WithCancel(context.Background())

		c := make(chan os.Signal, 3)
		signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGHUP)
		go func() {
			for range c {
				stp()
				break
			}
		}()

		events := mon.Start(ctx)

		if *consulEnabled {
			consulEventsConsumer(events)
		} else {
			dummyConsumer(events)
		}

	case "run":
		exe := buildExec(*runMode, *runExecutable, *runArgs, *retries, *label)
		exe.RestartTimeout = *restartInterval
		exe.StartTimeout = *startTimeout
		exe.StopTimeout = *stopTimeout
		exe.WorkDir = *workdir

		mon := monexec.Monitor{Executables: []*monexec.Executable{exe}}
		ctx, stp := context.WithCancel(context.Background())

		c := make(chan os.Signal, 2)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		go func() {
			for range c {
				stp()
				break
			}
		}()

		events := mon.Start(ctx)

		if *consulEnabled {
			consulEventsConsumer(events)
		} else {
			dummyConsumer(events)
		}
	}

}
