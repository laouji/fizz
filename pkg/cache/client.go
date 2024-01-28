package cache

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

const KeyHitCountEndpoints = "hitcount:endpoints"

type Client struct {
	redis  *redis.Client
	logger logrus.FieldLogger
}

func NewClient(
	redisClient *redis.Client,
	logger logrus.FieldLogger,
) *Client {
	return &Client{redisClient, logger}
}

func (c *Client) IncrementEndpointHitCount(ctx context.Context, path string) error {
	_, err := c.redis.ZIncrBy(ctx, KeyHitCountEndpoints, 1, path).Result()
	return err
}

func (c *Client) GetEndpointHitCount(ctx context.Context, offset, limit int64) (
	results []redis.Z,
	err error,
) {
	args := redis.ZRangeArgs{
		Key:   KeyHitCountEndpoints,
		Start: offset,
		Stop:  limit - 1,
		Rev:   true,
	}
	return c.redis.ZRangeArgsWithScores(ctx, args).Result()
}
