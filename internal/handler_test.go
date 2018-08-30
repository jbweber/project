package internal

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDateTimeHandler(t *testing.T) {
	// "monkey patch"
	// this could go wrong with parallel tests
	now = func() time.Time {
		return time.Unix(unixRefTime, 0)
	}

	req, err := http.NewRequest(http.MethodGet, "/datetime", nil)
	assert.NoError(t, err)

	h := http.HandlerFunc(DateTimeHandler)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code, "return code was %d not %d", rr.Code, http.StatusOK)
	assert.Equal(t, "{\"timestamp\":1136239445}\n", rr.Body.String(), "expected body to be empty")

}

func TestHealthHandler(t *testing.T) {

	tests := []struct {
		name         string
		method       string
		expectedCode int
		expectedBody string
	}{
		{
			"Health GET",
			http.MethodGet,
			http.StatusOK,
			"",
		},
		{
			"Health HEAD",
			http.MethodHead,
			http.StatusMethodNotAllowed,
			"method not allowed\n",
		},
		{
			"Health POST",
			http.MethodPost,
			http.StatusMethodNotAllowed,
			"method not allowed\n",
		},
		{
			"Health PUT",
			http.MethodPut,
			http.StatusMethodNotAllowed,
			"method not allowed\n",
		},
		{
			"Health PATCH",
			http.MethodPatch,
			http.StatusMethodNotAllowed,
			"method not allowed\n",
		},
		{
			"Health DELETE",
			http.MethodDelete,
			http.StatusMethodNotAllowed,
			"method not allowed\n",
		},
		{
			"Health OPTIONS",
			http.MethodOptions,
			http.StatusMethodNotAllowed,
			"method not allowed\n",
		},
		{
			"Health TRACE",
			http.MethodTrace,
			http.StatusMethodNotAllowed,
			"method not allowed\n",
		},
		{
			"Health CONNECT",
			http.MethodConnect,
			http.StatusMethodNotAllowed,
			"method not allowed\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, "/health", nil)
			assert.NoError(t, err)

			h := http.HandlerFunc(HealthHandler)
			rr := httptest.NewRecorder()
			h.ServeHTTP(rr, req)
			assert.Equal(t, tt.expectedCode, rr.Code, "return code was %d not %d", rr.Code, tt.expectedCode)
			assert.Equal(t, tt.expectedBody, rr.Body.String(), "expected body to be empty")
		})
	}
}
