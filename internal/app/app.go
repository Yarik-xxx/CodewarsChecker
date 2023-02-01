package app

import (
	"CodeWarsCheckerAnalysis/internal/config"
	"CodeWarsCheckerAnalysis/internal/storage"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
)

type Bot struct {
	cfg              config.Config
	st               storage.Storage
	bot              *tgbotapi.BotAPI
	queue            chan *tgbotapi.Update
	appLogger        *log.Logger
	statisticsLogger *log.Logger
}

func (b *Bot) Start() {
	// logging
	file, err := os.OpenFile("log.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	b.appLogger = log.New(file, "APP: ", log.Ldate|log.Ltime|log.Lshortfile)
	b.statisticsLogger = log.New(file, "STATISTICS: ", log.Ldate|log.Ltime|log.Lshortfile)

	b.initSettings(true)

	// bot fall protection
	defer func() {
		if err := recover(); err != nil {
			b.bot.StopReceivingUpdates()
			b.appLogger.Println("PANIC:", err)
			panic(err)
		}
	}()

	// running background processes
	go b.st.FindChallengeNeedUpdate()
	go b.st.StartUpdate()
	go b.startQueue()

	// receiving updates
	updates := b.getUpdatesChan(60)
	for update := range updates {
		if update.Message != nil {
			b.statisticsLogger.Printf("Request (%d): %s\n", update.Message.Chat.ID, update.Message.Text)
			switch update.Message.Text {
			case "/start":
				b.handlerStart(&update)

			default:
				go func(u tgbotapi.Update) {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "The request has been added to the queue. Statistics will be ready in a couple of minutes")
					b.bot.Send(msg)
					b.queue <- &u
				}(update)
			}
		}
	}
}

func (b *Bot) initSettings(debug bool) {
	b.cfg = config.New("")
	b.st = storage.New()
	b.queue = make(chan *tgbotapi.Update)

	bot, err := tgbotapi.NewBotAPI(b.cfg.Bot.Token)
	if err != nil {
		b.appLogger.Fatal(err)
	}
	b.bot = bot

	bot.Debug = debug

	b.appLogger.Printf("Authorized on account %s\n", bot.Self.UserName)
}

func (b *Bot) startQueue() {
	b.appLogger.Println("Started queue handler")
	defer b.appLogger.Println("Closing a queue handler")

	for update := range b.queue {
		b.handlerStatistics(update)
	}
}

func (b *Bot) getUpdatesChan(timeout int) tgbotapi.UpdatesChannel {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = timeout

	updates := b.bot.GetUpdatesChan(u)

	return updates
}
