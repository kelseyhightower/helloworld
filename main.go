// Copyright 2017 Google Inc. All Rights Reserved.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
        "io/ioutil"
)

const (
	version = "v1.0.0"
)

func main() {
	log.Println("Starting helloworld application...")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello universe!\n")
	        response, err := http.Get("http://person.pool/info")

        	if err != nil {
            		fmt.Print(err.Error())
            		os.Exit(1)
        	}

        	responseData, err := ioutil.ReadAll(response.Body)
        	if err != nil {
            		log.Fatal(err)
        	}
        	fmt.Fprintf(w, string(responseData))
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})

	http.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, version)
	})

	s := http.Server{Addr: ":80"}
	go func() {
		log.Fatal(s.ListenAndServe())
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	log.Println("Shutdown signal received, exiting...")

	s.Shutdown(context.Background())
}
