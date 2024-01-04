package config

type cache struct {
	Store string
}

var Cache cache

func loadCacheConfig() {
	Cache = cache{
		Store: GetString("cache.store", "redis"),
	}
}
