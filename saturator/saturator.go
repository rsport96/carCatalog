package saturator

import (
	"catalog/models"
	"errors"
	jsonit "github.com/json-iterator/go"
	"github.com/valyala/fasthttp"
)

type Saturator interface {
	Saturate(regNum string) (*models.Car, error)
}

type saturator struct {
	url    string
	client *fasthttp.Client
}

var NilClientError = errors.New("client is nil")
var StatusBadRequestError = errors.New("external url returned status bad request")
var StatusInternalServerErrorError = errors.New("external url returned status internal server error")

func (s *saturator) Saturate(regNum string) (*models.Car, error) {
	if s.client == nil {
		return nil, NilClientError
	}

	req := fasthttp.AcquireRequest()
	req.SetRequestURI(s.url + "/info?regNum=" + regNum)
	req.Header.SetMethod(fasthttp.MethodGet)
	resp := fasthttp.AcquireResponse()
	if err := s.client.Do(req, resp); err != nil {
		return nil, err
	}
	if resp.StatusCode() == fasthttp.StatusBadRequest {
		return nil, StatusBadRequestError
	}
	if resp.StatusCode() == fasthttp.StatusInternalServerError {
		return nil, StatusInternalServerErrorError
	}
	car := &models.Car{}
	if err := jsonit.Unmarshal(resp.Body(), &car); err != nil {
		return nil, err
	}
	return car, nil
}

func NewSaturator(url string) Saturator {
	return &saturator{
		url:    url,
		client: &fasthttp.Client{},
	}
}
