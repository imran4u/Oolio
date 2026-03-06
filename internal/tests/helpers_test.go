package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/imran4u/Oolio/internal/router"
)

type mockCouponService struct {
	validate func(code string) bool
}

func (m mockCouponService) ValidateCoupon(code string) bool {
	if m.validate == nil {
		return false
	}
	return m.validate(code)
}

func newTestRouter(couponSvc interface {
	ValidateCoupon(string) bool
}) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	router.RegisterRoutes(r, couponSvc)
	return r
}

func newJSONRequest(t *testing.T, method, path string, body any) *http.Request {
	t.Helper()
	var r io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("marshal body: %v", err)
		}
		r = bytes.NewReader(b)
	}
	req, err := http.NewRequest(method, path, r)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	return req
}

func newRawRequest(t *testing.T, method, path, rawBody string) *http.Request {
	t.Helper()
	req, err := http.NewRequest(method, path, bytes.NewBufferString(rawBody))
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	return req
}

func setAPIKey(req *http.Request, apiKey string) {
	if apiKey == "" {
		return
	}
	req.Header.Set("api_key", apiKey)
}

func performRequest(r http.Handler, req *http.Request) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

