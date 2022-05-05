package main

import (
	"github.com/mateusolvr/web-scraper-go/domain/email"
	readpage "github.com/mateusolvr/web-scraper-go/domain/read_page"
	"github.com/mateusolvr/web-scraper-go/domain/validation"
)

func main() {
	validationService := validation.NewService()
	emailService := email.NewService()
	crawlerService := readpage.NewService(validationService, emailService)

	crawlerService.InitializeCrawler()
}
