package main

import (
	"gopkg.in/alecthomas/kingpin.v2"
	"github.com/Pallinder/go-randomdata"
	".."
	"gopkg.in/yaml.v2"
	"os"
	"io/ioutil"
	"strings"
	"context"
	"os/signal"
	"syscall"
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

func main() {
	app := kingpin.New("monexec", "Light supervisor for monitoring processes")
	app.Version("1.0.0").Author("Baryshnikov Alexander <dev@baryshnikov.net>")

	cmdGen := app.Command("gen", "Generate basic configuration file for executable")

	genLabel := cmdGen.Flag("label", "Label for service").Short('l').Default(randomdata.Noun() + "-" + randomdata.Adjective()).String()
	genMode := cmdGen.Arg("mode", "Mode types").Required().Enum("oneshot", "critical", "forever", "restart")
	genExecutable := cmdGen.Arg("executable", "Applications to start").Required().String()
	genArgs := cmdGen.Arg("arg", "Arguments").Strings()
	genRestart := cmdGen.Flag("retries", "Restart count").Short('r').Default("5").Int()

	start := app.Command("start", "Start supervisor")
	configLocations := start.Arg("config", "Config file or directory with .yaml/.yml files").Strings()

	run := app.Command("run", "Run single executable")
	runMode := run.Arg("mode", "Mode types").Required().Enum("oneshot", "critical", "forever", "restart")
	runExecutable := run.Arg("executable", "Applications to start").Required().String()
	runArgs := run.Arg("arg", "Arguments").Strings()
	runRestart := run.Flag("retries", "Restart count").Short('r').Default("5").Int()
	runLabel := run.Flag("label", "Label for service").Short('l').Default(randomdata.Noun() + "-" + randomdata.Adjective()).String()
	runRestartInterval := run.Flag("restart-timeout", "Timeout before restart").Default("5s").Duration()
	runStartTimeout := run.Flag("start-timeout", "Timeout to check that process is started").Default("3s").Duration()
	runStopTimeout := run.Flag("stop-timeout", "Timeout for graceful shutdown").Default("5s").Duration()
	runWorkdir := run.Flag("workdir", "Working directory").Short('w').String()

	app.DefaultEnvars()
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case "gen":
		exe := buildExec(*genMode, *genExecutable, *genArgs, *genRestart, *genLabel)
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
			for _ = range c {
				stp()
			}
		}()

		err = mon.Run(ctx)
		if err != nil {
			panic(err)
		}
	case "run":
		exe := buildExec(*runMode, *runExecutable, *runArgs, *runRestart, *runLabel)
		exe.RestartTimeout = *runRestartInterval
		exe.StartTimeout = *runStartTimeout
		exe.StopTimeout = *runStopTimeout
		exe.WorkDir = *runWorkdir

		mon := monexec.Monitor{Executables: []*monexec.Executable{exe}}
		ctx, stp := context.WithCancel(context.Background())

		c := make(chan os.Signal, 3)
		signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGHUP)
		go func() {
			for _ = range c {
				stp()
			}
		}()
		err := mon.Run(ctx)
		if err != nil {
			panic(err)
		}
	}

}
