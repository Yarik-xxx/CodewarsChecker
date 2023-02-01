package app

import (
	"CodeWarsCheckerAnalysis/pkg/CodeWars"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"os"
	"sort"
	"time"
)

func (b *Bot) counter(challenges CodeWars.RespCompletedChallenges) map[string]map[string]int {
	result := make(map[string]map[string]int)
	result["overall"] = make(map[string]int)

	for _, challenge := range challenges.ListChallenges {
		// Получение ранга
		rank, err := b.st.GetChallenge(challenge.ID)

		// Обработка many requests
		for err != nil && err.Error() == "many requests" {
			b.appLogger.Printf("Many of requests. Sleep for 15 seconds (ID %s)", challenge.ID)
			time.Sleep(15 * time.Second)
			rank, err = b.st.GetChallenge(challenge.ID)
		}

		if err != nil {
			b.appLogger.Println("Error (skip)", err)
			continue
		}

		// Подсчет кат
		for _, lang := range challenge.CompletedLanguages {
			if _, ok := result[lang]; !ok {
				result[lang] = make(map[string]int)
			}
			result[lang][rank.Rank]++
			result["overall"][rank.Rank]++
		}
	}

	return result
}

func (b *Bot) generateMessage(userInfo CodeWars.RespUser, stat map[string]map[string]int) string {
	message := fmt.Sprintf("Username: <a href=\"https://www.codewars.com/users/%s\">%s</a>\nHonor: %d\nLeaderboard Position: %d\nRank: %s\nTotal Completed Kata: %d\n\nSummary:\n", userInfo.Username, userInfo.Username, userInfo.Honor, userInfo.LeaderboardPosition, userInfo.Ranks.Overall.Name, userInfo.CodeChallenges.TotalCompleted)

	for lang, statLang := range stat {
		if lang == "overall" {
			continue
		}

		message += "◉ " + lang + ":\n"

		var total int
		for _, count := range statLang {
			total += count
		}

		message += fmt.Sprintf("\t Total Completed Kata: %d\n", total)
		message += "\t Rank: " + userInfo.Ranks.Languages[lang].Name + "\n"
		message += fmt.Sprintf("\t Score: %d\n", userInfo.Ranks.Languages[lang].Score)
		message += "\t By rank:\n"

		message += generateSummary(statLang)
	}

	message += "◉ overall:\n"
	statLang := stat["overall"]
	message += "\t By rank:\n"
	message += generateSummary(statLang)

	return message
}

func (b *Bot) generateDocument(message string, chatId int64, userName string) error {
	content := []byte(message)
	tmpFile, err := os.CreateTemp("", "statistics.txt")
	if err != nil {
		b.appLogger.Println("Error:", err)
		return err
	}
	defer os.Remove(tmpFile.Name()) // очистка

	if _, err := tmpFile.Write(content); err != nil {
		b.appLogger.Println("Error:", err)
		return err
	}

	if err = tmpFile.Close(); err != nil {
		b.appLogger.Println("Error:", err)
		return err
	}

	tmpFile, err = os.Open(tmpFile.Name())
	if err != nil {
		b.appLogger.Println("Error:", err)
		return err
	}

	documentUpload := tgbotapi.NewDocument(chatId, tgbotapi.FileReader{
		Name:   fmt.Sprintf("%s_stat.txt", userName),
		Reader: tmpFile,
	})
	documentUpload.Caption = fmt.Sprintf("Whoa! %s, solved so many problems that his statistics didn't fit in one post. But I did not lose my head and collected everything in a file)", userName)
	b.bot.Send(documentUpload)

	return nil
}

func (b *Bot) sendMessage(text string, chatId int64, userName string) {
	msg := tgbotapi.NewMessage(chatId, "Oops... Something went wrong) Try again later")
	msg.ParseMode = "HTML"
	if len(text) < 4000 {
		msg.Text = text
		b.bot.Send(msg)
		return
	}

	// Отправка ответа в виде файла
	err := b.generateDocument(text, chatId, userName)
	if err != nil {
		b.appLogger.Println("Error:", err)
		b.bot.Send(msg)
	}
}

func generateSummary(statLang map[string]int) string {
	var message string

	kyuSort := make([]string, 0, 5)
	for kyu, _ := range statLang {
		kyuSort = append(kyuSort, kyu)
	}
	sort.Strings(kyuSort)

	for _, kyu := range kyuSort {
		count := statLang[kyu]
		message += fmt.Sprintf("\t\t▪ %s: %d\n", kyu, count)
	}

	return message
}
