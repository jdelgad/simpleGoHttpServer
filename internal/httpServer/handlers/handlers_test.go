package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

func Test_hash_get(t *testing.T) {
	shutdownCh := make(chan bool)
	h := NewHandlers(shutdownCh)

	req, err := http.NewRequest("GET", "/hash", nil)
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}

	req.Form = make(url.Values)
	req.Form["password"] = []string{"angryMonkey"}

	response := httptest.NewRecorder()
	handler := http.HandlerFunc(h.Hash)

	handler.ServeHTTP(response, req)

	if response.Code != http.StatusMethodNotAllowed {
		t.Fatalf("GET method not supported for /hash")
	}
}

func Test_hash_put(t *testing.T) {
	shutdownCh := make(chan bool)
	h := NewHandlers(shutdownCh)

	req, err := http.NewRequest("PUT", "/hash", nil)
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}

	req.Form = make(url.Values)
	req.Form["password"] = []string{"angryMonkey"}

	response := httptest.NewRecorder()
	handler := http.HandlerFunc(h.Hash)

	handler.ServeHTTP(response, req)

	if response.Code != http.StatusMethodNotAllowed {
		t.Fatalf("PUT method not supported for /hash")
	}
}


func Test_hash_multiple_post_success(t *testing.T) {
	shutdownCh := make(chan bool)
	h := NewHandlers(shutdownCh)

	req, err := http.NewRequest("POST", "/hash", nil)
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}

	req.Form = make(url.Values)
	req.Form["password"] = []string{"angryMonkey"}

	response := httptest.NewRecorder()
	handler := http.HandlerFunc(h.Hash)

	handler.ServeHTTP(response, req)

	if response.Code != http.StatusAccepted {
		t.Fatalf("expected accepted status code not %v", response.Code)
	}

	body := response.Body.String()
	if body != "1" {
		t.Fatalf("exepected first post call to /hash to return 1, not %v", body)
	}
	response.Body.Reset()


	handler.ServeHTTP(response, req)
	if response.Code != http.StatusAccepted {
		t.Fatalf("expected 2nd post call accepted status code not %v", response.Code)
	}

	body = response.Body.String()
	if body != "2" {
		t.Fatalf("exepected second post call to /hash to return 2, not %v", body)
	}

}

func Test_hash_post_missing_password_form(t *testing.T) {
	shutdownCh := make(chan bool)
	h := NewHandlers(shutdownCh)

	req, err := http.NewRequest("POST", "/hash", nil)
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}

	req.Form = make(url.Values)

	response := httptest.NewRecorder()
	handler := http.HandlerFunc(h.Hash)

	handler.ServeHTTP(response, req)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected accepted status code bad request not %v", response.Code)
	}
}

func Test_hash_4s_wait_hashid(t *testing.T) {
	shutdownCh := make(chan bool)
	h := NewHandlers(shutdownCh)

	req, err := http.NewRequest("POST", "/hash", nil)
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}

	req.Form = make(url.Values)
	req.Form["password"] = []string{"angryMonkey"}

	response := httptest.NewRecorder()
	handler := http.HandlerFunc(h.Hash)

	handler.ServeHTTP(response, req)
	if response.Code != http.StatusAccepted {
		t.Fatalf("expected accepted status code not %v", response.Code)
	}

	body := response.Body.String()
	if body != "1" {
		t.Fatalf("exepected first post call to /hash to return 1, not %v", body)
	}
	response.Body.Reset()

	req, err = http.NewRequest("GET", "/hash/1", nil)
	if err != nil {
		t.Fatalf("could not create GET request to /hash/1: %v", err)
	}

	response = httptest.NewRecorder()
	handler = http.HandlerFunc(h.HashID)

	time.Sleep(h.waitToStore - 1 * time.Second)
	handler.ServeHTTP(response, req)
	if response.Code != http.StatusAccepted {
		t.Fatalf("expected accepted status code not %v", response.Code)
	}
	body = response.Body.String()
	if body != "" {
		t.Fatalf("Did not expect any return value")
	}
}

