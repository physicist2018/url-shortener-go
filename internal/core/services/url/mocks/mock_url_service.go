package mocks

import (
	"github.com/physicist2018/url-shortener-go/internal/core/models/urlmodels"
	"github.com/stretchr/testify/mock"
)

type MockURLService struct {
	mock.Mock
}

func (m *MockURLService) GenerateShortURL(originalURL string) (urlmodels.URL, error) {
	args := m.Called(originalURL)
	return args.Get(0).(urlmodels.URL), args.Error(1)
}

func (m *MockURLService) GetOriginalURL(shortURL string) (urlmodels.URL, error) {
	args := m.Called(shortURL)
	return args.Get(0).(urlmodels.URL), args.Error(1)
}
