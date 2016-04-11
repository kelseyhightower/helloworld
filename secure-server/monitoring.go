package main

import (
	"log"
	"net/http"
	"strconv"
)

var (
	healthzStatus   = 200
	readinessStatus = 200
)

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("/healthz returned %v", healthzStatus)
	w.WriteHeader(healthzStatus)
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	endpoint := r.FormValue("endpoint")
	status := r.FormValue("status")
	s, err := strconv.Atoi(status)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	switch endpoint {
	case "healthz":
		log.Printf("/healthz state set to %v", s)
		healthzStatus = s
	case "readiness":
		log.Printf("/readiness state set to %v", s)
		readinessStatus = s
	}

	w.WriteHeader(http.StatusOK)
}

func readinessHandler(w http.ResponseWriter, r *http.Request) {
	// check all subsystems
	log.Printf("/readiness returned %v", readinessStatus)
	w.WriteHeader(readinessStatus)
}
