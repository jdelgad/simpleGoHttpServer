package httpServer

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServer_NoRoutesRegistered_ServeHTTP(t *testing.T) {
	req, err := http.NewRequest("POST", "/hash", nil)
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}
	response := httptest.NewRecorder()

	s := NewServer()
	s.ServeHTTP(response, req)
	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected bad request status code not %v", response.Code)
	}
}

func TestServer_NoMatchingRoutesRegistered_ServeHTTP(t *testing.T) {
	req, err := http.NewRequest("POST", "/hash", nil)
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}
	response := httptest.NewRecorder()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})

	s := NewServer()
	s.Handle("/test123", handler)
	s.ServeHTTP(response, req)
	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected bad request status code not %v", response.Code)
	}
}

func TestServer_MatchingRegexRoute_ServeHTTP(t *testing.T) {
	req, err := http.NewRequest("POST", "/hash123", nil)
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}
	response := httptest.NewRecorder()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusContinue)
	})

	s := NewServer()
	s.Handle(`/hash\d+$`, handler)
	s.ServeHTTP(response, req)
	if response.Code != http.StatusContinue {
		t.Fatalf("expected status continue status code not %v", response.Code)
	}
}

func TestServer_MatchingPotentialConflictRegexRoute_ServeHTTP(t *testing.T) {
	req, err := http.NewRequest("POST", "/hash123", nil)
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}
	response := httptest.NewRecorder()

	handlerMatch := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusContinue)
	})
	handlerNotMatch := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	})

	s := NewServer()
	s.Handle(`/hash\d+$`, handlerMatch)
	s.Handle(`/hash`, handlerNotMatch)
	s.ServeHTTP(response, req)
	if response.Code != http.StatusContinue {
		t.Fatalf("expected status continue status code not %v", response.Code)
	}
}

func TestServer_TestHandler_ServeHTTP(t *testing.T) {

}
