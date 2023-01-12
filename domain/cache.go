package domain

type CacheStorage interface {
	SetKey(key, value string, expiration int) error
	GetKey(key string) (string, error)
	DelKey(key string) (int64, error)
}

type CacheService interface {
	GetActivitiesWithoutCache(activities []Activity) ([]Activity, error)
	AddActivitiesCache(activities []Activity) error
	DelKey(key string) (int64, error)
	SetKey(key, value string, expiration int) error
	CheckErrorCache(err error) (bool, error)
}
