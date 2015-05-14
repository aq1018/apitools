package errormapper

import (
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
)

func (s *mapperTestSuite) TestMappError() {
	s.Mapper.Register(func(err error) ErrorMessage {
		return nil
	})

	s.Nil(s.Mapper.MapError(testErr))

	s.Mapper.Register(func(err error) (message ErrorMessage) {
		if err == testErr {
			message = &errorMessage{500, testErr.Error()}
		}
		return
	})

	m := s.Mapper.MapError(testErr)

	message, ok := m.(*errorMessage)

	s.True(ok)
	s.Equal(500, message.Status)
	s.Equal("test", message.Message)
}

func (s *mapperTestSuite) TestMiddleware() {
	s.Mapper.Register(func(err error) (message ErrorMessage) {
		if err == testErr {
			message = &errorMessage{500, testErr.Error()}
		}
		return
	})

	handler := s.Mapper.Middleware(func(w http.ResponseWriter, r *http.Request) error {
		if r.Method == "GET" {
			w.WriteHeader(200)
			return nil
		}

		return testErr
	})

	// request that should not generate eror
	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "http://example.com", nil)
	s.NoError(err)

	handler.ServeHTTP(w, r)
	s.Equal(200, w.Code)

	// request that generates error
	w = httptest.NewRecorder()
	r, err = http.NewRequest("POST", "http://example.com", nil)
	s.NoError(err)

	handler.ServeHTTP(w, r)
	s.Equal(500, w.Code)
	s.Equal("application/json", w.HeaderMap.Get("Content-Type"))

	var message errorMessage
	err = json.NewDecoder(w.Body).Decode(&message)
	s.NoError(err)
	s.Equal(500, message.Status)
	s.Equal("test", message.Message)
}

func (message errorMessage) HTTPStatus() int {
	return message.Status
}

func TestMapperTestSuite(t *testing.T) {
	suite.Run(t, new(mapperTestSuite))
}

type (
	errorMessage struct {
		Status  int
		Message string
	}

	mapperTestSuite struct {
		suite.Suite
		Mapper Mapper
	}
)

var (
	testErr = errors.New("test")
)
