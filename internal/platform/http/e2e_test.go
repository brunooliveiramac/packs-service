package httpapi_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	httpapi "github.com/brunooliveiramac/packs-service/internal/platform/http"
	"github.com/gin-gonic/gin"
)

func TestE2E_PacksCalc(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	httpapi.RegisterRoutes(r)

	srv := httptest.NewServer(r)
	t.Cleanup(srv.Close)

	type respPack struct {
		Size  int `json:"size"`
		Count int `json:"count"`
	}
	type respBody struct {
		Requested int        `json:"requested"`
		Shipped   int        `json:"shipped"`
		Packs     []respPack `json:"packs"`
	}

	tests := []struct {
		name        string
		qty         int
		wantShipped int
		want        map[int]int
	}{
		{"qty=1", 1, 250, map[int]int{250: 1}},
		{"qty=250", 250, 250, map[int]int{250: 1}},
		{"qty=251", 251, 500, map[int]int{500: 1}},
		{"qty=501", 501, 750, map[int]int{500: 1, 250: 1}},
		{"qty=12001", 12001, 12250, map[int]int{5000: 2, 2000: 1, 250: 1}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			body := map[string]int{"quantity": tc.qty}
			b, _ := json.Marshal(body)
			res, err := http.Post(srv.URL+"/api/packs/calc", "application/json", bytes.NewReader(b))
			if err != nil {
				t.Fatalf("post error: %v", err)
			}
			defer res.Body.Close()
			if res.StatusCode != http.StatusOK {
				t.Fatalf("status=%d want=%d", res.StatusCode, http.StatusOK)
			}
			var rb respBody
			if err := json.NewDecoder(res.Body).Decode(&rb); err != nil {
				t.Fatalf("decode: %v", err)
			}
			if rb.Shipped != tc.wantShipped {
				t.Fatalf("shipped=%d want=%d", rb.Shipped, tc.wantShipped)
			}
			got := map[int]int{}
			for _, p := range rb.Packs {
				got[p.Size] = p.Count
			}
			for sz, ct := range tc.want {
				if got[sz] != ct {
					t.Fatalf("pack %d got=%d want=%d", sz, got[sz], ct)
				}
			}
		})
	}
}


