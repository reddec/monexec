package plugins

import (
	"log"
	"os"
	"net/smtp"
	"net"
	"errors"
	"path/filepath"
	"github.com/reddec/monexec/pool"
	"context"
)

type Email struct {
	Smtp     string   `yaml:"smtp"`
	From     string   `yaml:"from"`
	Password string   `yaml:"password"`
	To       []string `yaml:"to"`
	Services []string `yaml:"services"`
	withTemplate      `mapstructure:",squash" yaml:",inline"`

	log         *log.Logger
	hostname    string
	servicesSet map[string]bool
	workDir     string
}

func (c *Email) renderAndSend(message string) {
	c.log.Println(message)
	host, _, _ := net.SplitHostPort(c.Smtp)
	auth := smtp.PlainAuth("", c.From, c.Password, host)
	err := smtp.SendMail(c.Smtp, auth, c.From, c.To, []byte(message))
	if err != nil {
		c.log.Println("failed send mail:", err)
	} else {
		c.log.Println("sent")
	}
}

func (c *Email) OnSpawned(ctx context.Context, sv pool.Instance) {}

func (c *Email) OnStarted(ctx context.Context, sv pool.Instance) {
	label := sv.Config().Name
	if c.servicesSet[label] {
		content, renderErr := c.renderDefault("spawned", label, label, nil, c.log)
		if renderErr != nil {
			c.log.Println("failed render:", renderErr)
		} else {
			c.renderAndSend(content)
		}
	}
}

func (c *Email) OnStopped(ctx context.Context, sv pool.Instance, err error) {
	label := sv.Config().Name
	if c.servicesSet[label] {
		content, renderErr := c.renderDefault("stopped", label, label, err, c.log)
		if renderErr != nil {
			c.log.Println("failed render:", renderErr)
		} else {
			c.renderAndSend(content)
		}
	}
}

func (p *Email) OnFinished(ctx context.Context, sv pool.Instance) {}

func (c *Email) Prepare(ctx context.Context, pl *pool.Pool) error {
	c.servicesSet = makeSet(c.Services)
	c.log = log.New(os.Stderr, "[email] ", log.LstdFlags)
	c.hostname, _ = os.Hostname()
	return nil
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
func (a *Email) Close() error { return nil }
func init() {
	registerPlugin("email", func(file string) PluginConfigNG {
		return &Email{workDir: filepath.Dir(file)}
	})
}