func Test_hash_6s_wait_hashid(t *testing.T) {
	shutdownCh := make(chan bool)
	h := NewHandlers(shutdownCh)

	req, err := http.NewRequest("POST", "/hash", nil)
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}

	req.Form = make(url.Values)
	req.Form["password"] = []string{"angryMonkey"}

	response := httptest.NewRecorder()
	handler := http.HandlerFunc(h.Hash)

	handler.ServeHTTP(response, req)
	if response.Code != http.StatusAccepted {
		t.Fatalf("expected accepted status code not %v", response.Code)
	}

	body := response.Body.String()
	if body != "1" {
		t.Fatalf("exepected first post call to /hash to return 1, not %v", body)
	}
	response.Body.Reset()


	req, err = http.NewRequest("GET", "/hash/1", nil)
	if err != nil {
		t.Fatalf("could not create GET request to /hash/1: %v", err)
	}

	response = httptest.NewRecorder()
	handler = http.HandlerFunc(h.HashID)

	time.Sleep(h.waitToStore + 1 *time.Second)
	handler.ServeHTTP(response, req)
	if response.Code != http.StatusAccepted {
		t.Fatalf("expected accepted status code not %v", response.Code)
	}
	body = response.Body.String()
	if body != "ZEHhWB65gUlzdVwtDQArEyx+KVLzp/aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0Fd7su2A+gf7Q==" {
		t.Fatalf("Did not expect any return value, received %v", body)
	}
}

func Test_hash_3x_stats(t *testing.T) {
	shutdownCh := make(chan bool)
	h := NewHandlers(shutdownCh)

	req, err := http.NewRequest("POST", "/hash", nil)
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}

	req.Form = make(url.Values)
	req.Form["password"] = []string{"angryMonkey"}

	response := httptest.NewRecorder()
	handler := http.HandlerFunc(h.Hash)

	handler.ServeHTTP(response, req)
	handler.ServeHTTP(response, req)
	handler.ServeHTTP(response, req)

	req, err = http.NewRequest("GET", "/stats", nil)
	if err != nil {
		t.Fatalf("could not create GET request to /stats: %v", err)
	}

	response = httptest.NewRecorder()
	handler = http.HandlerFunc(h.Stats)

	// since the time taken is dependent on machine time, override for testing
	h.metrics.timeTaken = 20000
	handler.ServeHTTP(response, req)
	if response.Code != http.StatusAccepted {
		t.Fatalf("expected accepted status code not %v", response.Code)
	}
	body := response.Body.String()
	if body != `{"total":3,"average":6}` {
		t.Fatalf("Did not expect any return value, received %v", body)
	}
}

func Test_stats_no_post(t *testing.T) {
	shutdownCh := make(chan bool)
	h := NewHandlers(shutdownCh)

	req, err := http.NewRequest("GET", "/stats", nil)
	if err != nil {
		t.Fatalf("could not create GET request to /stats: %v", err)
	}

	response := httptest.NewRecorder()
	handler := http.HandlerFunc(h.Stats)

	handler.ServeHTTP(response, req)
	if response.Code != http.StatusAccepted {
		t.Fatalf("expected accepted status code not %v", response.Code)
	}
	body := response.Body.String()
	if body != `{"total":0,"average":0}` {
		t.Fatalf("Did not expect any return value, received %v", body)
	}
}

func Test_stats_post(t *testing.T) {
	shutdownCh := make(chan bool)
	h := NewHandlers(shutdownCh)

	req, err := http.NewRequest("POST", "/stats", nil)
	if err != nil {
		t.Fatalf("could not create GET request to /stats: %v", err)
	}

	response := httptest.NewRecorder()
	handler := http.HandlerFunc(h.Stats)

	handler.ServeHTTP(response, req)
	if response.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected method not allowed code not %v", response.Code)
	}
}