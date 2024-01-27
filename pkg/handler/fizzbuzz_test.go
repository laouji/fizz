package handler_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/laouji/fizz/pkg/handler"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
)

type handlerTestSuite struct {
	suite.Suite
	fizzBuzzHandler http.HandlerFunc
}

func TestHandler(t *testing.T) {
	suite.Run(t, new(handlerTestSuite))
}

func (s *handlerTestSuite) SetupTest() {
	logger := &logrus.Logger{
		Out: ioutil.Discard,
	}
	s.fizzBuzzHandler = handler.FizzBuzz(logger)
}

func (s *handlerTestSuite) TestFizzBuzz_MissingParams() {
	req, err := http.NewRequest("GET", "/", nil)
	s.Require().NoError(err)

	rr := httptest.NewRecorder()

	s.fizzBuzzHandler(rr, req)
	s.Equal(http.StatusBadRequest, rr.Code)
	s.Regexp("required", rr.Body.String())
}

func (s *handlerTestSuite) TestFizzBuzz_InvalidStr() {
	req, err := http.NewRequest("GET", "/?int1=3&int2=5&str1=ASCII意外の文字&str2=valid&limit=30", nil)
	s.Require().NoError(err)

	rr := httptest.NewRecorder()

	s.fizzBuzzHandler(rr, req)
	s.Equal(http.StatusBadRequest, rr.Code)
	s.Regexp("str1", rr.Body.String())
	s.Regexp("printascii", rr.Body.String())
}

func (s *handlerTestSuite) TestFizzBuzz_InvalidInt() {
	req, err := http.NewRequest("GET", "/?int1=3&int2=five&str1=somestr&str2=valid&limit=30", nil)
	s.Require().NoError(err)

	rr := httptest.NewRecorder()

	s.fizzBuzzHandler(rr, req)
	s.Equal(http.StatusBadRequest, rr.Code)
	s.Regexp("int2", rr.Body.String())
	s.Regexp("number", rr.Body.String())
}

func (s *handlerTestSuite) TestFizzBuzz_LimitOutOfRange() {
	req, err := http.NewRequest("GET", "/?int1=3&int2=6&str1=fizz&str2=buzz&limit=10000000", nil)
	s.Require().NoError(err)

	rr := httptest.NewRecorder()

	s.fizzBuzzHandler(rr, req)
	s.Equal(http.StatusBadRequest, rr.Code)
	s.Regexp("limit", rr.Body.String())
	s.Regexp("max", rr.Body.String())
}

func (s *handlerTestSuite) TestFizzBuzz_StringOutOfRange() {
	u := url.URL{}
	q := u.Query()
	q.Set("str1", ";.")
	q.Set("str2", `As Athos and Porthos had foreseen, at the expiration of a half hour, D'Artagnan returned. He had again missed his man, who had disappeared as if by enchantment.`)
	q.Set("int1", "1")
	q.Set("int2", "3")
	q.Set("limit", "25")

	req, err := http.NewRequest("GET", "?"+q.Encode(), nil)
	s.Require().NoError(err)

	rr := httptest.NewRecorder()

	s.fizzBuzzHandler(rr, req)
	s.Equal(http.StatusBadRequest, rr.Code)
	s.Regexp("str2", rr.Body.String())
	s.Regexp("max", rr.Body.String())
}

func (s *handlerTestSuite) TestFizzBuzz_OK() {
	req, err := http.NewRequest("GET", "/?int1=3&int2=5&str1=fizz&str2=buzz&limit=16", nil)
	s.Require().NoError(err)

	rr := httptest.NewRecorder()

	s.fizzBuzzHandler(rr, req)
	s.Equal(http.StatusOK, rr.Code)
	var res = []string{}
	err = json.Unmarshal(rr.Body.Bytes(), &res)
	s.Require().NoError(err)
	s.Len(res, 16)
}
