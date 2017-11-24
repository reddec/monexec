package plugins

import (
	"bytes"
	"github.com/reddec/container"
	"time"
	"log"
	"os"
	"net/smtp"
	"net"
	"errors"
	"path/filepath"
)

type Email struct {
	Smtp         string   `yaml:"smtp"`
	From         string   `yaml:"from"`
	Password     string   `yaml:"password"`
	To           []string `yaml:"to"`
	Template     string   `yaml:"template"`
	TemplateFile string   `yaml:"templateFile"` // template file (relative to config dir) has priority. Template supports basic utils
	Services     []string `yaml:"services"`

	log         *log.Logger
	hostname    string
	servicesSet map[string]bool
	workDir     string
}

func (c *Email) renderAndSend(params map[string]interface{}) {
	message := &bytes.Buffer{}

	parser, err := parseFileOrTemplate(c.TemplateFile, c.Template, c.log)
	if err != nil {
		c.log.Println("failed parse template:", err)
		return
	}
	renderErr := parser.Execute(message, params)
	if renderErr != nil {
		c.log.Println("failed render:", renderErr, "; params:", params)
		return
	}

	c.log.Println(message.String())
	host, _, _ := net.SplitHostPort(c.Smtp)
	auth := smtp.PlainAuth("", c.From, c.Password, host)
	err = smtp.SendMail(c.Smtp, auth, c.From, c.To, message.Bytes())
	if err != nil {
		c.log.Println("failed send mail:", err)
	} else {
		c.log.Println("sent")
	}
}

func (c *Email) Spawned(runnable container.Runnable, id container.ID) {
	if c.servicesSet[runnable.Label()] {
		params := map[string]interface{}{
			"action":   "spawned",
			"id":       id,
			"label":    runnable.Label(),
			"hostname": c.hostname,
			"time":     time.Now().String(),
		}
		c.renderAndSend(params)
	}
}

func (c *Email) Prepare() error {
	c.servicesSet = makeSet(c.Services)
	c.log = log.New(os.Stderr, "[email] ", log.LstdFlags)
	c.hostname, _ = os.Hostname()
	return nil
}

func (c *Email) Stopped(runnable container.Runnable, id container.ID, err error) {
	if c.servicesSet[runnable.Label()] {
		params := map[string]interface{}{
			"action":   "stopped",
			"id":       id,
			"error":    err,
			"label":    runnable.Label(),
			"hostname": c.hostname,
			"time":     time.Now().String(),
		}
		c.renderAndSend(params)
	}
}

func (a *Email) MergeFrom(other interface{}) (error) {
	b := other.(*Email)
	if a.From == "" {
		a.From = b.From
	}
	if a.From != b.From {
		return errors.New("different from address")
	}
	if a.Smtp == "" {
		a.Smtp = b.Smtp
	}
	if a.Smtp != b.Smtp {
		return errors.New("different smtp servers")
	}
	if a.Template == "" {
		a.Template = b.Template
	}
	if a.Template != b.Template {
		return errors.New("different templates")
	}
	a.TemplateFile = realPath(a.TemplateFile, a.workDir)
	b.TemplateFile = realPath(b.TemplateFile, b.workDir)
	if a.TemplateFile == "" {
		a.TemplateFile = b.TemplateFile
	}
	if a.TemplateFile != b.TemplateFile {
		return errors.New("different template files")
	}
	if a.Password == "" {
		a.Password = b.Password
	}
	if a.Password != b.Password {
		return errors.New("different password")
	}
	a.To = unique(append(a.To, b.To...))
	a.Services = append(a.Services, b.Services...)
	return nil
}

func init() {
	registerPlugin("email", func(file string) PluginConfig {
		return &Email{workDir: filepath.Dir(file)}
	})
}
