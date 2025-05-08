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

var ctx = context.Background()

var rdb = redis.NewClient(&redis.Options{
	Addr: os.Getenv("REDIS_URL"), // service name in docker-compose
})

// CacheBestFlightResponse stores best flight response for 30 sec, using the set of params as id
func CacheBestFlightResponse(searchCriteria string, data pkg.GetBestFlightOffersResponse) error {
	bodyBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return rdb.Set(ctx, searchCriteria, bodyBytes, 30*time.Second).Err()
}

// CacheBestFlightResponse restores best flight response, using the set of params as id
func GetCachedBestFlightResponse(searchCriteria string) (*pkg.GetBestFlightOffersResponse, error) {
	var response pkg.GetBestFlightOffersResponse

	bodyBytes, err := rdb.Get(ctx, searchCriteria).Bytes()
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
