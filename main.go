package main

import (
	"log"
	"time"

	"github.com/mateusolvr/burnaby-spot-notifier/domain/cache"
	"github.com/mateusolvr/burnaby-spot-notifier/domain/config"
	"github.com/mateusolvr/burnaby-spot-notifier/domain/email"
	readapi "github.com/mateusolvr/burnaby-spot-notifier/domain/read_api"
	"github.com/mateusolvr/burnaby-spot-notifier/domain/validation"
	"github.com/mateusolvr/burnaby-spot-notifier/internal/infrastructure/redis"
)

func main() {
	start := time.Now()
	log.Printf("Starting notifier at %s", time.Now())

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
	apiService := readapi.NewService(validationService, emailService, cacheService, cfg)

	apiService.Initialize()

	log.Printf("Took: %f secs\n", time.Since(start).Seconds())
}
