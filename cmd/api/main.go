package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, "ok")
	})

	log.Println("🚀 API running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
