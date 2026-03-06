package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/imran4u/Oolio/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testCouponSvc struct {
	ok map[string]bool
}

func (t testCouponSvc) ValidateCoupon(code string) bool { return t.ok[code] }

func TestOrderHandler_PlaceOrder(t *testing.T) {
	gin.SetMode(gin.TestMode)

	h := NewOrderHandler(testCouponSvc{ok: map[string]bool{
		"HAPPYHRS": true,
		"FIFTYOFF": true,
	}})

	r := gin.New()
	r.POST("/api/order", h.PlaceOrder)

	t.Run("bad json body", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, "/api/order", strings.NewReader(`{"items":`))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "message")
	})

	t.Run("invalid coupon", func(t *testing.T) {
		body := `{"couponCode":"SUPER100","items":[{"productId":"1","quantity":1}]}`
		req, err := http.NewRequest(http.MethodPost, "/api/order", strings.NewReader(body))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusUnprocessableEntity, w.Code)
		assert.Contains(t, w.Body.String(), "invalid coupon")
	})

	t.Run("invalid product", func(t *testing.T) {
		body := `{"items":[{"productId":"999","quantity":1}]}`
		req, err := http.NewRequest(http.MethodPost, "/api/order", strings.NewReader(body))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "invalid product")
	})

	t.Run("success with coupon", func(t *testing.T) {
		body := `{"couponCode":"HAPPYHRS","items":[{"productId":"1","quantity":1}]}`
		req, err := http.NewRequest(http.MethodPost, "/api/order", strings.NewReader(body))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusOK, w.Code)

		var got model.Order
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &got))
		assert.NotEmpty(t, got.ID)
		require.Len(t, got.Items, 1)
		require.Len(t, got.Products, 1)
		assert.Equal(t, "1", got.Products[0].ID)
	})
}

