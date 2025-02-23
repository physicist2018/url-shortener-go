// Code generated by MockGen. DO NOT EDIT.
// Source: internal/domain/service.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	domain "github.com/physicist2018/url-shortener-go/internal/domain"
)

// MockURLLinkService is a mock of URLLinkService interface.
type MockURLLinkService struct {
	ctrl     *gomock.Controller
	recorder *MockURLLinkServiceMockRecorder
}

// MockURLLinkServiceMockRecorder is the mock recorder for MockURLLinkService.
type MockURLLinkServiceMockRecorder struct {
	mock *MockURLLinkService
}

// NewMockURLLinkService creates a new mock instance.
func NewMockURLLinkService(ctrl *gomock.Controller) *MockURLLinkService {
	mock := &MockURLLinkService{ctrl: ctrl}
	mock.recorder = &MockURLLinkServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockURLLinkService) EXPECT() *MockURLLinkServiceMockRecorder {
	return m.recorder
}

// CreateShortURL mocks base method.
func (m *MockURLLinkService) CreateShortURL(ctx context.Context, link domain.URLLink) (domain.URLLink, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateShortURL", ctx, link)
	ret0, _ := ret[0].(domain.URLLink)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateShortURL indicates an expected call of CreateShortURL.
func (mr *MockURLLinkServiceMockRecorder) CreateShortURL(ctx, link interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateShortURL", reflect.TypeOf((*MockURLLinkService)(nil).CreateShortURL), ctx, link)
}

// FindAll mocks base method.
func (m *MockURLLinkService) FindAll(ctx context.Context, userID string) ([]domain.URLLink, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAll", ctx, userID)
	ret0, _ := ret[0].([]domain.URLLink)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindAll indicates an expected call of FindAll.
func (mr *MockURLLinkServiceMockRecorder) FindAll(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAll", reflect.TypeOf((*MockURLLinkService)(nil).FindAll), ctx, userID)
}

// GetOriginalURL mocks base method.
func (m *MockURLLinkService) GetOriginalURL(ctx context.Context, link domain.URLLink) (domain.URLLink, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOriginalURL", ctx, link)
	ret0, _ := ret[0].(domain.URLLink)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOriginalURL indicates an expected call of GetOriginalURL.
func (mr *MockURLLinkServiceMockRecorder) GetOriginalURL(ctx, link interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOriginalURL", reflect.TypeOf((*MockURLLinkService)(nil).GetOriginalURL), ctx, link)
}

// MarkURLsAsDeleted mocks base method.
func (m *MockURLLinkService) MarkURLsAsDeleted(ctx context.Context, links domain.DeleteRecordTask) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MarkURLsAsDeleted", ctx, links)
	ret0, _ := ret[0].(error)
	return ret0
}

// MarkURLsAsDeleted indicates an expected call of MarkURLsAsDeleted.
func (mr *MockURLLinkServiceMockRecorder) MarkURLsAsDeleted(ctx, links interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MarkURLsAsDeleted", reflect.TypeOf((*MockURLLinkService)(nil).MarkURLsAsDeleted), ctx, links)
}

// Ping mocks base method.
func (m *MockURLLinkService) Ping(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Ping", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// Ping indicates an expected call of Ping.
func (mr *MockURLLinkServiceMockRecorder) Ping(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ping", reflect.TypeOf((*MockURLLinkService)(nil).Ping), ctx)
}
