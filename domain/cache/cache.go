package cache

import (
	"github.com/mateusolvr/web-scraper-go/domain"
)

type service struct {
	cacheStorage domain.CacheStorage
	cfg          domain.Config
}

func NewService(cacheStorage domain.CacheStorage, cfg domain.Config) *service {
	return &service{
		cacheStorage: cacheStorage,
		cfg:          cfg,
	}
}

func (s *service) GetActivitiesWithoutCache(activities []domain.Activity) ([]domain.Activity, error) {
	var filteredAct []domain.Activity

	if !s.cfg.Redis.Enabled {
		return activities, nil
	}

	for _, v := range activities {
		keyValueStr, err := s.cacheStorage.GetKey(v.ActKeyCache)
		if err != nil && err.Error() != "redis: nil" {
			return nil, err
		}
		if keyValueStr == "" {
			filteredAct = append(filteredAct, v)
		}
	}

	return filteredAct, nil
}

func (s *service) AddActivitiesCache(activities []domain.Activity) error {
	if s.cfg.Redis.Enabled {
		for _, v := range activities {
			err := s.cacheStorage.SetKey(v.ActKeyCache, "true", s.cfg.Redis.ExpireMinutes)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *service) DelKey(key string) (int64, error) {
	if s.cfg.Redis.Enabled {
		return s.cacheStorage.DelKey(key)
	}
	return 0, nil
}

func (s *service) SetKey(key, value string, expiration int) error {
	if s.cfg.Redis.Enabled {
		s.cacheStorage.SetKey(key, value, expiration)
	}
	return nil
}

func (s *service) CheckErrorCache(err error) (bool, error) {
	if !s.cfg.Redis.Enabled {
		return false, nil
	}

	keyValueStr, cacheErr := s.cacheStorage.GetKey(err.Error())
	if cacheErr != nil && cacheErr.Error() != "redis: nil" {
		return false, cacheErr
	}
	if cacheErr == nil || keyValueStr != "" {
		return true, nil
	}
	return false, nil
}
