package tests

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/imran4u/Oolio/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGETProduct_ListProducts(t *testing.T) {
	r := newTestRouter(mockCouponService{})

	tests := []struct {
		name       string
		apiKey     string
		wantStatus int
		wantMsg    string
	}{
		{
			name:       "successful request with valid header",
			apiKey:     "apitest",
			wantStatus: http.StatusOK,
		},
		{
			name:       "missing api_key header",
			apiKey:     "",
			wantStatus: http.StatusUnauthorized,
			wantMsg:    "missing api key",
		},
		{
			name:       "invalid api_key",
			apiKey:     "nope",
			wantStatus: http.StatusUnauthorized,
			wantMsg:    "invalid api key",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := newJSONRequest(t, http.MethodGet, "/api/product", nil)
			setAPIKey(req, tt.apiKey)

			w := performRequest(r, req)
			require.Equal(t, tt.wantStatus, w.Code)

			if tt.wantStatus == http.StatusOK {
				var got []model.Product
				require.NoError(t, json.Unmarshal(w.Body.Bytes(), &got))
				require.Len(t, got, 3)
				assert.Equal(t, "1", got[0].ID)
				assert.Equal(t, "2", got[1].ID)
				assert.Equal(t, "3", got[2].ID)
				return
			}

			var got map[string]any
			require.NoError(t, json.Unmarshal(w.Body.Bytes(), &got))
			assert.Equal(t, tt.wantMsg, got["message"])
		})
	}
}

func TestGETProduct_GetProduct(t *testing.T) {
	r := newTestRouter(mockCouponService{})

	tests := []struct {
		name       string
		apiKey     string
		productID  string
		wantStatus int
		wantMsg    string
		wantID     string
	}{
		{
			name:       "valid productId",
			apiKey:     "apitest",
			productID:  "1",
			wantStatus: http.StatusOK,
			wantID:     "1",
		},
		{
			name:       "non-existing productId",
			apiKey:     "apitest",
			productID:  "999",
			wantStatus: http.StatusNotFound,
			wantMsg:    "product not found",
		},
		{
			name:       "invalid productId format",
			apiKey:     "apitest",
			productID:  "abc",
			wantStatus: http.StatusBadRequest,
			wantMsg:    "invalid product id",
		},
		{
			name:       "missing api_key header",
			apiKey:     "",
			productID:  "1",
			wantStatus: http.StatusUnauthorized,
			wantMsg:    "missing api key",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := newJSONRequest(t, http.MethodGet, "/api/product/"+tt.productID, nil)
			setAPIKey(req, tt.apiKey)

			w := performRequest(r, req)
			require.Equal(t, tt.wantStatus, w.Code)

			if tt.wantStatus == http.StatusOK {
				var got model.Product
				require.NoError(t, json.Unmarshal(w.Body.Bytes(), &got))
				assert.Equal(t, tt.wantID, got.ID)
				assert.NotEmpty(t, got.Name)
				return
			}

			var got map[string]any
			require.NoError(t, json.Unmarshal(w.Body.Bytes(), &got))
			assert.Equal(t, tt.wantMsg, got["message"])
		})
	}
}

