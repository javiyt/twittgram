package bot

import (
	"strings"

	tb "gopkg.in/tucnak/telebot.v2"
)

type Bot struct {
	bot *tb.Bot
}

type BotOption func(b *Bot)

type botHandler struct {
	handlerFunc func(*tb.Message)
	help        string
	filters     []filterFunc
}

func WithTelegramBot(tb *tb.Bot) BotOption {
	return func(b *Bot) {
		b.bot = tb
	}
}

func NewBot(options ...BotOption) *Bot {
	b := &Bot{}

	for _, o := range options {
		o(b)
	}

	return b
}

func (b *Bot) Start() {
	b.setCommands()
	b.setUpHandlers()

	b.bot.Start()
}

func (b *Bot) Stop() {
	b.bot.Stop()
}

func (b *Bot) getHandlers() map[string]botHandler {
	return map[string]botHandler{
		"/start": {
			handlerFunc: b.handleStartCommand,
			help:        "Start a conversation with the bot",
			filters: []filterFunc{
				b.onlyPrivate,
			},
		},
		"/help": {
			handlerFunc: b.handleHelpCommand,
			help:        "Show help",
			filters: []filterFunc{
				b.onlyPrivate,
			},
		},
	}
}

func (b *Bot) setCommands() {
	var cmds []tb.Command
	for c, h := range b.getHandlers() {
		cmds = append(cmds, tb.Command{
			Text:        strings.Replace(c, "/", "", 1),
			Description: h.help,
		})
	}

	b.bot.SetCommands(cmds)
}

func (b *Bot) setUpHandlers() {
	for c, h := range b.getHandlers() {
		exec := h.handlerFunc

		for _, v := range h.filters {
			exec = v(exec)
		}

		b.bot.Handle(c, h.handlerFunc)
	}
}
