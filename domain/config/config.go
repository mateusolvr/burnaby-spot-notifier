package config

import (
	"log"
	"os"

	"github.com/mateusolvr/web-scraper-go/domain"
	"gopkg.in/yaml.v2"
)

type service struct {
}

func NewService() *service {
	return &service{}
}

func (s *service) GetConfig() domain.Config {
	f, err := os.Open("config.yml")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	var cfg domain.Config

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	return cfg
}
