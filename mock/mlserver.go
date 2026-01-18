package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		filename := r.URL.Query().Get("filename")

		// Generate a random score between 0.0 and 1.0
		score := rand.Float64()

		fmt.Printf("Scoring image: %s -> %.4f\n", filename, score)
		fmt.Fprintf(w, "%f", score)
	})

	fmt.Println("Mock ML server running on :8001")
	if err := http.ListenAndServe(":8001", nil); err != nil {
		panic(err)
	}
}
