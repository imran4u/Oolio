package tests

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/imran4u/Oolio/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPOSTOrder_PlaceOrder(t *testing.T) {
	tests := []struct {
		name       string
		apiKey     string
		body       any
		rawBody    *string
		couponFn   func(code string) bool
		wantStatus int
		wantMsg    string
	}{
		{
			name:   "valid order with coupon HAPPYHRS",
			apiKey: "apitest",
			body: model.OrderReq{
				CouponCode: "HAPPYHRS",
				Items: []model.OrderItem{
					{ProductID: "1", Quantity: 1},
					{ProductID: "2", Quantity: 2},
				},
			},
			couponFn: func(code string) bool { return code == "HAPPYHRS" || code == "FIFTYOFF" },
			wantStatus: http.StatusOK,
		},
		{
			name:   "valid order with coupon FIFTYOFF",
			apiKey: "apitest",
			body: model.OrderReq{
				CouponCode: "FIFTYOFF",
				Items: []model.OrderItem{
					{ProductID: "1", Quantity: 1},
				},
			},
			couponFn: func(code string) bool { return code == "HAPPYHRS" || code == "FIFTYOFF" },
			wantStatus: http.StatusOK,
		},
		{
			name:   "invalid coupon SUPER100",
			apiKey: "apitest",
			body: model.OrderReq{
				CouponCode: "SUPER100",
				Items: []model.OrderItem{
					{ProductID: "1", Quantity: 1},
				},
			},
			couponFn:   func(code string) bool { return code == "HAPPYHRS" || code == "FIFTYOFF" },
			wantStatus: http.StatusUnprocessableEntity,
			wantMsg:    "invalid coupon",
		},
		{
			name:   "missing coupon",
			apiKey: "apitest",
			body: model.OrderReq{
				Items: []model.OrderItem{
					{ProductID: "3", Quantity: 3},
				},
			},
			couponFn:   func(code string) bool { return false },
			wantStatus: http.StatusOK,
		},
		{
			name:       "invalid request body",
			apiKey:     "apitest",
			rawBody:    ptr(`{"couponCode":`),
			couponFn:   func(code string) bool { return true },
			wantStatus: http.StatusBadRequest,
			wantMsg:    "EOF",
		},
		{
			name:       "missing api_key header",
			apiKey:     "",
			body:       model.OrderReq{Items: []model.OrderItem{{ProductID: "1", Quantity: 1}}},
			couponFn:   func(code string) bool { return true },
			wantStatus: http.StatusUnauthorized,
			wantMsg:    "missing api key",
		},
		{
			name:   "invalid product in items",
			apiKey: "apitest",
			body: model.OrderReq{
				Items: []model.OrderItem{
					{ProductID: "999", Quantity: 1},
				},
			},
			couponFn:   func(code string) bool { return true },
			wantStatus: http.StatusBadRequest,
			wantMsg:    "invalid product",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := newTestRouter(mockCouponService{validate: tt.couponFn})

			var req *http.Request
			if tt.rawBody != nil {
				req = newRawRequest(t, http.MethodPost, "/api/order", *tt.rawBody)
			} else {
				req = newJSONRequest(t, http.MethodPost, "/api/order", tt.body)
			}
			setAPIKey(req, tt.apiKey)

			w := performRequest(r, req)
			require.Equal(t, tt.wantStatus, w.Code)

			if tt.wantStatus == http.StatusOK {
				var got model.Order
				require.NoError(t, json.Unmarshal(w.Body.Bytes(), &got))
				assert.NotEmpty(t, got.ID)
				assert.NotEmpty(t, got.Items)
				assert.Len(t, got.Products, len(got.Items))
				return
			}

			var got map[string]any
			require.NoError(t, json.Unmarshal(w.Body.Bytes(), &got))
			if tt.wantStatus == http.StatusBadRequest && tt.rawBody != nil {
				// Gin's JSON binding error contains more detail; assert by substring.
				require.IsType(t, "", got["message"])
				assert.Contains(t, got["message"].(string), tt.wantMsg)
				return
			}
			assert.Equal(t, tt.wantMsg, got["message"])
		})
	}
}

func ptr(s string) *string { return &s }

