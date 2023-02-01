package main

import (
	"CodeWarsCheckerAnalysis/internal/config"
	"CodeWarsCheckerAnalysis/internal/models"
	"CodeWarsCheckerAnalysis/internal/storage"
	"CodeWarsCheckerAnalysis/internal/storage/database"
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	// БД
	cfg := config.New("")
	db, err := database.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	err = db.InitTables()
	if err != nil {
		log.Fatal(err)
	}

	// Файл с ифнормацией о катах
	file, err := os.Open("/home/codewarsbot/CodeWarsCheckerAnalysis/db.txt")

	if err != nil {
		log.Fatalf("Error when opening file: %s", err)
	}
	defer file.Close()

	// Получение информации
	fileScanner := bufio.NewScanner(file)

	for fileScanner.Scan() {
		row := fileScanner.Text()
		data := strings.Split(row, ",")

		err = db.InsertChallenge(models.Challenge{
			Id:         data[0],
			Rank:       data[1],
			LastUpdate: time.Now(),
		})
		if err != nil {
			log.Fatal(err)
		}
	}
}

func initFile() {
	st := storage.New()

	// Файл для кат
	f, err := os.OpenFile("cmd/tools/filler/db.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	for _, c := range st.Cache.Challenges() {
		if _, err = f.WriteString(fmt.Sprintf("%s,%s,%s\n", c.Id, c.Rank, c.LastUpdate)); err != nil {
			log.Fatal(err)
		}
	}
}
