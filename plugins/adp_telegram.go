package plugins

import (
	"github.com/reddec/container"
	"log"
	"os"
	"gopkg.in/telegram-bot-api.v4"
	"errors"
	"path/filepath"
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

func (c *Telegram) Prepare() error {
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

func (c *Telegram) Stopped(runnable container.Runnable, id container.ID, err error) {
	if c.servicesSet[runnable.Label()] {
		content, renderErr := c.renderDefault("stopped", string(id), runnable.Label(), err, c.logger)
		if renderErr != nil {
			c.logger.Println("failed render:", renderErr)
		} else {
			c.renderAndSend(content)
		}
	}
}

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

func (c *Telegram) Spawned(runnable container.Runnable, id container.ID) {
	if c.servicesSet[runnable.Label()] {
		content, renderErr := c.renderDefault("spawned", string(id), runnable.Label(), nil, c.logger)
		if renderErr != nil {
			c.logger.Println("failed render:", renderErr)
		} else {
			c.renderAndSend(content)
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

func init() {
	registerPlugin("telegram", func(file string) PluginConfig {
		return &Telegram{workDir: filepath.Dir(file)}
	})
}
