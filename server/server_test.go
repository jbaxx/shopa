package server

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/jbaxx/shopa/db"
)

type StubCommerceStore struct {
}

func String(a string) *string {
	return &a
}

func (s *StubCommerceStore) GetChains() []db.Chain {
	c := db.Chain{
		Id: db.String("1"),
	}
	return []db.Chain{c}
}

func TestGetChains(t *testing.T) {
	store := &StubCommerceStore{}
	srv := NewCommerceServer(store, "", log.New(os.Stdout, "", log.LstdFlags))

	t.Run("returns all the chains", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/chains"), nil)
		response := httptest.NewRecorder()

		srv.server.Handler.ServeHTTP(response, request)

		checkStatus(t, response, http.StatusOK)
		checkContentType(t, response, ContentApplicationJSON)

		want := []db.Chain{{Id: String("1")}}
		got := getChainsFromResponse(t, response.Body)
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("ParseOutput() mismatch (-want +got):\n%s", diff)
		}
	})
}

func getChainsFromResponse(t *testing.T, body io.Reader) []db.Chain {
	t.Helper()
	c := []db.Chain{}
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
