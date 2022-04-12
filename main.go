package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var (
	commitSha   string
	version     string
	backendHost string
)

func main() {
	backendHost = os.Getenv("BACKEND_HOST")
	if backendHost == "" {
		backendHost = "http://localhost:8081"
	}

	go serveProbe()
	serveHTTP()
}

func serveHTTP() {
	fmt.Printf("serving http :8080 | version: %s | commit sha: %s | backend "+
		"host: %s\n", version, commitSha, backendHost)

	m := http.NewServeMux()

	m.HandleFunc("/", rootHandler)

	s := http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: m,
	}

	log.Fatal(s.ListenAndServe())
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("GET /")

	resp, err := http.Get(backendHost)

	if err != nil {
		log.Printf("Request Failed: %s\n", err)
		fmt.Fprintf(w, "version: %s\ncommit sha: %s\nfeature enabled: %s",
			version, commitSha, "request failed")
		return
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed reading content from backend: %s\n", err)
		return
	}

	fmt.Fprintf(w, "feature enabled: %s\n", string(body))
}

func serveProbe() {
	fmt.Println("serving probe :8090")

	m := http.NewServeMux()

	m.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	s := http.Server{
		Addr:    "0.0.0.0:8090",
		Handler: m,
	}

	log.Fatal(s.ListenAndServe())
}
