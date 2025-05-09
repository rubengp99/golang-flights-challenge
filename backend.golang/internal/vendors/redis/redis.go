package redis

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rubengp99/golang-flights-challenge/pkg"
)

// Service provides functionality to interact with redis
type Service struct {
	disabled bool
	ctx      context.Context
	rdb      *redis.Client
}

// NewRedisService creates a new redis service
func NewRedisService(disabled bool) Service {
	return Service{
		disabled: disabled,
		ctx:      context.Background(),
		rdb: redis.NewClient(&redis.Options{
			Addr: os.Getenv("REDIS_URL"), // service name in docker-compose
		}),
	}
}

// CacheBestFlightResponse stores best flight response for 30 sec, using the set of params as id
func (s Service) CacheBestFlightResponse(searchCriteria string, data pkg.GetBestFlightOffersResponse) error {
	if s.disabled {
		return nil
	}

	bodyBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return s.rdb.Set(s.ctx, searchCriteria, bodyBytes, 30*time.Second).Err()
}

// CacheBestFlightResponse restores best flight response, using the set of params as id
func (s Service) GetCachedBestFlightResponse(searchCriteria string) (*pkg.GetBestFlightOffersResponse, error) {
	if s.disabled {
		return nil, nil
	}

	var response pkg.GetBestFlightOffersResponse

	bodyBytes, err := s.rdb.Get(s.ctx, searchCriteria).Bytes()
	// means not found
	if err == redis.Nil {
		log.Println("no redis cache found")
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(bodyBytes, &response); err != nil {
		return nil, err
	}

	return &response, nil
}
