package main

import (
	"github.com/mateusolvr/web-scraper-go/domain/config"
	"github.com/mateusolvr/web-scraper-go/domain/email"
	readpage "github.com/mateusolvr/web-scraper-go/domain/read_page"
	"github.com/mateusolvr/web-scraper-go/domain/validation"
)

func main() {
	configService := config.NewService()
	cfg := configService.GetConfig()

	validationService := validation.NewService()
	emailService := email.NewService()
	crawlerService := readpage.NewService(validationService, emailService, cfg)

	crawlerService.InitializeCrawler()
}
