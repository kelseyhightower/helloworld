package main

// https://golang.org/pkg/net/http/httptest

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_health(t *testing.T) {
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatalf("Test_health(%v) err:%v", "success", err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(health)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func Test_helloWorld(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatalf("Test_hello(%v) err:%v", "success", err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(helloWorld)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	resp := rr.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	// fmt.Printf("Status: %d\n", resp.StatusCode)
	// fmt.Printf("Content-Type: %s\n", resp.Header.Get("Content-Type"))
	// fmt.Printf("Body: %s\n", string(body))

	if success := `{"message": "Hello world!"}`; string(body) != success {
		t.Errorf("wrong response received, got %v want %v",
			string(body), success)
	}
}

func Test_getVersion(t *testing.T) {
	handler := http.HandlerFunc(getVersion)

	req := httptest.NewRequest("GET", "/version", nil)
	w := httptest.NewRecorder()
	handler(w, req)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	if success := fmt.Sprintf(`{"version": "%s"}`, version); string(body) != success {
		t.Errorf("wrong response received, got %v want %v",
			string(body), success)
	}
}

func Benchmark_health(b *testing.B) {
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		b.Fatalf("Test_health(%v) err:%v", "success", err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(health)

	for n := 0; n < b.N; n++ {
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			b.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}
	}
}

// func Benchmark_helloWorld(b *testing.B) {

// }

// func Benchmark_getVersion(b *testing.B) {

// }
