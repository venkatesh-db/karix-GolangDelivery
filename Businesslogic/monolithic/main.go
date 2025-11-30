package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"karix.com/monolith/caching"
	"karix.com/monolith/configs"
	"karix.com/monolith/handlers"
	"karix.com/monolith/repos"
	"karix.com/monolith/schemas"
	"karix.com/monolith/service"
)

/*

 curl -s -X POST http://localhost:8080/api/users \
-H "Content-Type: application/json" \
-d '{"username":"tiktock boy","email":"alice@ticktock.com" , "password":"smiles12"}'


curl -s http://localhost:8080/api/users/1


*/

func main() {

	cfg := configs.LoadConfig() // from env with defaults

	// Setup DB (SQLite)
	db, err := sql.Open("sqlite3", cfg.SQLiteDBPath)
	if err != nil {
		log.Fatalf("open sqlite: %v", err)
	}
	db.SetMaxOpenConns(1) // sqlite concurrent writes are limited

	/*
		if err := migrate(db); err != nil {
			log.Fatalf("migrate: %v", err)
		}
	*/

	// Redis client (optional — app still works with redis down via local cache)
	redisOpt := &redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	}
	rcli := redis.NewClient(redisOpt)
	// do a quick ping with timeout
	//	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	//defer cancel()
	if err := rcli.Ping().Err(); err != nil {
		log.Printf("redis ping failed (continuing with local cache): %v", err)
	}

	schemas.MigrateDB(db)

	// Create caches & repository & handlers
	localCache := caching.NewLocalCache(5 * time.Minute) // small in-memory cache
	repo := repos.NewUserRepo(db)
	cache := caching.NewCacheService(rcli, localCache, cfg.CacheTTL)
	svc := service.NewUserService(repo, cache)
	h := handlers.NewHTTPHandler(svc)

	// Gin setup
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	api := r.Group("/api")
	{
		api.POST("/users", h.CreateUser)
		api.GET("/users/:id", h.GetUser)
	}

	srv := &http.Server{
		Addr:    cfg.BindAddr,
		Handler: r,
	}

	// Start server
	go func() {
		log.Printf("listening on %s", cfg.BindAddr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server listen: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutting down...")

	ctxShutdown, cancelShutdown := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelShutdown()

	if err := srv.Shutdown(ctxShutdown); err != nil {
		log.Fatalf("server shutdown: %v", err)
	}
	if err := db.Close(); err != nil {
		log.Printf("db close: %v", err)
	}
	if err := rcli.Close(); err != nil {
		log.Printf("redis close: %v", err)
	}
	log.Println("graceful exit")
}
