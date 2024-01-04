package config

type redis struct {
	Host   string
	Port   int
	Pass   string
	Prefix string
	DB     int
}

var Redis redis

func loadRedisConfig() {
	Redis = redis{
		Host:   GetString("redis.host", "127.0.0.1"),
		Port:   GetInt("redis.port", 6379),
		Pass:   GetString("redis.pass", ""),
		Prefix: GetString("redis.prefix", "blog:"),
		DB:     GetInt("redis.db", 0),
	}
}
