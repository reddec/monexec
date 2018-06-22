package plugins

import (
	"log"
	"os"
	"gopkg.in/telegram-bot-api.v4"
	"errors"
	"path/filepath"
	"github.com/reddec/monexec/pool"
	"context"
)

type Telegram struct {
	Token      string   `yaml:"token"`
	Recipients []int64  `yaml:"recipients"`
	Services   []string `yaml:"services"`
	withTemplate        `mapstructure:",squash" yaml:",inline"`

	servicesSet map[string]bool  `yaml:"-"`
	logger      *log.Logger      `yaml:"-"`
	bot         *tgbotapi.BotAPI `yaml:"-"`
	workDir     string
	hostname    string
}

func (c *Telegram) Prepare(ctx context.Context, pl *pool.Pool) error {
	c.servicesSet = make(map[string]bool)
	for _, srv := range c.Services {
		c.servicesSet[srv] = true
	}
	c.logger = log.New(os.Stderr, "[telegram] ", log.LstdFlags)
	bot, err := tgbotapi.NewBotAPI(c.Token)
	if err != nil {
		return err
	}
	c.bot = bot
	c.hostname, _ = os.Hostname()
	return nil
}

func (p *Telegram) OnSpawned(ctx context.Context, sv pool.Instance) {}

func (c *Telegram) OnStarted(ctx context.Context, sv pool.Instance) {
	if c.servicesSet[sv.Config().Name] {
		content, renderErr := c.renderDefault("spawned", string(sv.Config().Name), sv.Config().Name, nil, c.logger)
		if renderErr != nil {
			c.logger.Println("failed render:", renderErr)
		} else {
			c.renderAndSend(content)
		}
	}
}

func (c *Telegram) OnStopped(ctx context.Context, sv pool.Instance, err error) {
	if c.servicesSet[sv.Config().Name] {
		content, renderErr := c.renderDefault("stopped", string(sv.Config().Name), sv.Config().Name, err, c.logger)
		if renderErr != nil {
			c.logger.Println("failed render:", renderErr)
		} else {
			c.renderAndSend(content)
		}
	}
}

func (p *Telegram) OnFinished(ctx context.Context, sv pool.Instance) {}

func (c *Telegram) renderAndSend(message string) {
	msg := tgbotapi.NewMessage(0, message)
	msg.ParseMode = "markdown"
	for _, r := range c.Recipients {
		msg.ChatID = r
		_, err := c.bot.Send(msg)
		if err != nil {
			c.logger.Println("failed send message to", r, "due to", err)
		}
	}
}

func (a *Telegram) MergeFrom(other interface{}) (error) {
	b := other.(*Telegram)
	if a.Token == "" {
		a.Token = b.Token
	}
	if a.Token != b.Token {
		return errors.New("token are different")
	}
	a.withTemplate.resolvePath(a.workDir)
	b.withTemplate.resolvePath(b.workDir)
	if err := a.withTemplate.MergeFrom(&b.withTemplate); err != nil {
		return err
	}
	a.Recipients = append(a.Recipients, b.Recipients...)
	a.Services = append(a.Services, b.Services...)
	return nil
}
func (a *Telegram) Close() error { return nil }
func init() {
	registerPlugin("telegram", func(file string) PluginConfigNG {
		return &Telegram{workDir: filepath.Dir(file)}
	})
}
