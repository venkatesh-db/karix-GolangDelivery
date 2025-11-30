package service

import (
	"context"
	"strconv"
	"time"

	"karix.com/monolith/caching"
	"karix.com/monolith/repos"
	"karix.com/monolith/schemas"
)

type UserService struct {
	repo  *repos.UserRepo
	cache *caching.Cacheservice
	// Add fields as necessary, e.g., database connections, configurations, etc.
}

func NewUserService(repo *repos.UserRepo, cache *caching.Cacheservice) *UserService {
	return &UserService{
		repo:  repo,
		cache: cache,
	}
}

// Add methods for UserService to handle user-related operations.

func (s *UserService) CreateUser(ctx context.Context, name, email string) (*schemas.User, error) {

	user, err := s.repo.CreateUser(ctx, name, email)

	if err != nil {
		return nil, err
	}

	key := userCachekey(int64(user.ID))
	// Cache the newly created user
	s.cache.SetUser(ctx, key, user)

	return user, nil

}

func userCachekey(id int64) string {
	return "user:" + strconv.FormatInt(id, 10)
}

func (s *UserService) GetUser(ctx context.Context, id int64) (*schemas.User, error) {

	key := userCachekey(id)

	// 1. Try caches with a short context
	ctxCache, cancel := context.WithTimeout(ctx, 500*time.Millisecond)

	defer cancel()

	var uq *schemas.User

	if uq, ok, err := s.cache.GetUser(ctxCache, key); err == nil && ok {
		return uq, nil

	}

	return uq, nil
}
