package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"time"
)

type Config struct {
	Bot struct {
		Token   string `yaml:"Token"`
		AdminID int    `yaml:"AdminID"`
	}
	Settings struct {
		StorageTimeChallenges time.Duration `yaml:"StorageTimeChallenges"`
		StorageTimeMessages   time.Duration `yaml:"StorageTimeMessages"`
	}
	Database struct {
		DBName   string `yaml:"DBName"`
		UserName string `yaml:"UserName"`
		Password string `yaml:"Password"`
		Host     string `yaml:"Host"`
		Port     int    `yaml:"Port"`
	}
}

func New(pathToFile string) Config {
	c := Config{}

	if pathToFile == "" {
		pathToFile = "/home/codewarsbot/CodeWarsCheckerAnalysis/config.yaml"
	}

	f, err := os.ReadFile(pathToFile)
	if err != nil {
		log.Fatal(err)
	}

	if err := yaml.Unmarshal(f, &c); err != nil {
		log.Fatal(err)
	}

	return c
}

func (c *Config) GetConnStr() string {
	return fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%d sslmode=disable",
		c.Database.UserName, c.Database.Password, c.Database.DBName, c.Database.Host, c.Database.Port,
	)
}
