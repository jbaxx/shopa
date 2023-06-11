package db

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type StubCommerceStore struct {
}

func Int(a int) *int {
	return &a
}

func (s *StubCommerceStore) GetChains() []Chain {
	c := Chain{
		Id: Int(1),
	}
	return []Chain{c}
}

func TestGetChains(t *testing.T) {
	store := &StubCommerceStore{}
	server := NewCommerceServer(store)

	t.Run("returns all the chains", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/chains"), nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		checkStatus(t, response, http.StatusOK)
		checkContentType(t, response, ContentApplicationJSON)

		want := []Chain{{Id: Int(1)}}
		got := getChainsFromResponse(t, response.Body)
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("ParseOutput() mismatch (-want +got):\n%s", diff)
		}
	})
}

func getChainsFromResponse(t *testing.T, body io.Reader) []Chain {
	t.Helper()
	c := []Chain{}
	err := json.NewDecoder(body).Decode(&c)
	if err != nil {
		t.Fatalf("unable to parse response from server %q into slice of Chain, '%v'", body, err)
	}
	return c
}

func checkStatus(t *testing.T, response *httptest.ResponseRecorder, want int) {
	t.Helper()
	got := response.Code
	if got != want {
		t.Errorf("incorrect status, got %d want %d", got, want)
	}
}

func checkContentType(t *testing.T, response *httptest.ResponseRecorder, want string) {
	t.Helper()
	got := response.Result().Header.Get("content-type")
	if got != want {
		t.Errorf("incorrect status, got %s want %s", got, want)
	}
}
