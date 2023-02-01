package cache

import (
	"CodeWarsCheckerAnalysis/internal/config"
	"CodeWarsCheckerAnalysis/internal/models"
	"CodeWarsCheckerAnalysis/internal/storage/database"
	"sync"
)

type Cache struct {
	challenges   map[string]models.Challenge
	messages     map[string]models.Message
	mtChallenges *sync.Mutex
	mtMessages   *sync.Mutex
}

func (c *Cache) Init() (int, error) {
	cfg := config.New("")
	c.mtChallenges = &sync.Mutex{}
	c.mtMessages = &sync.Mutex{}

	db, err := database.New(cfg)
	if err != nil {
		return 0, err
	}

	challenges, err := db.GetAllChallenges()
	if err != nil {
		return 0, err
	}

	c.challenges = challenges
	c.messages = make(map[string]models.Message)

	return len(challenges), nil
}

func (c *Cache) GetChallenge(id string) (models.Challenge, bool) {
	c.mtChallenges.Lock()
	defer c.mtChallenges.Unlock()
	res, ok := c.challenges[id]
	return res, ok
}

func (c *Cache) UpdateChallenge(challenge models.Challenge) {
	c.mtChallenges.Lock()
	defer c.mtChallenges.Unlock()
	c.challenges[challenge.Id] = challenge
}

func (c *Cache) Challenges() map[string]models.Challenge {
	c.mtChallenges.Lock()
	defer c.mtChallenges.Unlock()

	copyChallenges := make(map[string]models.Challenge, len(c.challenges))
	for id, challenge := range c.challenges {
		copyChallenges[id] = challenge
	}

	return copyChallenges
}

func (c *Cache) GetMessage(userName string) (models.Message, bool) {
	c.mtMessages.Lock()
	defer c.mtMessages.Unlock()
	res, ok := c.messages[userName]
	return res, ok
}

func (c *Cache) UpdateMessage(message models.Message) {
	c.mtMessages.Lock()
	defer c.mtMessages.Unlock()
	c.messages[message.UserName] = message
}




