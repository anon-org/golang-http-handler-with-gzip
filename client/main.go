package main

import (
	"context"
	"encoding/json"
	"example"
	"log"
	"net/http"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	body := &example.Payload{
		Number: 100,
	}

	log.Printf("create compressed request with %#v", body)
	req := example.MustCreateCompressedRequest(ctx, http.MethodPost, "http://localhost:8000/", body)
	defer req.Body.Close()

	log.Printf("send compressed request")
	resp, err := http.DefaultClient.Do(req)
	example.PanicIfErr(err)
	defer resp.Body.Close()

	log.Println("resp.Uncompressed?", resp.Uncompressed)
	var responsePayload *example.Payload
	if resp.Uncompressed {
		err = json.NewDecoder(resp.Body).Decode(&responsePayload)
	} else {
		responsePayload = example.MustReadCompressedBody[example.Payload](resp.Body)
	}
	log.Printf("read response %#v", responsePayload)
}
