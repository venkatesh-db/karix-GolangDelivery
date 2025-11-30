package configs

import "time"

type Config struct {
	BindAddr      string
	SQLiteDBPath  string
	RedisAddr     string
	RedisPassword string
	RedisDB       int
	CacheTTL      time.Duration
}

func LoadConfig() *Config {

	bind := getenv("BIND_ADDR", ":8080")
	sqlite := getenv("SQLITE_DB_PATH", "app.db")
	redisAddr := getenv("REDIS_ADDR", "localhost:6379")
	redisPassword := getenv("REDIS_PASSWORD", "")
	redisDB := getenvInt("REDIS_DB", 0)
	cacheTTL := getenvDuration("CACHE_TTL", 5*time.Minute)

	return &Config{
		BindAddr:      bind,
		SQLiteDBPath:  sqlite,
		RedisAddr:     redisAddr,
		RedisPassword: redisPassword,
		RedisDB:       redisDB,
		CacheTTL:      cacheTTL,
	}
}
