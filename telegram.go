package main

import (
	"github.com/reddec/container"
	"text/template"
	"log"
	"os"
	"bytes"
	"gopkg.in/telegram-bot-api.v4"
	"github.com/pkg/errors"
	"time"
)

type Telegram struct {
	Token      string   `yaml:"token"`
	Recipients []int64  `yaml:"recipients"`
	Services   []string `yaml:"services"`
	Template   string   `yaml:"template"`

	servicesSet map[string]bool     `yaml:"-"`
	templateBin *template.Template  `yaml:"-"`
	logger      *log.Logger         `yaml:"-"`
	bot         *tgbotapi.BotAPI    `yaml:"-"`
	hostname    string
}

func (c *Telegram) Prepare() error {
	c.servicesSet = make(map[string]bool)
	for _, srv := range c.Services {
		c.servicesSet[srv] = true
	}
	t, err := template.New("").Parse(c.Template)
	if err != nil {
		return err
	}
	c.templateBin = t
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

func (c *Telegram) renderAndSend(params map[string]interface{}) {
	message := &bytes.Buffer{}
	renderErr := c.templateBin.Execute(message, params)
	if renderErr != nil {
		c.logger.Println("failed render:", renderErr, "; params:", params)
	} else {
		msg := tgbotapi.NewMessage(0, message.String())
		msg.ParseMode = "markdown"
		for _, r := range c.Recipients {
			msg.ChatID = r
			_, err := c.bot.Send(msg)
			if err != nil {
				c.logger.Println("failed send message to", r, "due to", err)
			}
		}
	}
}

func (c *Telegram) Spawned(runnable container.Runnable, id container.ID) {
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

func mergeTelegram(a *Telegram, b *Telegram) (*Telegram, error) {
	if a == nil {
		return b, nil
	}
	if b == nil {
		return a, nil
	}
	if a.Token == "" {
		a.Token = b.Token
	}
	if b.Token == "" {
		b.Token = a.Token
	}
	if a.Token != b.Token {
		return nil, errors.New("token are different")
	}
	if a.Template == "" {
		a.Template = b.Template
	}
	if b.Template == "" {
		b.Template = a.Template
	}
	if a.Template != b.Template {
		return nil, errors.New("different templates")
	}
	a.Recipients = append(a.Recipients, b.Recipients...)
	a.Services = append(a.Services, b.Services...)
	return a, nil
}
