package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

type alwaysValidCoupon struct{}

func (alwaysValidCoupon) ValidateCoupon(string) bool { return true }

func TestRegisterRoutes_WiresEndpoints(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	RegisterRoutes(r, alwaysValidCoupon{})

	tests := []struct {
		method string
		path   string
	}{
		{method: http.MethodGet, path: "/api/product"},
		{method: http.MethodGet, path: "/api/product/1"},
		{method: http.MethodPost, path: "/api/order"},
	}

	for _, tt := range tests {
		req, err := http.NewRequest(tt.method, tt.path, nil)
		require.NoError(t, err)
		req.Header.Set("api_key", "apitest")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		// We only assert that the route exists (i.e. not 404).
		require.NotEqual(t, http.StatusNotFound, w.Code)
	}
}

