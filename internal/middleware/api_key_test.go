package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAPIKeyAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		apiKey     string
		wantStatus int
		wantMsg    string
	}{
		{
			name:       "missing api_key header",
			apiKey:     "",
			wantStatus: http.StatusUnauthorized,
			wantMsg:    "missing api key",
		},
		{
			name:       "invalid api_key header",
			apiKey:     "nope",
			wantStatus: http.StatusUnauthorized,
			wantMsg:    "invalid api key",
		},
		{
			name:       "valid api_key header",
			apiKey:     "apitest",
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.New()
			r.GET("/x", APIKeyAuth(), func(c *gin.Context) { c.Status(http.StatusOK) })

			req, err := http.NewRequest(http.MethodGet, "/x", nil)
			require.NoError(t, err)
			if tt.apiKey != "" {
				req.Header.Set("api_key", tt.apiKey)
			}
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			require.Equal(t, tt.wantStatus, w.Code)
			if tt.wantStatus == http.StatusOK {
				return
			}
			assert.Contains(t, w.Body.String(), tt.wantMsg)
		})
	}
}

