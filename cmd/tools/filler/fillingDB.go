package main

import (
	"CodeWarsCheckerAnalysis/internal/storage"
	"CodeWarsCheckerAnalysis/pkg/CodeWars"
	"bufio"
	"log"
	"os"
	"time"
)

func main() {
}

func filling() {
	// Инициализация хранилища
	st := storage.New()

	// Файл с ID
	file, err := os.Open("cmd/tools/id.txt")
	if err != nil {
		log.Fatalf("Error when opening file: %s", err)
	}
	defer file.Close()

	// Получение ID
	fileScanner := bufio.NewScanner(file)

	for fileScanner.Scan() {
		id := fileScanner.Text()
		// Проверка на наличие информации
		_, ok := st.Cache.GetChallenge(id)
		if ok {
			//log.Default().Println("Cache:", c.Id, "-", c.Rank)
			continue
		}

		// Получение информации, если ее нет
		_, err := st.UpdateThroughAPI(id)
		if err != nil {
			if err.Error() == "many requests" {
				time.Sleep(1 * time.Minute)
			}
			log.Default().Println("Get error:", err)
			continue
		}
		//log.Default().Println("Request:", c.Id, "-", c.Rank)
		time.Sleep(400 * time.Millisecond)
	}
}

func getId() {
	// Файл для ID
	f, err := os.OpenFile("id.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// Файл с именами из топа
	file, err := os.Open("names.txt")

	if err != nil {
		log.Fatalf("Error when opening file: %s", err)
	}
	defer file.Close()

	// Получение ID
	fileScanner := bufio.NewScanner(file)
	ids := make(map[string]interface{})

	for fileScanner.Scan() {
		// Запрос списка кат
		userName := fileScanner.Text()
		r, err := CodeWars.GetCompletedChallenges(userName)
		if err != nil {
			log.Default().Println("error:", err)
			if err.Error() == "many requests" {
				time.Sleep(1 * time.Minute)
			}
			continue
		}

		// Обработка ID из списка
		for _, challengeComp := range r.ListChallenges {
			_, ok := ids[challengeComp.ID]
			if !ok {
				ids[challengeComp.ID] = nil
				if _, err = f.WriteString(challengeComp.ID + "\n"); err != nil {
					log.Fatal(err)
				}
			}
		}
		log.Default().Println(userName, len(ids))
		time.Sleep(30 * time.Second)
	}
}
