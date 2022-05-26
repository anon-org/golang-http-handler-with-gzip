package main

import (
	"example"
	"log"
	"net/http"
)

func main() {
	log.Println("server listening at :8000")
	http.ListenAndServe(":8000", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body := example.MustReadCompressedBody[example.Payload](r.Body)
		body.Number++

		example.MustWriteCompressedResponse(w, body)
	}))
}
