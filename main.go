package main

import (
	"log"
	"time"

	"github.com/mateusolvr/web-scraper-go/domain/config"
	"github.com/mateusolvr/web-scraper-go/domain/email"
	readpage "github.com/mateusolvr/web-scraper-go/domain/read_page"
	"github.com/mateusolvr/web-scraper-go/domain/validation"
)

func main() {
	start := time.Now()
	log.Printf("Starting crawler at %s", time.Now())

	configService := config.NewService()
	cfg := configService.GetConfig()
	emailService := email.NewService(cfg)
	validationService := validation.NewService(emailService)
	crawlerService := readpage.NewService(validationService, emailService, cfg)

	crawlerService.InitializeCrawler()

	log.Printf("Took: %f secs\n", time.Since(start).Seconds())
}
