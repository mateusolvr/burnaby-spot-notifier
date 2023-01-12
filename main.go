package main

import (
	"log"
	"time"

	"github.com/mateusolvr/web-scraper-go/domain/cache"
	"github.com/mateusolvr/web-scraper-go/domain/config"
	"github.com/mateusolvr/web-scraper-go/domain/email"
	readpage "github.com/mateusolvr/web-scraper-go/domain/read_page"
	"github.com/mateusolvr/web-scraper-go/domain/validation"
	"github.com/mateusolvr/web-scraper-go/internal/infrastructure/redis"
)

func main() {
	start := time.Now()
	log.Printf("Starting crawler at %s", time.Now())

	// Config
	configService := config.NewService()
	cfg := configService.GetConfig()

	// Storages...
	dbRedis, dbErr := redis.NewConnection(cfg.Redis.Url)
	cacheStorage := redis.NewCacheStorage(dbRedis)

	// Services...
	cacheService := cache.NewService(cacheStorage, cfg)
	emailService := email.NewService(cacheService, cfg)
	if dbErr != nil && cfg.Redis.Enabled {
		emailService.SendErrorEmail(dbErr)
		log.Fatal(dbErr)
	}
	validationService := validation.NewService(emailService)
	crawlerService := readpage.NewService(validationService, emailService, cacheService, cfg)

	crawlerService.InitializeCrawler()

	log.Printf("Took: %f secs\n", time.Since(start).Seconds())
}
