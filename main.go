package main

import (
	readpage "github.com/mateusolvr/web-scraper-go/domain/read_page"
	"github.com/mateusolvr/web-scraper-go/domain/validation"
)

func main() {
	validationService := validation.NewService()
	crawlerService := readpage.NewService(validationService)

	crawlerService.InitializeCrawler()
}
