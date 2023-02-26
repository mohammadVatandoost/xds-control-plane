package redis_test

import (
	"context"
	"fmt"
	"github.com/alicebob/miniredis/v2"
	"github.com/elliotchance/redismock/v8"
	"github.com/go-redis/redis/v8"
	internalRedis "github.com/mohammadVatandoost/xds-conrol-plane/pkg/redis"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type ClientTestSuite struct {
	suite.Suite
}

func (s *ClientTestSuite) TestShouldSetKeyCorrectly() {
	redisMock, _ := getRedisMock()
	client := internalRedis.NewClient(redisMock)

	_, err := client.Set(context.Background(), "key", "value", time.Hour)
	result, _ := redisMock.Get(context.Background(), "key").Result()

	s.Nil(err)
	s.Equal("value", result)
}

func (s *ClientTestSuite) TestShouldSetKeyCorrectlyWithCircuitBreaker() {
	redisMock, _ := getRedisMock()
	client := internalRedis.NewClient(redisMock)

	_, err := client.Set(context.Background(), "key", "value", time.Hour)
	result, _ := redisMock.Get(context.Background(), "key").Result()

	s.Nil(err)
	s.Equal("value", result)
}

func (s *ClientTestSuite) TestShouldGetKeyCorrectly() {
	redisMock, _ := getRedisMock()
	redisMock.Set(context.Background(), "key", "value", time.Hour)
	client := internalRedis.NewClient(redisMock)

	result, err := client.Get(context.Background(), "key")

	s.Nil(err)
	s.Equal("value", result)
}

func (s *ClientTestSuite) TestShouldGetKeyCorrectlyWithCircuitBreaker() {
	redisMock, _ := getRedisMock()
	redisMock.Set(context.Background(), "key", "value", time.Hour)
	client := internalRedis.NewClient(redisMock)

	result, err := client.Get(context.Background(), "key")

	s.Nil(err)
	s.Equal("value", result)
}

func (s *ClientTestSuite) TestShouldReturnNotFoundIfKeyDoesNotExist() {
	redisMock, _ := getRedisMock()
	client := internalRedis.NewClient(redisMock)

	result, err := client.Get(context.Background(), "key")

	s.Empty(result)
	s.Equal(internalRedis.ErrNotFound, err)
}

func (s *ClientTestSuite) TestShouldReturnUnavailableForTimeout() {
	redisMock := getUnavailableRedisMock()
	client := internalRedis.NewClient(redisMock)

	result, err := client.Get(context.Background(), "key")
	s.Empty(result)
	s.Equal(internalRedis.ErrNotAvailable, err)

	result, err = client.Set(context.Background(), "key", "value", time.Hour)
	s.Empty(result)
	s.Equal(internalRedis.ErrNotAvailable, err)
}

func (s *ClientTestSuite) TestShouldReturnUnavailableForCircuitBreakerError() {
	redisMock := getUnavailableRedisMock()
	client := internalRedis.NewClient(redisMock)

	result, err := client.Get(context.Background(), "key")
	s.Empty(result)
	s.Equal(internalRedis.ErrNotAvailable, err)

	result, err = client.Set(context.Background(), "key", "value", time.Hour)
	s.Empty(result)
	s.Equal(internalRedis.ErrNotAvailable, err)
}

func (s *ClientTestSuite) TestShouldReturnErrorIfRedisFailed() {
	redisMock := getFaultyRedisMock()
	client := internalRedis.NewClient(redisMock)

	result, err := client.Get(context.Background(), "key")
	s.Empty(result)
	s.Equal("unknown redis error", err.Error())

	result, err = client.Set(context.Background(), "key", "value", time.Hour)
	s.Empty(result)
	s.Equal("unknown redis error", err.Error())
}

func TestRedisTestSuite(t *testing.T) {
	suite.Run(t, new(ClientTestSuite))
}

func getRedisMock() (*redis.Client, *miniredis.Miniredis) {
	mr, err := miniredis.Run()
	if err != nil {
		return nil, nil
	}
	m := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	return m, mr
}

func getUnavailableRedisMock() redis.Cmdable {
	redisMock, _ := getRedisMock()
	clientMock := redismock.NewNiceMock(redisMock)

	clientMock.On("Get", mock.Anything, mock.Anything).Return(
		redis.NewStringResult("", context.DeadlineExceeded),
	)
	clientMock.On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(
		redis.NewStatusResult("", context.DeadlineExceeded),
	)

	return clientMock
}

func getFaultyRedisMock() redis.Cmdable {
	redisMock, _ := getRedisMock()
	clientMock := redismock.NewNiceMock(redisMock)

	clientMock.On("Get", mock.Anything, mock.Anything).Return(
		redis.NewStringResult("", fmt.Errorf("unknown redis error")),
	)
	clientMock.On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(
		redis.NewStatusResult("", fmt.Errorf("unknown redis error")),
	)

	return clientMock
}
