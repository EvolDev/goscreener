package index_test

import (
	"errors"
	"goscreener/internal/handlers/index"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockResponseWriter struct {
	mock.Mock
	status int
}

func (m *MockResponseWriter) Header() http.Header {
	return http.Header{}
}

func (m *MockResponseWriter) Write(data []byte) (int, error) {
	args := m.Called(data)
	return args.Int(0), args.Error(1)
}

func (m *MockResponseWriter) WriteHeader(statusCode int) {
	m.status = statusCode
}

// Success complete ServeHTTP
func TestHandler_ServeHTTP_Success(t *testing.T) {
	handler := index.NewIndexHandler()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Greater(t, rr.Body.Len(), 0)
}

// Test error for write http body
func TestHandler_ServeHTTP_WriteError(t *testing.T) {
	handler := index.NewIndexHandler()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	mockWriter := &MockResponseWriter{}
	mockWriter.On("Write", mock.Anything).Return(0, errors.New("write error"))

	handler.ServeHTTP(mockWriter, req)

	mockWriter.AssertExpectations(t)
	assert.Equal(t, http.StatusInternalServerError, mockWriter.status)
}
