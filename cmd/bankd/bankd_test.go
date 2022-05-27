package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestLookupBSB(t *testing.T) {
	tests := []struct {
		name string
		bsb  string
		code int
	}{
		{name: "ok", bsb: "012-209", code: http.StatusOK},
		{name: "not found", bsb: "111-111", code: http.StatusNotFound},
		{name: "bad bsb", bsb: "abc", code: http.StatusBadRequest},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wr := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/bsb/", nil)
			req = mux.SetURLVars(req, map[string]string{
				"bsb": tt.bsb,
			})
			LookupBSB(wr, req)
			if tt.code != wr.Code {
				t.Errorf("LookupBSB wanted status code %d, but got %d", tt.code, wr.Code)
			}
		})
	}
}
