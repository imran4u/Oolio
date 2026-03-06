package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/imran4u/Oolio/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListProducts_Handler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/api/product", ListProducts)

	req, err := http.NewRequest(http.MethodGet, "/api/product", nil)
	require.NoError(t, err)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	var got []model.Product
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &got))
	require.Len(t, got, 3)
}

func TestGetProduct_Handler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/api/product/:productId", GetProduct)

	tests := []struct {
		name       string
		id         string
		wantStatus int
		wantMsg    string
	}{
		{name: "valid id", id: "1", wantStatus: http.StatusOK},
		{name: "invalid format", id: "abc", wantStatus: http.StatusBadRequest, wantMsg: "invalid product id"},
		{name: "not found", id: "999", wantStatus: http.StatusNotFound, wantMsg: "product not found"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, "/api/product/"+tt.id, nil)
			require.NoError(t, err)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			require.Equal(t, tt.wantStatus, w.Code)
			if tt.wantStatus == http.StatusOK {
				var got model.Product
				require.NoError(t, json.Unmarshal(w.Body.Bytes(), &got))
				assert.Equal(t, tt.id, got.ID)
				return
			}
			assert.Contains(t, w.Body.String(), tt.wantMsg)
		})
	}
}

