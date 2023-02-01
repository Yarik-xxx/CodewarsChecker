package app

import (
	"CodeWarsCheckerAnalysis/pkg/CodeWars"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"time"
)

func (b *Bot) handlerStart(update *tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Hey! Send me a user's nickname in CodeWars and I'll collect his account statistics")

	b.bot.Send(msg)
}

func (b *Bot) handlerStatistics(update *tgbotapi.Update) {
	username := update.Message.Text
	chatId := update.Message.Chat.ID

	// Попытка получения из кэша
	cacheMessage, ok := b.st.GetMessage(username)
	if ok {
		b.sendMessage(cacheMessage.Text, chatId, cacheMessage.UserName)
		return
	}

	msg := tgbotapi.NewMessage(chatId, "Oops... Something went wrong) Try again later")

	// Информация об аккаунте
	userInfo, err := CodeWars.GetUser(username)
	if err != nil {
		if err.Error() == "not found" {
			msg.Text = fmt.Sprintf("Пользователь %s не найден", username)
		}

		b.appLogger.Println("Error:", err)

		if err.Error() == "many requests" {
			time.Sleep(30 * time.Second)

			go func() {
				b.queue <- update
			}()

		}
		b.bot.Send(msg)
		return
	}

	// Список выполненных кат
	listChallenges, err := CodeWars.GetCompletedChallenges(username)
	if err != nil {
		b.appLogger.Println("Error:", err)

		if err.Error() == "many requests" {
			time.Sleep(30 * time.Second)
			go func(u tgbotapi.Update) {
				b.queue <- update
			}(*update)
			return
		}

		b.bot.Send(msg)
		return
	}

	// Статистика выполненных кат
	stat := b.counter(listChallenges)

	// Генерация сообщения
	text := b.generateMessage(userInfo, stat)

	// Обновление в кэше
	b.st.UpdateMessage(username, text)

	// Отправка сообщения (файл | сообщение)
	b.sendMessage(text, chatId, username)
}
