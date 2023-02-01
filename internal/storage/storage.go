package storage

import (
	"CodeWarsCheckerAnalysis/internal/config"
	"CodeWarsCheckerAnalysis/internal/models"
	"CodeWarsCheckerAnalysis/internal/storage/cache"
	"CodeWarsCheckerAnalysis/internal/storage/database"
	"CodeWarsCheckerAnalysis/pkg/CodeWars"
	"log"
	"os"
	"time"
)

type Storage struct {
	Db                    database.Database
	Cache                 cache.Cache
	chanNeedUpdate        chan string
	storageTimeChallenges time.Duration
	storageTimeMessages   time.Duration
	storageLogger         *log.Logger
}

func New() Storage {
	s := Storage{}

	file, err := os.OpenFile("log.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	s.storageLogger = log.New(file, "STORAGE: ", log.Ldate|log.Ltime|log.Lshortfile)
	s.storageLogger.Println("Parsing the config file")
	cfg := config.New("")

	s.storageTimeChallenges = cfg.Settings.StorageTimeChallenges
	s.storageTimeMessages = cfg.Settings.StorageTimeMessages
	s.chanNeedUpdate = make(chan string, 1024)

	s.storageLogger.Println("Database initialization")
	db, err := database.New(cfg)
	if err != nil {
		s.storageLogger.Fatal(err)
	}
	s.Db = db

	s.storageLogger.Println("Cache initialization")
	n, err := s.Cache.Init()
	if err != nil {
		s.storageLogger.Fatal(err)
	}

	s.storageLogger.Println("Storage init. Total uploaded", n, "challenges")
	return s
}

func (s *Storage) StartUpdate() {
	defer s.storageLogger.Println("Finished background update of challenges information")
	s.storageLogger.Println("Started background update of challenges information")

	for id := range s.chanNeedUpdate {
		_, err := s.UpdateThroughAPI(id)
		if err != nil {
			s.storageLogger.Println("Background update error:", err)
		}
		time.Sleep(700 * time.Millisecond)
	}

}

func (s *Storage) FindChallengeNeedUpdate() {
	defer s.storageLogger.Println("Finished background search for tasks needed in the update")
	s.storageLogger.Println("Started background search for tasks needed in the update")
	for {
		for _, challenge := range s.Cache.Challenges() {
			if s.isNeedUpdateChallenge(challenge) {
				s.chanNeedUpdate <- challenge.Id
			}
			time.Sleep(20 * time.Millisecond)
		}
	}
}

func (s *Storage) GetChallenge(id string) (models.Challenge, error) {
	resCache, ok := s.Cache.GetChallenge(id)
	if ok {
		return resCache, nil
	}

	return s.UpdateThroughAPI(id)
}

func (s *Storage) updateChallenge(id string, rank string) (models.Challenge, error) {
	challenge := models.Challenge{Id: id, Rank: rank, LastUpdate: time.Now()}
	if err := s.Db.InsertChallenge(challenge); err != nil {
		return challenge, err
	}

	s.Cache.UpdateChallenge(challenge)
	return challenge, nil
}

func (s *Storage) UpdateThroughAPI(id string) (models.Challenge, error) {
	resRequest, err := CodeWars.GetCodeChallenge(id)
	if err != nil {
		return models.Challenge{}, err
	}

	challenge, err := s.updateChallenge(resRequest.ID, resRequest.Rank.Name)
	if err != nil {
		return models.Challenge{}, err
	}

	return challenge, nil
}

func (s *Storage) isNeedUpdateChallenge(challenge models.Challenge) bool {
	difference := time.Now().Sub(challenge.LastUpdate)
	return difference > s.storageTimeChallenges
}

func (s *Storage) GetMessage(username string) (models.Message, bool) {
	resCache, ok := s.Cache.GetMessage(username)
	return resCache, ok && !s.isNeedUpdateMessage(resCache)
}

func (s *Storage) UpdateMessage(userName string, text string) {
	s.Cache.UpdateMessage(
		models.Message{
			UserName:   userName,
			Text:       text,
			LastUpdate: time.Now(),
		})
}

func (s *Storage) isNeedUpdateMessage(message models.Message) bool {
	difference := time.Now().Sub(message.LastUpdate)
	return difference > s.storageTimeMessages
}
