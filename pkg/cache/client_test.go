package cache_test

import (
	"context"
	"io/ioutil"
	"testing"

	"github.com/go-redis/redismock/v9"
	"github.com/laouji/fizz/pkg/cache"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
)

type clientTestSuite struct {
	suite.Suite
	db     *redis.Client
	mock   redismock.ClientMock
	logger logrus.FieldLogger
}

func TestClient(t *testing.T) {
	suite.Run(t, new(clientTestSuite))
}

func (s *clientTestSuite) SetupTest() {
	s.db, s.mock = redismock.NewClientMock()
	s.logger = &logrus.Logger{
		Out: ioutil.Discard,
	}
}

func (s *clientTestSuite) TearDownTest() {
	err := s.mock.ExpectationsWereMet()
	s.Require().NoError(err)
}

func (s *clientTestSuite) TestGetEndpointHitCount() {
	client := cache.NewClient(s.db, s.logger)
	var offset int64 = 0
	var limit int64 = 10
	results := []redis.Z{{Score: 1, Member: "somestr"}}
	args := redis.ZRangeArgs{
		Key:   cache.KeyHitCountEndpoints,
		Start: offset,
		Stop:  limit - 1,
		Rev:   true,
	}
	s.mock.ExpectZRangeArgsWithScores(args).SetVal(results)

	r, err := client.GetEndpointHitCount(context.Background(), offset, limit)
	s.Require().NoError(err)
	s.Equal(results, r)
}

func (s *clientTestSuite) TestIncrementEndpointHitCount() {
	client := cache.NewClient(s.db, s.logger)
	path := "somepath"
	s.mock.ExpectZIncrBy(cache.KeyHitCountEndpoints, 1, path).SetVal(float64(0))

	err := client.IncrementEndpointHitCount(context.Background(), path)
	s.Require().NoError(err)
}
