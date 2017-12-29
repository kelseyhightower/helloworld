package main

import (
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

// func Test_helloWorld(t *testing.T) {

// }

// func Test_getVersion(t *testing.T) {

// }

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
