package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-redis/redismock/v9"
	"github.com/laouji/fizz/pkg/cache"
	"github.com/laouji/fizz/pkg/handler"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
)

type statsHandlerTestSuite struct {
	suite.Suite
	statsHandler http.HandlerFunc
	db           *redis.Client
	mock         redismock.ClientMock
}

func TestStatsHandler(t *testing.T) {
	suite.Run(t, new(statsHandlerTestSuite))
}

func (s *statsHandlerTestSuite) SetupTest() {
	logger := &logrus.Logger{
		//		Out: ioutil.Discard,
		Out: os.Stderr,
	}
	s.db, s.mock = redismock.NewClientMock()
	c := cache.NewClient(s.db, logger)
	s.statsHandler = handler.Stats(c, logger)
}

func (s *statsHandlerTestSuite) TestStats_OK() {
	req, err := http.NewRequest("GET", "/stats", nil)
	s.Require().NoError(err)

	expectedResult := []redis.Z{{Score: 25, Member: "param=val"}}
	args := redis.ZRangeArgs{
		Key:   cache.KeyHitCountEndpoints,
		Start: int64(0),
		Stop:  int64(0),
		Rev:   true,
	}
	s.mock.ExpectZRangeArgsWithScores(args).SetVal(expectedResult)
	rr := httptest.NewRecorder()

	s.statsHandler(rr, req)
	err = s.mock.ExpectationsWereMet()
	s.Require().NoError(err)

	s.Require().Equal(http.StatusOK, rr.Code)

	out := []handler.StatsOutput{}
	err = json.Unmarshal(rr.Body.Bytes(), &out)
	s.Require().NoError(err)
	s.Require().Len(out, 1)
	s.Equal(expectedResult[0].Member, out[0].Query)
	s.Equal(expectedResult[0].Score, out[0].HitCount)
}
